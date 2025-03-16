package history

import (
	"DevEx/internal/docker"
	"DevEx/internal/network"
	"DevEx/internal/system"
	"context"
	"log"
	"strconv"
	"strings"
	"time"
)

// Collector periodically collects system metrics and stores them in the database
type Collector struct {
	db         *DB
	interval   time.Duration
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewCollector creates a new metrics collector
func NewCollector(db *DB, interval time.Duration) *Collector {
	ctx, cancel := context.WithCancel(context.Background())
	return &Collector{
		db:         db,
		interval:   interval,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// Start begins collecting metrics at the specified interval
func (c *Collector) Start() {
	log.Printf("Starting metrics collector with interval: %v", c.interval)

	// Collect metrics immediately
	c.collectMetrics()

	// Then collect at regular intervals
	ticker := time.NewTicker(c.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				c.collectMetrics()
			case <-c.ctx.Done():
				log.Println("Metrics collector stopped")
				return
			}
		}
	}()
}

// Stop halts the metrics collection
func (c *Collector) Stop() {
	c.cancelFunc()
}

// collectMetrics gathers all system metrics and stores them in the database
func (c *Collector) collectMetrics() {
	now := time.Now()

	// Collect and store CPU metrics
	c.collectCPUMetrics(now)

	// Collect and store RAM metrics
	c.collectRAMMetrics(now)

	// Collect and store disk metrics
	c.collectDiskMetrics(now)

	// Collect and store Docker metrics
	c.collectDockerMetrics(now)

	// Collect and store network metrics
	c.collectNetworkMetrics(now)
}

// collectCPUMetrics collects and stores CPU metrics
func (c *Collector) collectCPUMetrics(timestamp time.Time) {
	cpuInfo := system.GetCPUInfo()
	if cpuInfo == "CPU: Error" || cpuInfo == "CPU: No data" {
		log.Println("Error collecting CPU metrics")
		return
	}

	// Extract percentage from string like "CPU: 45.2%"
	parts := strings.Split(cpuInfo, " ")
	if len(parts) != 2 {
		log.Printf("Unexpected CPU info format: %s", cpuInfo)
		return
	}

	// Remove the % sign
	usageStr := strings.TrimSuffix(parts[1], "%")
	usage, err := strconv.ParseFloat(usageStr, 64)
	if err != nil {
		log.Printf("Error parsing CPU usage: %v", err)
		return
	}

	metrics := CPUMetrics{
		Timestamp: timestamp,
		Usage:     usage,
	}

	if err := c.db.StoreCPUMetrics(metrics); err != nil {
		log.Printf("Error storing CPU metrics: %v", err)
	}
}

// collectRAMMetrics collects and stores RAM metrics
func (c *Collector) collectRAMMetrics(timestamp time.Time) {
	ramInfo := system.GetRAMInfo()
	ramDetails := system.GetRAMDetails()

	if ramInfo == "RAM: Error" {
		log.Println("Error collecting RAM metrics")
		return
	}

	// Extract percentage from string like "RAM: 65.7%"
	parts := strings.Split(ramInfo, " ")
	if len(parts) != 2 {
		log.Printf("Unexpected RAM info format: %s", ramInfo)
		return
	}

	// Remove the % sign
	usageStr := strings.TrimSuffix(parts[1], "%")
	usage, err := strconv.ParseFloat(usageStr, 64)
	if err != nil {
		log.Printf("Error parsing RAM usage: %v", err)
		return
	}

	// Parse RAM details
	// Format: "Used: 8.5 GB\nTotal: 16.0 GB"
	detailsLines := strings.Split(ramDetails, "\n")
	if len(detailsLines) != 2 {
		log.Printf("Unexpected RAM details format: %s", ramDetails)
		return
	}

	// Extract used bytes (approximate)
	usedParts := strings.Split(detailsLines[0], " ")
	if len(usedParts) < 3 {
		log.Printf("Unexpected RAM used format: %s", detailsLines[0])
		return
	}

	// Extract total bytes (approximate)
	totalParts := strings.Split(detailsLines[1], " ")
	if len(totalParts) < 3 {
		log.Printf("Unexpected RAM total format: %s", detailsLines[1])
		return
	}

	// Convert to bytes (approximate)
	usedValue, err := strconv.ParseFloat(usedParts[1], 64)
	if err != nil {
		log.Printf("Error parsing RAM used value: %v", err)
		return
	}

	totalValue, err := strconv.ParseFloat(totalParts[1], 64)
	if err != nil {
		log.Printf("Error parsing RAM total value: %v", err)
		return
	}

	// Convert to bytes based on unit (GB, MB, etc.)
	usedUnit := usedParts[2]
	totalUnit := totalParts[2]

	usedBytes := convertToBytes(usedValue, usedUnit)
	totalBytes := convertToBytes(totalValue, totalUnit)

	metrics := RAMMetrics{
		Timestamp:  timestamp,
		Usage:      usage,
		UsedBytes:  usedBytes,
		TotalBytes: totalBytes,
	}

	if err := c.db.StoreRAMMetrics(metrics); err != nil {
		log.Printf("Error storing RAM metrics: %v", err)
	}
}

// collectDiskMetrics collects and stores disk metrics
func (c *Collector) collectDiskMetrics(timestamp time.Time) {
	diskInfo := system.GetDiskInfo()
	diskDetails := system.GetDiskDetails()

	if diskInfo == "Disk: Error" {
		log.Println("Error collecting disk metrics")
		return
	}

	// Extract percentage from string like "Disk: 75.3%"
	parts := strings.Split(diskInfo, " ")
	if len(parts) != 2 {
		log.Printf("Unexpected disk info format: %s", diskInfo)
		return
	}

	// Remove the % sign
	usageStr := strings.TrimSuffix(parts[1], "%")
	usage, err := strconv.ParseFloat(usageStr, 64)
	if err != nil {
		log.Printf("Error parsing disk usage: %v", err)
		return
	}

	// Parse disk details
	// Format: "Used: 256.7 GB\nTotal: 512.0 GB"
	detailsLines := strings.Split(diskDetails, "\n")
	if len(detailsLines) != 2 {
		log.Printf("Unexpected disk details format: %s", diskDetails)
		return
	}

	// Extract used bytes (approximate)
	usedParts := strings.Split(detailsLines[0], " ")
	if len(usedParts) < 3 {
		log.Printf("Unexpected disk used format: %s", detailsLines[0])
		return
	}

	// Extract total bytes (approximate)
	totalParts := strings.Split(detailsLines[1], " ")
	if len(totalParts) < 3 {
		log.Printf("Unexpected disk total format: %s", detailsLines[1])
		return
	}

	// Convert to bytes (approximate)
	usedValue, err := strconv.ParseFloat(usedParts[1], 64)
	if err != nil {
		log.Printf("Error parsing disk used value: %v", err)
		return
	}

	totalValue, err := strconv.ParseFloat(totalParts[1], 64)
	if err != nil {
		log.Printf("Error parsing disk total value: %v", err)
		return
	}

	// Convert to bytes based on unit (GB, MB, etc.)
	usedUnit := usedParts[2]
	totalUnit := totalParts[2]

	usedBytes := convertToBytes(usedValue, usedUnit)
	totalBytes := convertToBytes(totalValue, totalUnit)

	metrics := DiskMetrics{
		Timestamp:  timestamp,
		Usage:      usage,
		UsedBytes:  usedBytes,
		TotalBytes: totalBytes,
	}

	if err := c.db.StoreDiskMetrics(metrics); err != nil {
		log.Printf("Error storing disk metrics: %v", err)
	}
}

// collectDockerMetrics collects and stores Docker metrics
func (c *Collector) collectDockerMetrics(timestamp time.Time) {
	dockerStatus := docker.GetStatus()
	dockerMetrics := docker.GetMetrics()

	metrics := DockerMetrics{
		Timestamp:         timestamp,
		DaemonRunning:     dockerStatus.DaemonRunning,
		ContainersRunning: dockerMetrics.ContainersUp,
		ContainersTotal:   dockerMetrics.ContainersAll,
		ImagesCount:       dockerMetrics.ImagesCount,
	}

	if err := c.db.StoreDockerMetrics(metrics); err != nil {
		log.Printf("Error storing Docker metrics: %v", err)
	}
}

// collectNetworkMetrics collects and stores network metrics
func (c *Collector) collectNetworkMetrics(timestamp time.Time) {
	networkStatus := network.GetStatus()

	metrics := NetworkMetrics{
		Timestamp:         timestamp,
		InternetConnected: networkStatus.InternetConnected,
		PingLatency:       networkStatus.PingLatency,
		DNSWorking:        networkStatus.DNSStatus == "Working",
	}

	if err := c.db.StoreNetworkMetrics(metrics); err != nil {
		log.Printf("Error storing network metrics: %v", err)
	}
}

// convertToBytes converts a value with a unit to bytes
func convertToBytes(value float64, unit string) uint64 {
	unit = strings.ToUpper(unit)

	switch {
	case strings.HasPrefix(unit, "B"):
		return uint64(value)
	case strings.HasPrefix(unit, "KB"):
		return uint64(value * 1024)
	case strings.HasPrefix(unit, "MB"):
		return uint64(value * 1024 * 1024)
	case strings.HasPrefix(unit, "GB"):
		return uint64(value * 1024 * 1024 * 1024)
	case strings.HasPrefix(unit, "TB"):
		return uint64(value * 1024 * 1024 * 1024 * 1024)
	default:
		return uint64(value)
	}
}
