package devtools

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// GitRepoInfo represents information about a Git repository
type GitRepoInfo struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Branch       string    `json:"branch"`
	Status       string    `json:"status"`
	LastCommit   string    `json:"lastCommit"`
	LastCommitBy string    `json:"lastCommitBy"`
	LastUpdated  time.Time `json:"lastUpdated"`
	Changes      int       `json:"changes"`
	URL          string    `json:"url,omitempty"`
	Description  string    `json:"description,omitempty"`
}

// GitRepoManager manages Git repositories
type GitRepoManager struct {
	repos       map[string]*GitRepoInfo
	mutex       sync.Mutex
	configPath  string
	db          *sql.DB
	initialized bool
}

var (
	gitManager     *GitRepoManager
	gitManagerOnce sync.Once
)

// GetGitRepoManager returns the singleton instance of GitRepoManager
func GetGitRepoManager() *GitRepoManager {
	gitManagerOnce.Do(func() {
		// Create the .devex directory if it doesn't exist
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting user home directory: %v\n", err)
			home = "."
		}

		devexDir := filepath.Join(home, ".devex")
		if _, err := os.Stat(devexDir); os.IsNotExist(err) {
			if err := os.MkdirAll(devexDir, 0755); err != nil {
				fmt.Printf("Error creating .devex directory: %v\n", err)
			}
		}

		dbPath := filepath.Join(devexDir, "devex.db")
		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			fmt.Printf("Error opening database: %v\n", err)
		}

		gitManager = &GitRepoManager{
			repos:       make(map[string]*GitRepoInfo),
			configPath:  filepath.Join(devexDir, "gitrepos.json"),
			db:          db,
			initialized: false,
		}

		// Initialize the database
		gitManager.initDB()

		// Load repositories from the database
		gitManager.loadReposFromDB()
	})
	return gitManager
}

// initDB initializes the database schema
func (gm *GitRepoManager) initDB() {
	if gm.db == nil {
		return
	}

	// Create the git_repos table if it doesn't exist
	_, err := gm.db.Exec(`
		CREATE TABLE IF NOT EXISTS git_repos (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			path TEXT NOT NULL,
			branch TEXT,
			status TEXT,
			last_commit TEXT,
			last_commit_by TEXT,
			last_updated TEXT,
			changes INTEGER,
			url TEXT,
			description TEXT
		)
	`)
	if err != nil {
		fmt.Printf("Error creating git_repos table: %v\n", err)
	}
}

