package devtools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ServerInfo represents information about a development server
type ServerInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Port        int    `json:"port"`
	Path        string `json:"path"`
	Command     string `json:"command"`
	Status      string `json:"status"`
	PID         int    `json:"pid"`
	StartTime   string `json:"startTime,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

// ServerManager manages development servers
type ServerManager struct {
	servers     map[string]*ServerInfo
	processes   map[string]*exec.Cmd
	mutex       sync.Mutex
	configPath  string
	initialized bool
}

var (
	manager *ServerManager
	once    sync.Once
)

// GetServerManager returns the singleton instance of ServerManager
func GetServerManager() *ServerManager {
	once.Do(func() {
		manager = &ServerManager{
			servers:    make(map[string]*ServerInfo),
			processes:  make(map[string]*exec.Cmd),
			configPath: filepath.Join(os.Getenv("HOME"), ".devex", "servers.json"),
		}
		// Add some default servers for demonstration
		manager.addDefaultServers()
	})
	return manager
}

// addDefaultServers adds some default servers for demonstration
func (sm *ServerManager) addDefaultServers() {
	sm.servers["node-server"] = &ServerInfo{
		ID:          "node-server",
		Name:        "Node.js API Server",
		Type:        "nodejs",
		Port:        3000,
		Path:        "~/projects/node-api",
		Command:     "npm start",
		Status:      "stopped",
		Description: "RESTful API server built with Express.js",
		URL:         "http://localhost:3000",
	}

	sm.servers["react-app"] = &ServerInfo{
		ID:          "react-app",
		Name:        "React Frontend",
		Type:        "react",
		Port:        3001,
		Path:        "~/projects/react-app",
		Command:     "npm start",
		Status:      "stopped",
		Description: "React.js frontend application",
		URL:         "http://localhost:3001",
	}

	sm.servers["python-flask"] = &ServerInfo{
		ID:          "python-flask",
		Name:        "Python Flask API",
		Type:        "python",
		Port:        5000,
		Path:        "~/projects/flask-api",
		Command:     "flask run",
		Status:      "stopped",
		Description: "Python Flask API server",
		URL:         "http://localhost:5000",
	}
}

// GetAllServers returns all registered servers
func (sm *ServerManager) GetAllServers() []ServerInfo {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	servers := make([]ServerInfo, 0, len(sm.servers))
	for _, server := range sm.servers {
		servers = append(servers, *server)
	}
	return servers
}

// StartServer starts a development server
func (sm *ServerManager) StartServer(serverID string) (ServerInfo, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	server, exists := sm.servers[serverID]
	if !exists {
		return ServerInfo{}, fmt.Errorf("server with ID %s not found", serverID)
	}

	if server.Status == "running" {
		return *server, fmt.Errorf("server is already running")
	}

	// Expand path if it contains ~
	path := server.Path
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}

	// For demonstration, we'll simulate starting the server
	// In a real implementation, you would execute the command
	server.Status = "running"
	server.PID = 12345 // This would be the actual PID in a real implementation
	server.StartTime = time.Now().Format(time.RFC3339)

	// Update the server info
	sm.servers[serverID] = server

	return *server, nil
}

// StopServer stops a development server
func (sm *ServerManager) StopServer(serverID string) (ServerInfo, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	server, exists := sm.servers[serverID]
	if !exists {
		return ServerInfo{}, fmt.Errorf("server with ID %s not found", serverID)
	}

	if server.Status != "running" {
		return *server, fmt.Errorf("server is not running")
	}

	// For demonstration, we'll simulate stopping the server
	// In a real implementation, you would kill the process
	server.Status = "stopped"
	server.PID = 0
	server.StartTime = ""

	// Update the server info
	sm.servers[serverID] = server

	return *server, nil
}

// AddServer adds a new server configuration
func (sm *ServerManager) AddServer(server ServerInfo) (ServerInfo, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Generate a unique ID if not provided
	if server.ID == "" {
		server.ID = fmt.Sprintf("server-%d", time.Now().UnixNano())
	}

	// Check if server with the same ID already exists
	if _, exists := sm.servers[server.ID]; exists {
		return ServerInfo{}, fmt.Errorf("server with ID %s already exists", server.ID)
	}

	// Set default status
	server.Status = "stopped"

	// Add the server
	sm.servers[server.ID] = &server

	return server, nil
}

// RemoveServer removes a server configuration
func (sm *ServerManager) RemoveServer(serverID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	server, exists := sm.servers[serverID]
	if !exists {
		return fmt.Errorf("server with ID %s not found", serverID)
	}

	// Check if the server is running
	if server.Status == "running" {
		return fmt.Errorf("cannot remove a running server, stop it first")
	}

	// Remove the server
	delete(sm.servers, serverID)

	return nil
}

// CheckPort checks if a port is in use
func (sm *ServerManager) CheckPort(port int) bool {
	// This is a simplified implementation
	// In a real implementation, you would check if the port is in use
	cmd := exec.Command("lsof", "-i", ":"+strconv.Itoa(port))
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If there's an error, the port is likely not in use
		return false
	}

	// If the output contains data, the port is in use
	return len(output) > 0
}
