package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

// Status represents the status of Docker-related checks
type Status struct {
	DaemonRunning bool   `json:"daemonRunning"`
	Version       string `json:"version"`
	Info          string `json:"info"`
}

// Metrics represents various Docker-related metrics
type Metrics struct {
	ImagesCount   int    `json:"imagesCount"`
	ContainersAll int    `json:"containersAll"`
	ContainersUp  int    `json:"containersUp"`
	DiskUsage     string `json:"diskUsage"`
	NetworkStatus string `json:"networkStatus"`
}

// GetStatus checks if Docker daemon is running and returns version info
func GetStatus() Status {
	status := Status{
		DaemonRunning: false,
		Version:       "Not Available",
		Info:          "Docker daemon not running",
	}

	// Check if docker daemon is running
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err == nil {
		status.DaemonRunning = true
		status.Info = "Docker daemon is running"
	}

	// Get Docker version
	versionCmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	if output, err := versionCmd.Output(); err == nil {
		status.Version = strings.TrimSpace(string(output))
	}

	return status
}

// GetMetrics returns various Docker-related metrics
func GetMetrics() Metrics {
	metrics := Metrics{
		ImagesCount:   0,
		ContainersAll: 0,
		ContainersUp:  0,
		DiskUsage:     "Not Available",
		NetworkStatus: "Not Available",
	}

	// Count images
	imagesCmd := exec.Command("docker", "images", "--format", "{{.ID}}")
	if output, err := imagesCmd.Output(); err == nil {
		metrics.ImagesCount = len(strings.Split(strings.TrimSpace(string(output)), "\n"))
	}

	// Count containers
	containersCmd := exec.Command("docker", "ps", "-a", "--format", "{{.Status}}")
	if output, err := containersCmd.Output(); err == nil {
		containers := strings.Split(strings.TrimSpace(string(output)), "\n")
		metrics.ContainersAll = len(containers)
		// Count running containers
		for _, status := range containers {
			if strings.Contains(status, "Up") {
				metrics.ContainersUp++
			}
		}
	}

	// Get disk usage
	diskCmd := exec.Command("docker", "system", "df", "--format", "{{.Size}}")
	if output, err := diskCmd.Output(); err == nil {
		metrics.DiskUsage = strings.TrimSpace(string(output))
	}

	// Check network status
	networkCmd := exec.Command("docker", "network", "ls", "--format", "{{.Name}}")
	if output, err := networkCmd.Output(); err == nil {
		networks := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(networks) > 0 {
			metrics.NetworkStatus = fmt.Sprintf("%d networks available", len(networks))
		} else {
			metrics.NetworkStatus = "No networks found"
		}
	}

	return metrics
}
