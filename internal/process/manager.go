package process

import (
	"context"
	"log"
	"sync"
	"time"
)

// Manager handles process monitoring and caching
type Manager struct {
	processes      []ProcessWithPorts
	lastUpdateTime time.Time
	updateInterval time.Duration
	ctx            context.Context
	cancelFunc     context.CancelFunc
	mutex          sync.RWMutex
	isUpdating     bool
	maxProcesses   int
}

// NewManager creates a new process manager
func NewManager(updateInterval time.Duration) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		updateInterval: updateInterval,
		ctx:            ctx,
		cancelFunc:     cancel,
		maxProcesses:   500, // Limit to 500 processes by default
	}
}

// Start begins the process monitoring
func (m *Manager) Start() {
	log.Printf("Starting process manager with interval: %v", m.updateInterval)

	// Update immediately
	m.updateProcesses()

	// Then update at regular intervals
	ticker := time.NewTicker(m.updateInterval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.updateProcesses()
			case <-m.ctx.Done():
				log.Println("Process manager stopped")
				return
			}
		}
	}()
}

// Stop halts the process monitoring
func (m *Manager) Stop() {
	m.cancelFunc()
}

// updateProcesses refreshes the process list
func (m *Manager) updateProcesses() {
	// Prevent concurrent updates
	if m.isUpdating {
		log.Println("Process update already in progress, skipping")
		return
	}

	m.mutex.Lock()
	m.isUpdating = true
	m.mutex.Unlock()

	defer func() {
		m.mutex.Lock()
		m.isUpdating = false
		m.mutex.Unlock()
	}()

	processes, err := GetProcessesWithPorts()
	if err != nil {
		log.Printf("Error updating processes: %v", err)
		return
	}

	// Limit the number of processes to avoid overwhelming the UI
	if len(processes) > m.maxProcesses {
		processes = processes[:m.maxProcesses]
	}

	m.mutex.Lock()
	m.processes = processes
	m.lastUpdateTime = time.Now()
	m.mutex.Unlock()
}

// GetProcesses returns the cached process list
func (m *Manager) GetProcesses() []ProcessWithPorts {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.processes
}

// SearchByPort searches for processes using a specific port
func (m *Manager) SearchByPort(port int) []ProcessWithPorts {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var result []ProcessWithPorts
	for _, p := range m.processes {
		for _, portInfo := range p.Ports {
			if portInfo.Port == port {
				result = append(result, p)
				break
			}
		}
	}
	return result
}

// KillProcessByPID terminates a process
func (m *Manager) KillProcessByPID(pid int) error {
	err := KillProcess(pid)
	if err != nil {
		return err
	}

	// Update the process list after killing
	go m.updateProcesses() // Run in background to avoid blocking
	return nil
}

// SetMaxProcesses sets the maximum number of processes to return
func (m *Manager) SetMaxProcesses(max int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.maxProcesses = max
}
