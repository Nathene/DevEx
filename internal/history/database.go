package history

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// DB represents the SQLite database connection
type DB struct {
	conn *sql.DB
}

// NewDB creates a new database connection and initializes the schema
func NewDB() (*DB, error) {
	// Create data directory if it doesn't exist
	dataDir := getDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Connect to SQLite database
	dbPath := filepath.Join(dataDir, "devex_metrics.db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Initialize database schema
	if err := initSchema(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Setup automatic cleanup of old data
	go startCleanupRoutine(conn)

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// getDataDir returns the directory where the database file will be stored
func getDataDir() string {
	// Use user's home directory for data storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory can't be determined
		return ".devex"
	}
	return filepath.Join(homeDir, ".devex")
}

// initSchema creates the necessary tables if they don't exist
func initSchema(db *sql.DB) error {
	// Create CPU metrics table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS cpu_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			usage REAL NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create RAM metrics table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ram_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			usage REAL NOT NULL,
			used_bytes INTEGER NOT NULL,
			total_bytes INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create disk metrics table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS disk_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			usage REAL NOT NULL,
			used_bytes INTEGER NOT NULL,
			total_bytes INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create docker metrics table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS docker_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			daemon_running BOOLEAN NOT NULL,
			containers_running INTEGER NOT NULL,
			containers_total INTEGER NOT NULL,
			images_count INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create network metrics table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS network_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			internet_connected BOOLEAN NOT NULL,
			ping_latency REAL NOT NULL,
			dns_working BOOLEAN NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Create indexes for faster querying by timestamp
	for _, table := range []string{"cpu_metrics", "ram_metrics", "disk_metrics", "docker_metrics", "network_metrics"} {
		_, err = db.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_%s_timestamp ON %s(timestamp)", table, table))
		if err != nil {
			return err
		}
	}

	return nil
}

// startCleanupRoutine periodically removes data older than the retention period
func startCleanupRoutine(db *sql.DB) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cleanupOldData(db)
	}
}

// cleanupOldData removes data older than the retention period (1 hour by default)
func cleanupOldData(db *sql.DB) {
	// Calculate timestamp for retention cutoff (1 hour ago)
	cutoff := time.Now().Add(-1 * time.Hour).UTC().Format("2006-01-02 15:04:05")

	// Delete old data from all tables
	tables := []string{"cpu_metrics", "ram_metrics", "disk_metrics", "docker_metrics", "network_metrics"}
	for _, table := range tables {
		query := fmt.Sprintf("DELETE FROM %s WHERE timestamp < ?", table)
		result, err := db.Exec(query, cutoff)
		if err != nil {
			log.Printf("Error cleaning up old data from %s: %v", table, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			log.Printf("Cleaned up %d old records from %s", rowsAffected, table)
		}
	}
}
