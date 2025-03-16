package devtools

import (
	"fmt"
	"sync"
	"time"
)

// DatabaseType represents the type of database
type DatabaseType string

const (
	MySQL      DatabaseType = "mysql"
	PostgreSQL DatabaseType = "postgresql"
	MongoDB    DatabaseType = "mongodb"
	SQLite     DatabaseType = "sqlite"
)

// DatabaseInfo represents information about a database connection
type DatabaseInfo struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Type        DatabaseType `json:"type"`
	Host        string       `json:"host"`
	Port        int          `json:"port"`
	Username    string       `json:"username"`
	Password    string       `json:"password,omitempty"`
	Database    string       `json:"database"`
	Status      string       `json:"status"`
	ConnectedAt string       `json:"connectedAt,omitempty"`
	URL         string       `json:"url,omitempty"`
	Description string       `json:"description,omitempty"`
}

// DatabaseManager manages database connections
type DatabaseManager struct {
	databases   map[string]*DatabaseInfo
	connections map[string]interface{} // This would be the actual database connection in a real implementation
	mutex       sync.Mutex
}

var (
	dbManager     *DatabaseManager
	dbManagerOnce sync.Once
)

// GetDatabaseManager returns the singleton instance of DatabaseManager
func GetDatabaseManager() *DatabaseManager {
	dbManagerOnce.Do(func() {
		dbManager = &DatabaseManager{
			databases:   make(map[string]*DatabaseInfo),
			connections: make(map[string]interface{}),
		}
		// Add some default databases for demonstration
		dbManager.addDefaultDatabases()
	})
	return dbManager
}

// addDefaultDatabases adds some default databases for demonstration
func (dm *DatabaseManager) addDefaultDatabases() {
	dm.databases["mysql-local"] = &DatabaseInfo{
		ID:          "mysql-local",
		Name:        "Local MySQL",
		Type:        MySQL,
		Host:        "localhost",
		Port:        3306,
		Username:    "root",
		Password:    "",
		Database:    "myapp",
		Status:      "disconnected",
		Description: "Local MySQL development database",
		URL:         "mysql://root@localhost:3306/myapp",
	}

	dm.databases["postgres-local"] = &DatabaseInfo{
		ID:          "postgres-local",
		Name:        "Local PostgreSQL",
		Type:        PostgreSQL,
		Host:        "localhost",
		Port:        5432,
		Username:    "postgres",
		Password:    "",
		Database:    "myapp",
		Status:      "disconnected",
		Description: "Local PostgreSQL development database",
		URL:         "postgresql://postgres@localhost:5432/myapp",
	}

	dm.databases["mongodb-local"] = &DatabaseInfo{
		ID:          "mongodb-local",
		Name:        "Local MongoDB",
		Type:        MongoDB,
		Host:        "localhost",
		Port:        27017,
		Username:    "",
		Password:    "",
		Database:    "myapp",
		Status:      "disconnected",
		Description: "Local MongoDB development database",
		URL:         "mongodb://localhost:27017/myapp",
	}
}

// GetAllDatabases returns all registered databases
func (dm *DatabaseManager) GetAllDatabases() []DatabaseInfo {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	databases := make([]DatabaseInfo, 0, len(dm.databases))
	for _, db := range dm.databases {
		// Don't include password in the response
		dbCopy := *db
		dbCopy.Password = ""
		databases = append(databases, dbCopy)
	}
	return databases
}

// ConnectDatabase connects to a database
func (dm *DatabaseManager) ConnectDatabase(databaseID string) (DatabaseInfo, error) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	db, exists := dm.databases[databaseID]
	if !exists {
		return DatabaseInfo{}, fmt.Errorf("database with ID %s not found", databaseID)
	}

	if db.Status == "connected" {
		return *db, fmt.Errorf("database is already connected")
	}

	// For demonstration, we'll simulate connecting to the database
	// In a real implementation, you would establish a connection to the database
	db.Status = "connected"
	db.ConnectedAt = time.Now().Format(time.RFC3339)

	// Update the database info
	dm.databases[databaseID] = db

	// Return a copy without the password
	dbCopy := *db
	dbCopy.Password = ""
	return dbCopy, nil
}

// DisconnectDatabase disconnects from a database
func (dm *DatabaseManager) DisconnectDatabase(databaseID string) (DatabaseInfo, error) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	db, exists := dm.databases[databaseID]
	if !exists {
		return DatabaseInfo{}, fmt.Errorf("database with ID %s not found", databaseID)
	}

	if db.Status != "connected" {
		return *db, fmt.Errorf("database is not connected")
	}

	// For demonstration, we'll simulate disconnecting from the database
	// In a real implementation, you would close the connection to the database
	db.Status = "disconnected"
	db.ConnectedAt = ""

	// Update the database info
	dm.databases[databaseID] = db

	// Return a copy without the password
	dbCopy := *db
	dbCopy.Password = ""
	return dbCopy, nil
}

// AddDatabase adds a new database configuration
func (dm *DatabaseManager) AddDatabase(db DatabaseInfo) (DatabaseInfo, error) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Generate a unique ID if not provided
	if db.ID == "" {
		db.ID = fmt.Sprintf("db-%d", time.Now().UnixNano())
	}

	// Check if database with the same ID already exists
	if _, exists := dm.databases[db.ID]; exists {
		return DatabaseInfo{}, fmt.Errorf("database with ID %s already exists", db.ID)
	}

	// Set default status
	db.Status = "disconnected"

	// Add the database
	dm.databases[db.ID] = &db

	// Return a copy without the password
	dbCopy := db
	dbCopy.Password = ""
	return dbCopy, nil
}

// RemoveDatabase removes a database configuration
func (dm *DatabaseManager) RemoveDatabase(databaseID string) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	db, exists := dm.databases[databaseID]
	if !exists {
		return fmt.Errorf("database with ID %s not found", databaseID)
	}

	// Check if the database is connected
	if db.Status == "connected" {
		return fmt.Errorf("cannot remove a connected database, disconnect it first")
	}

	// Remove the database
	delete(dm.databases, databaseID)

	return nil
}

// TestConnection tests a database connection without actually connecting
func (dm *DatabaseManager) TestConnection(db DatabaseInfo) (bool, string) {
	// For demonstration, we'll simulate testing the connection
	// In a real implementation, you would attempt to connect to the database

	// Simulate a successful connection for local databases
	if db.Host == "localhost" {
		return true, "Connection successful"
	}

	// Simulate a failed connection for other hosts
	return false, "Connection failed: could not connect to host"
}
