package process

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo represents information about a running process
type ProcessInfo struct {
	PID         int       `json:"pid"`
	Name        string    `json:"name"`
	CommandLine string    `json:"commandLine"`
	Username    string    `json:"username"`
	CPUPercent  float64   `json:"cpuPercent"`
	MemoryUsage uint64    `json:"memoryUsage"`
	StartTime   time.Time `json:"startTime"`
}

// PortInfo represents a port being used by a process
type PortInfo struct {
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"` // TCP or UDP
	LocalAddr string `json:"localAddr"`
	State     string `json:"state"` // LISTEN, ESTABLISHED, etc.
	PID       int    `json:"pid"`
}

// ProcessWithPorts combines process info with its open ports
type ProcessWithPorts struct {
	Process ProcessInfo `json:"process"`
	Ports   []PortInfo  `json:"ports"`
}

// GetAllProcesses returns information about all running processes
func GetAllProcesses() ([]ProcessInfo, error) {
	// Use gopsutil to get all processes
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("error getting processes: %w", err)
	}

	var result []ProcessInfo
	for _, p := range processes {
		// Skip processes that we can't access
		if !processAccessible(p) {
			continue
		}

		name, _ := p.Name()
		cmdline, _ := p.Cmdline()
		username, _ := p.Username()
		cpuPercent, _ := p.CPUPercent()
		memInfo, _ := p.MemoryInfo()
		createTime, _ := p.CreateTime()

		startTime := time.Unix(createTime/1000, 0)

		var memUsage uint64
		if memInfo != nil {
			memUsage = memInfo.RSS
		}

		result = append(result, ProcessInfo{
			PID:         int(p.Pid),
			Name:        name,
			CommandLine: cmdline,
			Username:    username,
			CPUPercent:  cpuPercent,
			MemoryUsage: memUsage,
			StartTime:   startTime,
		})
	}

	return result, nil
}

// processAccessible checks if we can access the process information
func processAccessible(p *process.Process) bool {
	_, err := p.Name()
	return err == nil
}

