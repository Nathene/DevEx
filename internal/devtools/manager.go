package devtools

import (
	"fmt"
	"os"
	"path/filepath"
)

// DevToolsManager is the main manager for developer tools
type DevToolsManager struct {
	serverManager   *ServerManager
	databaseManager *DatabaseManager
	apiTester       *APITester
	gitRepoManager  *GitRepoManager
	initialized     bool
}

// NewDevToolsManager creates a new DevToolsManager
func NewDevToolsManager() *DevToolsManager {
	return &DevToolsManager{
		serverManager:   GetServerManager(),
		databaseManager: GetDatabaseManager(),
		apiTester:       NewAPITester(),
		gitRepoManager:  GetGitRepoManager(),
	}
}

// Initialize initializes the DevToolsManager
func (dtm *DevToolsManager) Initialize() error {
	if dtm.initialized {
		return nil
	}

	// Create the .devex directory if it doesn't exist
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %v", err)
	}

	devexDir := filepath.Join(home, ".devex")
	if _, err := os.Stat(devexDir); os.IsNotExist(err) {
		if err := os.MkdirAll(devexDir, 0755); err != nil {
			return fmt.Errorf("error creating .devex directory: %v", err)
		}
	}

	dtm.initialized = true
	return nil
}

// GetAllServers returns all registered servers
func (dtm *DevToolsManager) GetAllServers() []ServerInfo {
	return dtm.serverManager.GetAllServers()
}

// StartServer starts a development server
func (dtm *DevToolsManager) StartServer(serverID string) (ServerInfo, error) {
	return dtm.serverManager.StartServer(serverID)
}

// StopServer stops a development server
func (dtm *DevToolsManager) StopServer(serverID string) (ServerInfo, error) {
	return dtm.serverManager.StopServer(serverID)
}

// AddServer adds a new server configuration
func (dtm *DevToolsManager) AddServer(server ServerInfo) (ServerInfo, error) {
	return dtm.serverManager.AddServer(server)
}

// RemoveServer removes a server configuration
func (dtm *DevToolsManager) RemoveServer(serverID string) error {
	return dtm.serverManager.RemoveServer(serverID)
}

// GetAllDatabases returns all registered databases
func (dtm *DevToolsManager) GetAllDatabases() []DatabaseInfo {
	return dtm.databaseManager.GetAllDatabases()
}

// ConnectDatabase connects to a database
func (dtm *DevToolsManager) ConnectDatabase(databaseID string) (DatabaseInfo, error) {
	return dtm.databaseManager.ConnectDatabase(databaseID)
}

// DisconnectDatabase disconnects from a database
func (dtm *DevToolsManager) DisconnectDatabase(databaseID string) (DatabaseInfo, error) {
	return dtm.databaseManager.DisconnectDatabase(databaseID)
}

// AddDatabase adds a new database configuration
func (dtm *DevToolsManager) AddDatabase(db DatabaseInfo) (DatabaseInfo, error) {
	return dtm.databaseManager.AddDatabase(db)
}

// RemoveDatabase removes a database configuration
func (dtm *DevToolsManager) RemoveDatabase(databaseID string) error {
	return dtm.databaseManager.RemoveDatabase(databaseID)
}

// TestDatabaseConnection tests a database connection without actually connecting
func (dtm *DevToolsManager) TestDatabaseConnection(db DatabaseInfo) (bool, string) {
	return dtm.databaseManager.TestConnection(db)
}

// SendAPIRequest sends an API request and returns the response
func (dtm *DevToolsManager) SendAPIRequest(req APIRequest) APIResponse {
	return dtm.apiTester.SendRequest(req)
}

// GetSavedAPIRequests returns a list of saved API requests
func (dtm *DevToolsManager) GetSavedAPIRequests() []APIRequest {
	return dtm.apiTester.GetSavedRequests()
}

// GetAllGitRepos returns all registered Git repositories
func (dtm *DevToolsManager) GetAllGitRepos() []GitRepoInfo {
	return dtm.gitRepoManager.GetAllRepos()
}

// RefreshGitRepo refreshes the status of a Git repository
func (dtm *DevToolsManager) RefreshGitRepo(repoID string) (GitRepoInfo, error) {
	return dtm.gitRepoManager.RefreshRepo(repoID)
}

// AddGitRepo adds a new Git repository
func (dtm *DevToolsManager) AddGitRepo(repo GitRepoInfo) (GitRepoInfo, error) {
	return dtm.gitRepoManager.AddRepo(repo)
}

// RemoveGitRepo removes a Git repository
func (dtm *DevToolsManager) RemoveGitRepo(repoID string) error {
	return dtm.gitRepoManager.RemoveRepo(repoID)
}

// GetGitRepoChanges returns the changes in a Git repository
func (dtm *DevToolsManager) GetGitRepoChanges(repoID string) ([]string, error) {
	return dtm.gitRepoManager.GetRepoChanges(repoID)
}

// GetGitRepoManager returns the Git repository manager
func (dtm *DevToolsManager) GetGitRepoManager() *GitRepoManager {
	return dtm.gitRepoManager
}
