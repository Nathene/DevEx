package main

import (
	"context"
	_ "net/http/pprof"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"DevEx/internal/docker"
	"DevEx/internal/network"
	"DevEx/internal/system"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

// Shutdown is called when the app is closing
func (a *App) Shutdown() {
	runtime.Quit(a.ctx)
}