// GetAllPorts returns information about all open ports
func GetAllPorts() ([]PortInfo, error) {
	switch runtime.GOOS {
	case "darwin":
		return getMacPorts()
	case "linux":
		return getLinuxPorts()
	case "windows":
		return getWindowsPorts()
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// executeCommandWithTimeout runs a command with a timeout
func executeCommandWithTimeout(timeout time.Duration, name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Output()
}

// getMacPorts gets port information on macOS using lsof
func getMacPorts() ([]PortInfo, error) {
	var result []PortInfo

	// On macOS, use lsof to get port information with a 5-second timeout
	output, err := executeCommandWithTimeout(5*time.Second, "lsof", "-i", "-P", "-n")
	if err != nil {
		return nil, fmt.Errorf("error executing lsof: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		// Skip header line
		if i == 0 || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		// Parse PID
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		// Parse address information
		addrInfo := fields[8]
		if !strings.Contains(addrInfo, ":") {
			continue
		}

		parts := strings.Split(addrInfo, "->")
		localAddr := parts[0]

		// Extract port from local address
		addrParts := strings.Split(localAddr, ":")
		if len(addrParts) < 2 {
			continue
		}

		portStr := addrParts[len(addrParts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		// Determine protocol (TCP/UDP)
		protocol := "TCP"
		if strings.Contains(fields[7], "UDP") {
			protocol = "UDP"
		}

		// Determine state
		state := "UNKNOWN"
		if len(fields) >= 10 {
			state = fields[9]
		} else if strings.Contains(fields[7], "LISTEN") {
			state = "LISTEN"
		}

		result = append(result, PortInfo{
			Port:      port,
			Protocol:  protocol,
			LocalAddr: localAddr,
			State:     state,
			PID:       pid,
		})
	}

	return result, nil
}

// getLinuxPorts gets port information on Linux using ss
func getLinuxPorts() ([]PortInfo, error) {
	var result []PortInfo

	// On Linux, use ss to get port information with a 5-second timeout
	output, err := executeCommandWithTimeout(5*time.Second, "ss", "-tuln", "-p")
	if err != nil {
		return nil, fmt.Errorf("error executing ss: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		// Skip header line
		if i == 0 || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// Parse protocol
		protocol := "TCP"
		if strings.Contains(fields[0], "udp") {
			protocol = "UDP"
		}

		// Parse local address
		localAddr := fields[4]
		if !strings.Contains(localAddr, ":") {
			continue
		}

		// Extract port from local address
		addrParts := strings.Split(localAddr, ":")
		if len(addrParts) < 2 {
			continue
		}

		portStr := addrParts[len(addrParts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		// Parse PID
		pid := 0
		if len(fields) >= 6 {
			pidInfo := fields[5]
			if strings.Contains(pidInfo, "pid=") {
				pidParts := strings.Split(pidInfo, "pid=")
				if len(pidParts) >= 2 {
					pidStr := strings.Split(pidParts[1], ",")[0]
					pid, _ = strconv.Atoi(pidStr)
				}
			}
		}

		result = append(result, PortInfo{
			Port:      port,
			Protocol:  protocol,
			LocalAddr: localAddr,
			State:     "LISTEN",
			PID:       pid,
		})
	}

	return result, nil
}

// getWindowsPorts gets port information on Windows using netstat
func getWindowsPorts() ([]PortInfo, error) {
	var result []PortInfo

	// On Windows, use netstat to get port information with a 5-second timeout
	output, err := executeCommandWithTimeout(5*time.Second, "netstat", "-ano")
	if err != nil {
		return nil, fmt.Errorf("error executing netstat: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		// Skip header lines
		if i < 4 || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// Parse protocol
		protocol := fields[0]
		if protocol != "TCP" && protocol != "UDP" {
			continue
		}

		// Parse local address
		localAddr := fields[1]
		if !strings.Contains(localAddr, ":") {
			continue
		}

		// Extract port from local address
		addrParts := strings.Split(localAddr, ":")
		if len(addrParts) < 2 {
			continue
		}

		portStr := addrParts[len(addrParts)-1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		// Parse state
		state := "UNKNOWN"
		if protocol == "TCP" && len(fields) >= 4 {
			state = fields[3]
		}

		// Parse PID
		pid := 0
		if len(fields) >= 5 {
			pid, _ = strconv.Atoi(fields[4])
		}

		result = append(result, PortInfo{
			Port:      port,
			Protocol:  protocol,
			LocalAddr: localAddr,
			State:     state,
			PID:       pid,
		})
	}

	return result, nil
}

// GetProcessesWithPorts combines process and port information
func GetProcessesWithPorts() ([]ProcessWithPorts, error) {
	processes, err := GetAllProcesses()
	if err != nil {
		return nil, err
	}

	ports, err := GetAllPorts()
	if err != nil {
		// If we can't get port information, just return processes without ports
		var result []ProcessWithPorts
		for _, p := range processes {
			result = append(result, ProcessWithPorts{
				Process: p,
				Ports:   []PortInfo{},
			})
		}
		return result, nil
	}

	// Create a map of PID to process info
	processMap := make(map[int]ProcessInfo)
	for _, p := range processes {
		processMap[p.PID] = p
	}

	// Create a map of PID to ports
	portMap := make(map[int][]PortInfo)
	for _, port := range ports {
		portMap[port.PID] = append(portMap[port.PID], port)
	}

	// Combine the information
	var result []ProcessWithPorts
	for pid, process := range processMap {
		result = append(result, ProcessWithPorts{
			Process: process,
			Ports:   portMap[pid],
		})
	}

	return result, nil
}

// SearchProcessesByPort finds processes using a specific port
func SearchProcessesByPort(port int) ([]ProcessWithPorts, error) {
	allProcesses, err := GetProcessesWithPorts()
	if err != nil {
		return nil, err
	}

	var result []ProcessWithPorts
	for _, p := range allProcesses {
		for _, portInfo := range p.Ports {
			if portInfo.Port == port {
				result = append(result, p)
				break
			}
		}
	}

	return result, nil
}

// KillProcess terminates a process by PID
func KillProcess(pid int) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	default: // macOS and Linux
		cmd = exec.Command("kill", "-9", strconv.Itoa(pid))
	}

	// Set a timeout for the kill command
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	return cmd.Run()
}

// FormatBytes converts bytes to a human-readable string
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
