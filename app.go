package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"DevEx/internal/docker"
	"DevEx/internal/history"
	"DevEx/internal/network"
	"DevEx/internal/process"
	"DevEx/internal/system"
)

// App struct
type App struct {
	ctx            context.Context
	db             *history.DB
	collector      *history.Collector
	processManager *process.Manager
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Initialize the database
	db, err := history.NewDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create a new collector with 10-second interval
	collector := history.NewCollector(db, 10*time.Second)

	// Create a new process manager with 30-second update interval (increased from 5 seconds)
	processManager := process.NewManager(30 * time.Second)

	// Set a limit on the number of processes to display
	processManager.SetMaxProcesses(300)

	return &App{
		db:             db,
		collector:      collector,
		processManager: processManager,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start the metrics collector
	a.collector.Start()

	// Start the process manager
	a.processManager.Start()

	// Start pprof server
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// Stop the metrics collector
	a.collector.Stop()

	// Stop the process manager
	a.processManager.Stop()

	// Close the database connection
	if err := a.db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}

// GetCPUInfo returns formatted CPU usage percentage
func (a *App) GetCPUInfo() string {
	return system.GetCPUInfo()
}

// GetCPUDetails returns detailed CPU information
func (a *App) GetCPUDetails() string {
	return system.GetCPUDetails()
}

// GetRAMInfo returns formatted RAM usage percentage
func (a *App) GetRAMInfo() string {
	return system.GetRAMInfo()
}

// GetRAMDetails returns detailed RAM information
func (a *App) GetRAMDetails() string {
	return system.GetRAMDetails()
}

// GetDiskInfo returns formatted disk usage percentage
func (a *App) GetDiskInfo() string {
	return system.GetDiskInfo()
}

// GetDiskDetails returns detailed disk information
func (a *App) GetDiskDetails() string {
	return system.GetDiskDetails()
}

// GetDockerStatus returns the Docker daemon status
func (a *App) GetDockerStatus() docker.Status {
	return docker.GetStatus()
}

// GetDockerMetrics returns Docker-related metrics
func (a *App) GetDockerMetrics() docker.Metrics {
	return docker.GetMetrics()
}

// GetNetworkStatus returns network status information
func (a *App) GetNetworkStatus() network.Status {
	return network.GetStatus()
}

// GetCPUHistory returns CPU usage history for the specified duration
func (a *App) GetCPUHistory(minutes int) []history.TimeSeriesPoint {
	duration := time.Duration(minutes) * time.Minute
	data, err := a.db.GetCPUHistory(duration)
	if err != nil {
		log.Printf("Error retrieving CPU history: %v", err)
		return []history.TimeSeriesPoint{}
	}
	return data
}

// GetRAMHistory returns RAM usage history for the specified duration
func (a *App) GetRAMHistory(minutes int) []history.TimeSeriesPoint {
	duration := time.Duration(minutes) * time.Minute
	data, err := a.db.GetRAMHistory(duration)
	if err != nil {
		log.Printf("Error retrieving RAM history: %v", err)
		return []history.TimeSeriesPoint{}
	}
	return data
}

// GetDiskHistory returns disk usage history for the specified duration
func (a *App) GetDiskHistory(minutes int) []history.TimeSeriesPoint {
	duration := time.Duration(minutes) * time.Minute
	data, err := a.db.GetDiskHistory(duration)
	if err != nil {
		log.Printf("Error retrieving disk history: %v", err)
		return []history.TimeSeriesPoint{}
	}
	return data
}

// GetDockerHistory returns Docker metrics history for the specified duration
func (a *App) GetDockerHistory(minutes int) []history.DockerMetrics {
	duration := time.Duration(minutes) * time.Minute
	data, err := a.db.GetDockerHistory(duration)
	if err != nil {
		log.Printf("Error retrieving Docker history: %v", err)
		return []history.DockerMetrics{}
	}
	return data
}

// GetNetworkHistory returns network metrics history for the specified duration
func (a *App) GetNetworkHistory(minutes int) []history.NetworkMetrics {
	duration := time.Duration(minutes) * time.Minute
	data, err := a.db.GetNetworkHistory(duration)
	if err != nil {
		log.Printf("Error retrieving network history: %v", err)
		return []history.NetworkMetrics{}
	}
	return data
}

// GetAllProcesses returns all running processes with port information
func (a *App) GetAllProcesses() []process.ProcessWithPorts {
	return a.processManager.GetProcesses()
}

// SearchProcessesByPort finds processes using a specific port
func (a *App) SearchProcessesByPort(port int) []process.ProcessWithPorts {
	return a.processManager.SearchByPort(port)
}

// KillProcess terminates a process by PID
func (a *App) KillProcess(pid int) error {
	return a.processManager.KillProcessByPID(pid)
}

// FormatProcessBytes formats bytes to human-readable string
func (a *App) FormatProcessBytes(bytes uint64) string {
	return process.FormatBytes(bytes)
}

// Shutdown is called when the app is closing
func (a *App) Shutdown() {
	runtime.Quit(a.ctx)
}
