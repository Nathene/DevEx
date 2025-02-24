package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// App struct
type App struct {
	ctx    context.Context
	ticker *time.Ticker
}

func NewApp() *App {
	return &App{}
}

// Add after the NewApp() function
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Create a ticker that ticks every 2 seconds
	a.ticker = time.NewTicker(2 * time.Second)

	// Start the update loop in a goroutine
	go func() {
		for {
			select {
			case <-a.ctx.Done():
				if a.ticker != nil {
					a.ticker.Stop()
				}
				return
			case <-a.ticker.C:
				// The frontend will poll for updates, so we don't need to emit events
				continue
			}
		}
	}()
}

func formatBytes(bytes uint64) string {
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

// GetCPUInfo returns formatted CPU information
func (a *App) GetCPUInfo() string {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return "CPU: Error"
	}
	if len(cpuPercent) == 0 {
		return "CPU: No data"
	}
	return fmt.Sprintf("CPU: %.1f%%", math.Round(cpuPercent[0]*10)/10)
}

// GetCPUDetails returns detailed CPU information
func (a *App) GetCPUDetails() string {
	cores, _ := cpu.Counts(false)
	threads, _ := cpu.Counts(true)
	return fmt.Sprintf("Cores: %d\nThreads: %d", cores, threads)
}

// GetRAMInfo returns formatted RAM usage
func (a *App) GetRAMInfo() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "RAM: Error"
	}
	return fmt.Sprintf("RAM: %.1f%%", math.Round(memInfo.UsedPercent*10)/10)
}

// GetRAMDetails returns detailed RAM information
func (a *App) GetRAMDetails() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "RAM: Error"
	}
	return fmt.Sprintf("Used: %s\nTotal: %s",
		formatBytes(memInfo.Used),
		formatBytes(memInfo.Total))
}

// GetDiskInfo returns formatted disk usage
func (a *App) GetDiskInfo() string {
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return "Disk: Error"
	}
	return fmt.Sprintf("Disk: %.1f%%", math.Round(diskInfo.UsedPercent*10)/10)
}

// GetDiskDetails returns detailed disk information
func (a *App) GetDiskDetails() string {
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return "Disk: Error"
	}
	return fmt.Sprintf("Used: %s\nTotal: %s",
		formatBytes(diskInfo.Used),
		formatBytes(diskInfo.Total))
}
