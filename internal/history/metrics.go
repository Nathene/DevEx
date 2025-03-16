package history

import (
	"time"
)

// TimeSeriesPoint represents a single data point in a time series
type TimeSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// CPUMetrics represents CPU usage metrics
type CPUMetrics struct {
	Timestamp time.Time `json:"timestamp"`
	Usage     float64   `json:"usage"`
}

// RAMMetrics represents RAM usage metrics
type RAMMetrics struct {
	Timestamp  time.Time `json:"timestamp"`
	Usage      float64   `json:"usage"`
	UsedBytes  uint64    `json:"usedBytes"`
	TotalBytes uint64    `json:"totalBytes"`
}

// DiskMetrics represents disk usage metrics
type DiskMetrics struct {
	Timestamp  time.Time `json:"timestamp"`
	Usage      float64   `json:"usage"`
	UsedBytes  uint64    `json:"usedBytes"`
	TotalBytes uint64    `json:"totalBytes"`
}

// DockerMetrics represents Docker metrics
type DockerMetrics struct {
	Timestamp         time.Time `json:"timestamp"`
	DaemonRunning     bool      `json:"daemonRunning"`
	ContainersRunning int       `json:"containersRunning"`
	ContainersTotal   int       `json:"containersTotal"`
	ImagesCount       int       `json:"imagesCount"`
}

// NetworkMetrics represents network metrics
type NetworkMetrics struct {
	Timestamp         time.Time `json:"timestamp"`
	InternetConnected bool      `json:"internetConnected"`
	PingLatency       float64   `json:"pingLatency"`
	DNSWorking        bool      `json:"dnsWorking"`
}

// StoreCPUMetrics stores CPU metrics in the database
func (db *DB) StoreCPUMetrics(metrics CPUMetrics) error {
	query := `INSERT INTO cpu_metrics (timestamp, usage) VALUES (?, ?)`
	_, err := db.conn.Exec(query, metrics.Timestamp.UTC(), metrics.Usage)
	return err
}

// StoreRAMMetrics stores RAM metrics in the database
func (db *DB) StoreRAMMetrics(metrics RAMMetrics) error {
	query := `INSERT INTO ram_metrics (timestamp, usage, used_bytes, total_bytes) VALUES (?, ?, ?, ?)`
	_, err := db.conn.Exec(query, metrics.Timestamp.UTC(), metrics.Usage, metrics.UsedBytes, metrics.TotalBytes)
	return err
}

// StoreDiskMetrics stores disk metrics in the database
func (db *DB) StoreDiskMetrics(metrics DiskMetrics) error {
	query := `INSERT INTO disk_metrics (timestamp, usage, used_bytes, total_bytes) VALUES (?, ?, ?, ?)`
	_, err := db.conn.Exec(query, metrics.Timestamp.UTC(), metrics.Usage, metrics.UsedBytes, metrics.TotalBytes)
	return err
}

// StoreDockerMetrics stores Docker metrics in the database
func (db *DB) StoreDockerMetrics(metrics DockerMetrics) error {
	query := `INSERT INTO docker_metrics (timestamp, daemon_running, containers_running, containers_total, images_count) 
			  VALUES (?, ?, ?, ?, ?)`
	_, err := db.conn.Exec(query, metrics.Timestamp.UTC(), metrics.DaemonRunning, metrics.ContainersRunning,
		metrics.ContainersTotal, metrics.ImagesCount)
	return err
}

// StoreNetworkMetrics stores network metrics in the database
func (db *DB) StoreNetworkMetrics(metrics NetworkMetrics) error {
	query := `INSERT INTO network_metrics (timestamp, internet_connected, ping_latency, dns_working) VALUES (?, ?, ?, ?)`
	_, err := db.conn.Exec(query, metrics.Timestamp.UTC(), metrics.InternetConnected, metrics.PingLatency, metrics.DNSWorking)
	return err
}

// GetCPUHistory retrieves CPU usage history for the specified duration
func (db *DB) GetCPUHistory(duration time.Duration) ([]TimeSeriesPoint, error) {
	cutoff := time.Now().Add(-duration).UTC()
	query := `SELECT timestamp, usage FROM cpu_metrics WHERE timestamp > ? ORDER BY timestamp ASC`

	rows, err := db.conn.Query(query, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TimeSeriesPoint
	for rows.Next() {
		var timestamp time.Time
		var usage float64
		if err := rows.Scan(&timestamp, &usage); err != nil {
			return nil, err
		}
		result = append(result, TimeSeriesPoint{
			Timestamp: timestamp,
			Value:     usage,
		})
	}

	return result, rows.Err()
}

// GetRAMHistory retrieves RAM usage history for the specified duration
func (db *DB) GetRAMHistory(duration time.Duration) ([]TimeSeriesPoint, error) {
	cutoff := time.Now().Add(-duration).UTC()
	query := `SELECT timestamp, usage FROM ram_metrics WHERE timestamp > ? ORDER BY timestamp ASC`

	rows, err := db.conn.Query(query, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TimeSeriesPoint
	for rows.Next() {
		var timestamp time.Time
		var usage float64
		if err := rows.Scan(&timestamp, &usage); err != nil {
			return nil, err
		}
		result = append(result, TimeSeriesPoint{
			Timestamp: timestamp,
			Value:     usage,
		})
	}

	return result, rows.Err()
}

// GetDiskHistory retrieves disk usage history for the specified duration
func (db *DB) GetDiskHistory(duration time.Duration) ([]TimeSeriesPoint, error) {
	cutoff := time.Now().Add(-duration).UTC()
	query := `SELECT timestamp, usage FROM disk_metrics WHERE timestamp > ? ORDER BY timestamp ASC`

	rows, err := db.conn.Query(query, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TimeSeriesPoint
	for rows.Next() {
		var timestamp time.Time
		var usage float64
		if err := rows.Scan(&timestamp, &usage); err != nil {
			return nil, err
		}
		result = append(result, TimeSeriesPoint{
			Timestamp: timestamp,
			Value:     usage,
		})
	}

	return result, rows.Err()
}

// GetDockerHistory retrieves Docker metrics history for the specified duration
func (db *DB) GetDockerHistory(duration time.Duration) ([]DockerMetrics, error) {
	cutoff := time.Now().Add(-duration).UTC()
	query := `SELECT timestamp, daemon_running, containers_running, containers_total, images_count 
			  FROM docker_metrics WHERE timestamp > ? ORDER BY timestamp ASC`

	rows, err := db.conn.Query(query, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []DockerMetrics
	for rows.Next() {
		var metrics DockerMetrics
		if err := rows.Scan(&metrics.Timestamp, &metrics.DaemonRunning, &metrics.ContainersRunning,
			&metrics.ContainersTotal, &metrics.ImagesCount); err != nil {
			return nil, err
		}
		result = append(result, metrics)
	}

	return result, rows.Err()
}

// GetNetworkHistory retrieves network metrics history for the specified duration
func (db *DB) GetNetworkHistory(duration time.Duration) ([]NetworkMetrics, error) {
	cutoff := time.Now().Add(-duration).UTC()
	query := `SELECT timestamp, internet_connected, ping_latency, dns_working 
			  FROM network_metrics WHERE timestamp > ? ORDER BY timestamp ASC`

	rows, err := db.conn.Query(query, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []NetworkMetrics
	for rows.Next() {
		var metrics NetworkMetrics
		if err := rows.Scan(&metrics.Timestamp, &metrics.InternetConnected, &metrics.PingLatency, &metrics.DNSWorking); err != nil {
			return nil, err
		}
		result = append(result, metrics)
	}

	return result, rows.Err()
}
