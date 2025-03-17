package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"DevEx/internal/devtools"
	"DevEx/internal/docker"
	"DevEx/internal/history"
	"DevEx/internal/network"
	"DevEx/internal/process"
	"DevEx/internal/system"
)

// App struct
type App struct {
	ctx             context.Context
	db              *history.DB
	collector       *history.Collector
	processManager  *process.Manager
	devToolsManager *devtools.DevToolsManager
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
		db:              db,
		collector:       collector,
		processManager:  processManager,
		devToolsManager: devtools.NewDevToolsManager(),
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

	// Initialize DevTools manager
	if err := a.devToolsManager.Initialize(); err != nil {
		log.Printf("Error initializing DevTools manager: %v", err)
	}
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// Stop the metrics collector
	a.collector.Stop()

	// Stop the process manager
	a.processManager.Stop()

	// Close the Git repository manager
	if gitManager := a.devToolsManager.GetGitRepoManager(); gitManager != nil {
		if err := gitManager.Close(); err != nil {
			log.Printf("Error closing Git repository manager: %v", err)
		}
	}

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

// GetTopMemoryProcesses returns the top processes by memory usage
func (a *App) GetTopMemoryProcesses() []process.ProcessWithPorts {
	processes := a.processManager.GetProcesses()

	// Sort processes by memory usage (highest first)
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Process.MemoryUsage > processes[j].Process.MemoryUsage
	})

	// Return top 30 or all if less than 30
	if len(processes) > 30 {
		return processes[:30]
	}
	return processes
}

// GetTopCPUProcesses returns the top processes by CPU usage
func (a *App) GetTopCPUProcesses() []process.ProcessWithPorts {
	processes := a.processManager.GetProcesses()

	// Sort processes by CPU usage (highest first)
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Process.CPUPercent > processes[j].Process.CPUPercent
	})

	// Return top 30 or all if less than 30
	if len(processes) > 30 {
		return processes[:30]
	}
	return processes
}

// GetTopDiskProcesses returns the top processes by disk I/O (approximated by process age)
// Note: Getting actual disk I/O per process requires additional monitoring
func (a *App) GetTopDiskProcesses() []process.ProcessWithPorts {
	processes := a.processManager.GetProcesses()

	// Sort processes by memory usage as a proxy for potential disk usage
	// In a real implementation, you would track actual disk I/O per process
	sort.Slice(processes, func(i, j int) bool {
		return processes[i].Process.MemoryUsage > processes[j].Process.MemoryUsage
	})

	// Return top 30 or all if less than 30
	if len(processes) > 30 {
		return processes[:30]
	}
	return processes
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

// DevTools methods

// GetAllServers returns all registered servers
func (a *App) GetAllServers() []devtools.ServerInfo {
	return a.devToolsManager.GetAllServers()
}

// StartServer starts a development server
func (a *App) StartServer(serverID string) (devtools.ServerInfo, error) {
	return a.devToolsManager.StartServer(serverID)
}

// StopServer stops a development server
func (a *App) StopServer(serverID string) (devtools.ServerInfo, error) {
	return a.devToolsManager.StopServer(serverID)
}

// AddServer adds a new server configuration
func (a *App) AddServer(server devtools.ServerInfo) (devtools.ServerInfo, error) {
	return a.devToolsManager.AddServer(server)
}

// RemoveServer removes a server configuration
func (a *App) RemoveServer(serverID string) error {
	return a.devToolsManager.RemoveServer(serverID)
}

// GetAllDatabases returns all registered databases
func (a *App) GetAllDatabases() []devtools.DatabaseInfo {
	return a.devToolsManager.GetAllDatabases()
}

// ConnectDatabase connects to a database
func (a *App) ConnectDatabase(databaseID string) (devtools.DatabaseInfo, error) {
	return a.devToolsManager.ConnectDatabase(databaseID)
}

// DisconnectDatabase disconnects from a database
func (a *App) DisconnectDatabase(databaseID string) (devtools.DatabaseInfo, error) {
	return a.devToolsManager.DisconnectDatabase(databaseID)
}

// AddDatabase adds a new database configuration
func (a *App) AddDatabase(db devtools.DatabaseInfo) (devtools.DatabaseInfo, error) {
	return a.devToolsManager.AddDatabase(db)
}

// RemoveDatabase removes a database configuration
func (a *App) RemoveDatabase(databaseID string) error {
	return a.devToolsManager.RemoveDatabase(databaseID)
}

// TestDatabaseConnection tests a database connection without actually connecting
func (a *App) TestDatabaseConnection(db devtools.DatabaseInfo) (bool, string) {
	return a.devToolsManager.TestDatabaseConnection(db)
}

// SendAPIRequest sends an API request and returns the response
func (a *App) SendAPIRequest(req devtools.APIRequest) devtools.APIResponse {
	return a.devToolsManager.SendAPIRequest(req)
}

// GetSavedAPIRequests returns a list of saved API requests
func (a *App) GetSavedAPIRequests() []devtools.APIRequest {
	return a.devToolsManager.GetSavedAPIRequests()
}

// GetAllGitRepos returns all registered Git repositories
func (a *App) GetAllGitRepos() []devtools.GitRepoInfo {
	return a.devToolsManager.GetAllGitRepos()
}

// RefreshGitRepo refreshes the status of a Git repository
func (a *App) RefreshGitRepo(repoID string) (devtools.GitRepoInfo, error) {
	return a.devToolsManager.RefreshGitRepo(repoID)
}

// RefreshAllGitRepos refreshes all Git repositories
func (a *App) RefreshAllGitRepos() []devtools.GitRepoInfo {
	return a.devToolsManager.RefreshAllGitRepos()
}

// AddGitRepo adds a new Git repository
func (a *App) AddGitRepo(repo devtools.GitRepoInfo) (devtools.GitRepoInfo, error) {
	return a.devToolsManager.AddGitRepo(repo)
}

// RemoveGitRepo removes a Git repository
func (a *App) RemoveGitRepo(repoID string) error {
	return a.devToolsManager.RemoveGitRepo(repoID)
}

// GetGitRepoChanges returns the changes in a Git repository
func (a *App) GetGitRepoChanges(repoID string) ([]string, error) {
	return a.devToolsManager.GetGitRepoChanges(repoID)
}

// OpenInVSCode opens a directory in Visual Studio Code
func (a *App) OpenInVSCode(path string) error {
	// Expand path if it contains ~
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %v", err)
		}
		path = filepath.Join(home, path[2:])
	}

	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("directory not found: %s", path)
	}

	// Open the directory in VS Code
	cmd := exec.Command("code", "-n", path)
	return cmd.Start()
}

// OpenFolderPicker opens a native folder picker dialog and returns the selected path
func (a *App) OpenFolderPicker() (string, error) {
	// Use the Wails runtime to open a directory selection dialog
	selectedDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Repository Directory",
	})

	if err != nil {
		return "", fmt.Errorf("error opening folder picker: %v", err)
	}

	// If user canceled, return empty string without error
	if selectedDir == "" {
		return "", nil
	}

	return selectedDir, nil
}