// loadReposFromDB loads repositories from the database
func (gm *GitRepoManager) loadReposFromDB() {
	if gm.db == nil {
		return
	}

	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	rows, err := gm.db.Query(`
		SELECT id, name, path, branch, status, last_commit, last_commit_by, last_updated, changes, url, description
		FROM git_repos
	`)
	if err != nil {
		fmt.Printf("Error querying git_repos: %v\n", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var repo GitRepoInfo
		var lastUpdatedStr string
		err := rows.Scan(
			&repo.ID,
			&repo.Name,
			&repo.Path,
			&repo.Branch,
			&repo.Status,
			&repo.LastCommit,
			&repo.LastCommitBy,
			&lastUpdatedStr,
			&repo.Changes,
			&repo.URL,
			&repo.Description,
		)
		if err != nil {
			fmt.Printf("Error scanning git_repo row: %v\n", err)
			continue
		}

		// Parse the last updated time
		if lastUpdatedStr != "" {
			lastUpdated, err := time.Parse(time.RFC3339, lastUpdatedStr)
			if err == nil {
				repo.LastUpdated = lastUpdated
			} else {
				repo.LastUpdated = time.Now()
			}
		} else {
			repo.LastUpdated = time.Now()
		}

		gm.repos[repo.ID] = &repo
		count++
	}

	// If no repositories were loaded, add default repositories
	if count == 0 {
		gm.addDefaultRepos()
	}

	gm.initialized = true
}

// saveRepoToDB saves a repository to the database
func (gm *GitRepoManager) saveRepoToDB(repo *GitRepoInfo) error {
	if gm.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Convert the last updated time to a string
	lastUpdatedStr := repo.LastUpdated.Format(time.RFC3339)

	// Insert or update the repository in the database
	_, err := gm.db.Exec(`
		INSERT OR REPLACE INTO git_repos (
			id, name, path, branch, status, last_commit, last_commit_by, last_updated, changes, url, description
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		repo.ID,
		repo.Name,
		repo.Path,
		repo.Branch,
		repo.Status,
		repo.LastCommit,
		repo.LastCommitBy,
		lastUpdatedStr,
		repo.Changes,
		repo.URL,
		repo.Description,
	)
	if err != nil {
		return fmt.Errorf("error saving repository to database: %v", err)
	}

	return nil
}

// deleteRepoFromDB deletes a repository from the database
func (gm *GitRepoManager) deleteRepoFromDB(repoID string) error {
	if gm.db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Delete the repository from the database
	_, err := gm.db.Exec("DELETE FROM git_repos WHERE id = ?", repoID)
	if err != nil {
		return fmt.Errorf("error deleting repository from database: %v", err)
	}

	return nil
}

// addDefaultRepos adds some default repositories for demonstration
func (gm *GitRepoManager) addDefaultRepos() {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "~/Documents/dev/app/DevEx" // Fallback to default path
	}

	// Add the current directory as a repository (assuming it's a Git repository)
	devexRepo := &GitRepoInfo{
		ID:           "devex",
		Name:         "DevEx",
		Path:         currentDir,
		Branch:       "main",
		Status:       "unknown", // Will be updated when refreshed
		LastCommit:   "Unknown",
		LastCommitBy: "Unknown",
		LastUpdated:  time.Now(),
		Changes:      0,
		URL:          "https://github.com/user/devex",
		Description:  "Developer Experience Tool",
	}
	gm.repos["devex"] = devexRepo
	gm.saveRepoToDB(devexRepo)

	// Add example repositories
	reactRepo := &GitRepoInfo{
		ID:           "react-app",
		Name:         "React App",
		Path:         "~/projects/react-app",
		Branch:       "develop",
		Status:       "modified",
		LastCommit:   "Update dependencies",
		LastCommitBy: "Jane Smith",
		LastUpdated:  time.Now().Add(-2 * 24 * time.Hour),
		Changes:      3,
		URL:          "https://github.com/user/react-app",
		Description:  "React.js frontend application",
	}
	gm.repos["react-app"] = reactRepo
	gm.saveRepoToDB(reactRepo)

	nodeRepo := &GitRepoInfo{
		ID:           "node-api",
		Name:         "Node API",
		Path:         "~/projects/node-api",
		Branch:       "feature/auth",
		Status:       "modified",
		LastCommit:   "Add authentication endpoints",
		LastCommitBy: "John Doe",
		LastUpdated:  time.Now().Add(-3 * 24 * time.Hour),
		Changes:      5,
		URL:          "https://github.com/user/node-api",
		Description:  "Node.js API server",
	}
	gm.repos["node-api"] = nodeRepo
	gm.saveRepoToDB(nodeRepo)
}

// GetAllRepos returns all registered repositories
func (gm *GitRepoManager) GetAllRepos() []GitRepoInfo {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	repos := make([]GitRepoInfo, 0, len(gm.repos))
	for _, repo := range gm.repos {
		repos = append(repos, *repo)
	}
	return repos
}

// RefreshRepo refreshes the status of a repository
func (gm *GitRepoManager) RefreshRepo(repoID string) (GitRepoInfo, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	repo, exists := gm.repos[repoID]
	if !exists {
		return GitRepoInfo{}, fmt.Errorf("repository with ID %s not found", repoID)
	}

	// Expand path if it contains ~
	path := repo.Path
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			repo.Status = "error"
			return *repo, fmt.Errorf("error getting home directory: %v", err)
		}
		path = filepath.Join(home, path[2:])
	}

	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		repo.Status = "not found"
		gm.saveRepoToDB(repo)
		return *repo, nil
	}

	// Check if it's a Git repository
	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		repo.Status = "not a git repository"
		gm.saveRepoToDB(repo)
		return *repo, nil
	}

	// Get the current branch
	cmd := exec.Command("git", "-C", path, "branch", "--show-current")
	output, err := cmd.CombinedOutput()
	if err == nil {
		branch := strings.TrimSpace(string(output))
		if branch != "" {
			repo.Branch = branch
		}
	}

	// Get the last commit
	cmd = exec.Command("git", "-C", path, "log", "-1", "--pretty=format:%s|%an|%ad", "--date=iso")
	output, err = cmd.CombinedOutput()
	if err == nil {
		parts := strings.Split(strings.TrimSpace(string(output)), "|")
		if len(parts) >= 3 {
			repo.LastCommit = parts[0]
			repo.LastCommitBy = parts[1]
			// Parse the date if possible
			if t, err := time.Parse("2006-01-02 15:04:05 -0700", parts[2]); err == nil {
				repo.LastUpdated = t
			}
		}
	}

	// Check if there are any changes
	cmd = exec.Command("git", "-C", path, "status", "--porcelain")
	output, err = cmd.CombinedOutput()
	if err == nil {
		changes := strings.TrimSpace(string(output))
		if changes == "" {
			repo.Status = "clean"
			repo.Changes = 0
		} else {
			repo.Status = "modified"
			repo.Changes = len(strings.Split(changes, "\n"))
		}
	} else {
		repo.Status = "error"
	}

	// Update the repository info
	gm.repos[repoID] = repo

	// Save the updated repository to the database
	gm.saveRepoToDB(repo)

	return *repo, nil
}

// RefreshAllRepos refreshes all repositories and returns the updated list
func (gm *GitRepoManager) RefreshAllRepos() []GitRepoInfo {
	gm.mutex.Lock()
	repoIDs := make([]string, 0, len(gm.repos))
	for id := range gm.repos {
		repoIDs = append(repoIDs, id)
	}
	gm.mutex.Unlock()

	// Refresh each repository
	for _, id := range repoIDs {
		// We don't need to handle errors here, as we want to continue refreshing other repos
		// even if one fails
		gm.RefreshRepo(id)
	}

	// Return the updated list
	return gm.GetAllRepos()
}

// AddRepo adds a new repository
func (gm *GitRepoManager) AddRepo(repo GitRepoInfo) (GitRepoInfo, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Generate a unique ID if not provided
	if repo.ID == "" {
		repo.ID = fmt.Sprintf("repo-%d", time.Now().UnixNano())
	}

	// Check if repository with the same ID already exists
	if _, exists := gm.repos[repo.ID]; exists {
		return GitRepoInfo{}, fmt.Errorf("repository with ID %s already exists", repo.ID)
	}

	// Set default values
	if repo.Status == "" {
		repo.Status = "unknown"
	}
	if repo.Branch == "" {
		repo.Branch = "unknown"
	}
	if repo.LastUpdated.IsZero() {
		repo.LastUpdated = time.Now()
	}

	// Add the repository
	gm.repos[repo.ID] = &repo

	// Save the repository to the database
	if err := gm.saveRepoToDB(&repo); err != nil {
		return GitRepoInfo{}, err
	}

	return repo, nil
}

// RemoveRepo removes a repository
func (gm *GitRepoManager) RemoveRepo(repoID string) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if _, exists := gm.repos[repoID]; !exists {
		return fmt.Errorf("repository with ID %s not found", repoID)
	}

	// Remove the repository
	delete(gm.repos, repoID)

	// Delete the repository from the database
	if err := gm.deleteRepoFromDB(repoID); err != nil {
		return err
	}

	return nil
}

// GetRepoChanges returns the changes in a Git repository
func (gm *GitRepoManager) GetRepoChanges(repoID string) ([]string, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	repo, exists := gm.repos[repoID]
	if !exists {
		return nil, fmt.Errorf("repository with ID %s not found", repoID)
	}

	// Expand path if it contains ~
	path := repo.Path
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting home directory: %v", err)
		}
		path = filepath.Join(home, path[2:])
	}

	// Check if the path exists and is a Git repository
	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a git repository or repository not found")
	}

	// Get the changes
	cmd := exec.Command("git", "-C", path, "status", "--porcelain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error getting git status: %v", err)
	}

	changes := strings.TrimSpace(string(output))
	if changes == "" {
		return []string{}, nil
	}

	return strings.Split(changes, "\n"), nil
}

// Close closes the database connection
func (gm *GitRepoManager) Close() error {
	if gm.db != nil {
		return gm.db.Close()
	}
	return nil
}
