<script>
  import { GetCPUInfo, GetCPUDetails, GetRAMInfo, GetRAMDetails, GetDiskInfo, GetDiskDetails, GetDockerStatus, GetDockerMetrics, GetNetworkStatus, GetCPUHistory, GetRAMHistory, GetDiskHistory, GetAllProcesses, SearchProcessesByPort, KillProcess, FormatProcessBytes } from '../wailsjs/go/main/App';
  import { BrowserOpenURL } from '../wailsjs/runtime/runtime';
  import { onMount } from 'svelte';

  let cpuInfo = "Loading...";
  let cpuDetails = "Loading...";
  let ramInfo = "Loading...";
  let ramDetails = "Loading...";
  let diskInfo = "Loading...";
  let diskDetails = "Loading...";

  let dockerStatus = {
    daemonRunning: false,
    version: "Loading...",
    info: "Loading..."
  };
  let dockerMetrics = {
    imagesCount: 0,
    containersAll: 0,
    containersUp: 0,
    diskUsage: "Loading...",
    networkStatus: "Loading..."
  };

  let networkStatus = {
    internetConnected: false,
    pingLatency: 0,
    pingStatus: "Loading...",
    dnsStatus: "Loading..."
  };

  let selectedCategory = 'dashboard';

  const categories = [
    { id: 'dashboard', name: 'Dashboard', icon: '📊' },
    { id: 'history', name: 'History', icon: '📈' },
    { id: 'processes', name: 'Processes', icon: '⚙️' },
    { id: 'logs', name: 'Logs', icon: '📝' },
    { id: 'settings', name: 'Settings', icon: '⚙️' }
  ];

  // Add history data state
  let cpuHistory = [];
  let ramHistory = [];
  let diskHistory = [];
  let selectedTimeRange = 15; // Default to 15 minutes
  let historyError = ""; // Add error state for debugging

  // Process management state
  let allProcesses = [];
  let filteredProcesses = [];
  let processSearchQuery = "";
  let isSearching = false;
  let showKillConfirmation = false;
  let processToKill = { pid: 0, name: "" };
  let processError = "";
  let isLoadingProcesses = false;
  let searchTimeout = null;
  let filterOption = "all"; // Changed from individual booleans to a single option

  function getStatusClass(value) {
    if (value === "Loading...") return "pending";
    const percentage = parseFloat(value.match(/\d+\.?\d*/)[0]);
    if (percentage >= 90) return "error";
    if (percentage >= 80) return "warning";
    return "success";
  }

  function getStatusText(value) {
    if (value === "Loading...") return "Checking...";
    const percentage = parseFloat(value.match(/\d+\.?\d*/)[0]);
    if (percentage >= 90) return "Critical";
    if (percentage >= 80) return "Warning";
    return "Healthy";
  }

  async function updateMetrics() {
    try {
      cpuInfo = await GetCPUInfo();
      cpuDetails = await GetCPUDetails();
      ramInfo = await GetRAMInfo();
      ramDetails = await GetRAMDetails();
      diskInfo = await GetDiskInfo();
      diskDetails = await GetDiskDetails();
      dockerStatus = await GetDockerStatus();
      dockerMetrics = await GetDockerMetrics();
      networkStatus = await GetNetworkStatus();
    } catch (error) {
      console.error('Error updating metrics:', error);
    }
  }

  async function openPprofDashboard() {
    await BrowserOpenURL('http://localhost:6060/debug/pprof/');
  }

  // Function to update history data
  async function updateHistoryData() {
    try {
      historyError = ""; // Reset error state
      console.log("Fetching history data for", selectedTimeRange, "minutes");
      
      const cpuData = await GetCPUHistory(selectedTimeRange);
      const ramData = await GetRAMHistory(selectedTimeRange);
      const diskData = await GetDiskHistory(selectedTimeRange);
      
      console.log("CPU History:", cpuData);
      console.log("RAM History:", ramData);
      console.log("Disk History:", diskData);
      
      cpuHistory = cpuData || [];
      ramHistory = ramData || [];
      diskHistory = diskData || [];
      
      if (cpuHistory.length === 0 && ramHistory.length === 0 && diskHistory.length === 0) {
        historyError = "No historical data available yet. Data is being collected every 10 seconds.";
      }
    } catch (error) {
      console.error('Error updating history data:', error);
      historyError = `Error fetching history data: ${error.message}`;
    }
  }

  // Function to format timestamp for display
  function formatTimestamp(timestamp) {
    const date = new Date(timestamp);
    return date.toLocaleTimeString();
  }

  // Function to handle time range change
  function handleTimeRangeChange(event) {
    selectedTimeRange = parseInt(event.target.value);
    updateHistoryData();
  }

  // Fetch all processes
  async function updateProcesses() {
    // Don't update if we're already loading
    if (isLoadingProcesses) {
      return;
    }
    
    try {
      isLoadingProcesses = true;
      processError = "";
      
      console.log("Fetching processes...");
      const processes = await GetAllProcesses();
      console.log(`Received ${processes.length} processes`);
      
      allProcesses = processes;
      
      // Apply filters (this now handles search and other filters)
      applyFilters();
    } catch (error) {
      console.error("Error fetching processes:", error);
      processError = `Error fetching processes: ${error.message}`;
      filteredProcesses = [];
    } finally {
      isLoadingProcesses = false;
    }
  }

  // Apply filters (based on selected filter option)
  function applyFilters() {
    let result = [...allProcesses]; // Start with all processes instead of filtered
    
    // Apply search if there's a query
    if (processSearchQuery.trim() !== "") {
      // Check if query is a number (potential port search)
      if (/^\d+$/.test(processSearchQuery)) {
        // Match any part of the port number, not just the beginning
        result = result.filter(p => {
          // Check if any of the process's ports contain the query
          if (p.ports && p.ports.length > 0) {
            return p.ports.some(port => 
              port.port.toString().includes(processSearchQuery)
            );
          }
          return false;
        });
        
        // Mark matching ports for highlighting
        if (result.length > 0) {
          result = result.map(p => {
            const matchingPorts = p.ports.map(port => ({
              ...port,
              isMatch: port.port.toString().includes(processSearchQuery)
            }));
            
            return {
              ...p,
              ports: matchingPorts
            };
          });
        }
      } else {
        // Search by process name or command line
        result = result.filter(p => 
          p.process.name.toLowerCase().includes(processSearchQuery.toLowerCase()) ||
          p.process.commandLine.toLowerCase().includes(processSearchQuery.toLowerCase())
        );
      }
    }
    
    // Apply filter based on selected option
    if (filterOption === "withPorts") {
      result = result.filter(p => p.ports && p.ports.length > 0);
    } else if (filterOption === "sortByCPU") {
      result.sort((a, b) => b.process.cpuPercent - a.process.cpuPercent);
    } else if (filterOption === "withPortsAndSortByCPU") {
      result = result.filter(p => p.ports && p.ports.length > 0);
      result.sort((a, b) => b.process.cpuPercent - a.process.cpuPercent);
    }
    
    filteredProcesses = result;
    console.log(`Applied filter: ${filterOption}, showing ${result.length} processes`);
  }

  // Handle filter option change
  function handleFilterChange(option) {
    filterOption = option;
    applyFilters();
  }

  // Handle column header click for filtering
  function handleColumnHeaderClick(column) {
    if (column === 'ports') {
      if (filterOption === 'withPorts') {
        // If already showing only processes with ports, switch back to all
        handleFilterChange('all');
      } else {
        // Show only processes with ports
        handleFilterChange('withPorts');
      }
    } else if (column === 'cpu') {
      if (filterOption === 'sortByCPU') {
        // If already sorting by CPU, switch back to all
        handleFilterChange('all');
      } else {
        // Sort by CPU usage
        handleFilterChange('sortByCPU');
      }
    } else if (column === 'portsAndCpu') {
      if (filterOption === 'withPortsAndSortByCPU') {
        // If already filtering and sorting, switch back to all
        handleFilterChange('all');
      } else {
        // Filter and sort
        handleFilterChange('withPortsAndSortByCPU');
      }
    }
  }

  // Format memory usage synchronously
  function formatMemoryUsage(bytes) {
    // This is a synchronous version to avoid Promise issues
    if (bytes < 1024) {
      return bytes + " B";
    } else if (bytes < 1024 * 1024) {
      return (bytes / 1024).toFixed(2) + " KB";
    } else if (bytes < 1024 * 1024 * 1024) {
      return (bytes / (1024 * 1024)).toFixed(2) + " MB";
    } else {
      return (bytes / (1024 * 1024 * 1024)).toFixed(2) + " GB";
    }
  }

  // Refresh processes manually
  function refreshProcesses() {
    updateProcesses();
  }

  // Handle process search with debounce
  function handleProcessSearch() {
    // Clear any existing timeout
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }
    
    // Set a new timeout to delay the search slightly while typing
    searchTimeout = setTimeout(() => {
      applyProcessSearch();
    }, 300); // 300ms delay
  }

  // Apply search filter to processes
  function applyProcessSearch() {
    isSearching = true;
    applyFilters();
    isSearching = false;
  }

  // Search processes by port - this is now only used for exact port searches
  async function searchProcessesByPort(port) {
    try {
      const processes = await SearchProcessesByPort(port);
      filteredProcesses = processes;
    } catch (error) {
      console.error(`Error searching for port ${port}:`, error);
      processError = `Error searching for port ${port}: ${error.message}`;
      filteredProcesses = [];
    }
  }

  // Show confirmation before killing a process
  function confirmKillProcess(pid, name) {
    processToKill = { pid, name };
    showKillConfirmation = true;
  }

  // Execute process termination
  async function executeKillProcess() {
    try {
      await KillProcess(processToKill.pid);
      showKillConfirmation = false;
      
      // Refresh the process list
      setTimeout(updateProcesses, 500);
    } catch (error) {
      console.error(`Error killing process ${processToKill.pid}:`, error);
      processError = `Error killing process: ${error.message}`;
    }
  }

  // Format timestamp for process start time
  function formatProcessTime(timestamp) {
    const date = new Date(timestamp);
    return date.toLocaleString();
  }

  onMount(() => {
    updateMetrics();
    updateHistoryData();
    
    // Delay process loading to avoid initial UI freeze
    setTimeout(updateProcesses, 1000);
    
    const interval = setInterval(updateMetrics, 2000);
    const historyInterval = setInterval(updateHistoryData, 30000); // Update history every 30 seconds
    const processInterval = setInterval(updateProcesses, 30000); // Update processes every 30 seconds instead of 10
    
    return () => {
      clearInterval(interval);
      clearInterval(historyInterval);
      clearInterval(processInterval);
    };
  });
</script>

<div class="app-container">
  <!-- Side Navigation -->
  <nav class="side-nav">
    <div class="nav-header">
      <h2>DevEx</h2>
    </div>
    <div class="nav-items">
      {#each categories as category}
        <button 
          class="nav-item {selectedCategory === category.id ? 'active' : ''}"
          on:click={() => selectedCategory = category.id}
        >
          <span class="icon">{category.icon}</span>
          {category.name}
        </button>
      {/each}
    </div>
  </nav>

  <!-- Main Content -->
  <main class="main-content">
    {#if selectedCategory === 'dashboard'}
      <!-- System Health Section -->
      <section class="dashboard-section">
        <h2 class="section-title">System Health</h2>
        <div class="dashboard-grid">
          <div class="card status-{getStatusClass(cpuInfo)}">
            <h3>CPU Usage</h3>
            <div class="metric status-indicator">
              <span class="status-dot {getStatusClass(cpuInfo)}"></span>
              {cpuInfo}
            </div>
            <div class="details">
              <span class="status-text">{getStatusText(cpuInfo)}</span>
              {cpuDetails}
            </div>
          </div>
          <div class="card status-{getStatusClass(ramInfo)}">
            <h3>Memory Usage</h3>
            <div class="metric status-indicator">
              <span class="status-dot {getStatusClass(ramInfo)}"></span>
              {ramInfo}
            </div>
            <div class="details">
              <span class="status-text">{getStatusText(ramInfo)}</span>
              {ramDetails}
            </div>
          </div>
          <div class="card status-{getStatusClass(diskInfo)}">
            <h3>Disk Usage</h3>
            <div class="metric status-indicator">
              <span class="status-dot {getStatusClass(diskInfo)}"></span>
              {diskInfo}
            </div>
            <div class="details">
              <span class="status-text">{getStatusText(diskInfo)}</span>
              {diskDetails}
            </div>
          </div>
        </div>
      </section>

      <!-- Docker Status Section -->
      <section class="dashboard-section">
        <h2 class="section-title">Docker Status</h2>
        <div class="dashboard-grid">
          <div class="card status-{dockerStatus.daemonRunning ? 'success' : 'error'}">
            <h3>Docker Daemon</h3>
            <div class="metric status-indicator">
              <span class="status-dot {dockerStatus.daemonRunning ? 'success' : 'error'}"></span>
              {dockerStatus.daemonRunning ? 'Running' : 'Not Running'}
            </div>
            <div class="details">
              <span class="status-text">{dockerStatus.daemonRunning ? 'Healthy' : 'Critical'}</span>
              Version: {dockerStatus.version}
            </div>
          </div>
          <div class="card status-{dockerMetrics.containersAll > 0 ? 'success' : 'pending'}">
            <h3>Docker Resources</h3>
            <div class="metric status-indicator">
              <span class="status-dot {dockerMetrics.containersAll > 0 ? 'success' : 'pending'}"></span>
              {dockerMetrics.containersUp}/{dockerMetrics.containersAll} Containers
            </div>
            <div class="details">
              <span class="status-text">{dockerMetrics.containersAll > 0 ? 'Active' : 'No Containers'}</span>
              {dockerMetrics.imagesCount} Images • {dockerMetrics.diskUsage} Used
            </div>
          </div>
          <div class="card status-{dockerMetrics.networkStatus !== 'Loading...' ? 'success' : 'pending'}">
            <h3>Docker Network</h3>
            <div class="metric status-indicator">
              <span class="status-dot {dockerMetrics.networkStatus !== 'Loading...' ? 'success' : 'pending'}"></span>
              {dockerMetrics.networkStatus}
            </div>
            <div class="details">
              <span class="status-text">{dockerMetrics.networkStatus !== 'Loading...' ? 'Available' : 'Checking...'}</span>
              Docker network status and configuration
            </div>
          </div>
        </div>
      </section>

      <!-- Network Status Section -->
      <section class="dashboard-section">
        <h2 class="section-title">Network Status</h2>
        <div class="dashboard-grid">
          <div class="card status-{networkStatus.internetConnected ? 'success' : 'error'}">
            <h3>Internet Connection</h3>
            <div class="metric status-indicator">
              <span class="status-dot {networkStatus.internetConnected ? 'success' : 'error'}"></span>
              {networkStatus.internetConnected ? 'Connected' : 'Disconnected'}
            </div>
            <div class="details">
              <span class="status-text">{networkStatus.internetConnected ? 'Online' : 'Offline'}</span>
              Connection to Google.com
            </div>
          </div>
          <div class="card status-{networkStatus.pingStatus !== 'Loading...' ? (networkStatus.pingLatency < 100 ? 'success' : (networkStatus.pingLatency < 200 ? 'warning' : 'error')) : 'pending'}">
            <h3>Ping Status</h3>
            <div class="metric status-indicator">
              <span class="status-dot {networkStatus.pingStatus !== 'Loading...' ? (networkStatus.pingLatency < 100 ? 'success' : (networkStatus.pingLatency < 200 ? 'warning' : 'error')) : 'pending'}"></span>
              {networkStatus.pingStatus}
            </div>
            <div class="details">
              <span class="status-text">
                {networkStatus.pingStatus !== 'Loading...' ? 
                  (networkStatus.pingLatency < 100 ? 'Good' : 
                   (networkStatus.pingLatency < 200 ? 'Fair' : 'Poor')) : 
                  'Checking...'}
              </span>
              Latency to Google.com
            </div>
          </div>
          <div class="card status-{networkStatus.dnsStatus === 'Working' ? 'success' : (networkStatus.dnsStatus === 'Loading...' ? 'pending' : 'error')}">
            <h3>DNS Resolution</h3>
            <div class="metric status-indicator">
              <span class="status-dot {networkStatus.dnsStatus === 'Working' ? 'success' : (networkStatus.dnsStatus === 'Loading...' ? 'pending' : 'error')}"></span>
              {networkStatus.dnsStatus}
            </div>
            <div class="details">
              <span class="status-text">{networkStatus.dnsStatus === 'Working' ? 'Healthy' : (networkStatus.dnsStatus === 'Loading...' ? 'Checking...' : 'Failed')}</span>
              DNS resolution for Google.com
            </div>
          </div>
        </div>
      </section>
    {:else if selectedCategory === 'history'}
      <section class="dashboard-section">
        <h2 class="section-title">Historical Metrics</h2>
        <div class="time-range-selector">
          <label for="timeRange">Time Range:</label>
          <select id="timeRange" value={selectedTimeRange} on:change={handleTimeRangeChange}>
            <option value="5">Last 5 minutes</option>
            <option value="15">Last 15 minutes</option>
            <option value="30">Last 30 minutes</option>
            <option value="60">Last 1 hour</option>
          </select>
          <button class="refresh-button" on:click={updateHistoryData}>Refresh Data</button>
        </div>
        
        {#if historyError}
          <div class="error-message">{historyError}</div>
        {/if}
        
        <div class="dashboard-grid">
          <div class="card chart-card">
            <h2>CPU Usage History</h2>
            <div class="chart-container">
              {#if cpuHistory.length === 0}
                <div class="no-data">No data available</div>
              {:else}
                <div class="chart">
                  <div class="chart-y-axis">
                    <span>100%</span>
                    <span>75%</span>
                    <span>50%</span>
                    <span>25%</span>
                    <span>0%</span>
                  </div>
                  <div class="chart-content">
                    <svg width="100%" height="200">
                      {#if cpuHistory.length > 1}
                        <polyline
                          points={cpuHistory.map((point, index) => {
                            // Scale x from 0-100% based on position in array
                            const x = (index / (cpuHistory.length - 1)) * 100;
                            // Scale y from 0-200px (inverted, 0% at bottom)
                            const y = 200 - (point.value * 2);
                            return `${x}%,${y}`;
                          }).join(' ')}
                          fill="none"
                          stroke="#4a90e2"
                          stroke-width="2"
                        />
                        
                        <!-- Add data points -->
                        {#each cpuHistory as point, index}
                          <circle 
                            cx={`${(index / (cpuHistory.length - 1)) * 100}%`} 
                            cy={200 - (point.value * 2)} 
                            r="3" 
                            fill="#4a90e2" 
                          />
                        {/each}
                      {/if}
                    </svg>
                    <div class="chart-x-axis">
                      {#each cpuHistory.filter((_, i) => i % Math.max(1, Math.floor(cpuHistory.length / 5)) === 0) as point}
                        <span>{formatTimestamp(point.timestamp)}</span>
                      {/each}
                    </div>
                  </div>
                </div>
              {/if}
            </div>
          </div>
          
          <div class="card chart-card">
            <h2>RAM Usage History</h2>
            <div class="chart-container">
              {#if ramHistory.length === 0}
                <div class="no-data">No data available</div>
              {:else}
                <div class="chart">
                  <div class="chart-y-axis">
                    <span>100%</span>
                    <span>75%</span>
                    <span>50%</span>
                    <span>25%</span>
                    <span>0%</span>
                  </div>
                  <div class="chart-content">
                    <svg width="100%" height="200">
                      {#if ramHistory.length > 1}
                        <polyline
                          points={ramHistory.map((point, index) => {
                            const x = (index / (ramHistory.length - 1)) * 100;
                            const y = 200 - (point.value * 2);
                            return `${x}%,${y}`;
                          }).join(' ')}
                          fill="none"
                          stroke="#f44336"
                          stroke-width="2"
                        />
                        
                        <!-- Add data points -->
                        {#each ramHistory as point, index}
                          <circle 
                            cx={`${(index / (ramHistory.length - 1)) * 100}%`} 
                            cy={200 - (point.value * 2)} 
                            r="3" 
                            fill="#f44336" 
                          />
                        {/each}
                      {/if}
                    </svg>
                    <div class="chart-x-axis">
                      {#each ramHistory.filter((_, i) => i % Math.max(1, Math.floor(ramHistory.length / 5)) === 0) as point}
                        <span>{formatTimestamp(point.timestamp)}</span>
                      {/each}
                    </div>
                  </div>
                </div>
              {/if}
            </div>
          </div>
          
          <div class="card chart-card">
            <h2>Disk Usage History</h2>
            <div class="chart-container">
              {#if diskHistory.length === 0}
                <div class="no-data">No data available</div>
              {:else}
                <div class="chart">
                  <div class="chart-y-axis">
                    <span>100%</span>
                    <span>75%</span>
                    <span>50%</span>
                    <span>25%</span>
                    <span>0%</span>
                  </div>
                  <div class="chart-content">
                    <svg width="100%" height="200">
                      {#if diskHistory.length > 1}
                        <polyline
                          points={diskHistory.map((point, index) => {
                            const x = (index / (diskHistory.length - 1)) * 100;
                            const y = 200 - (point.value * 2);
                            return `${x}%,${y}`;
                          }).join(' ')}
                          fill="none"
                          stroke="#4caf50"
                          stroke-width="2"
                        />
                        
                        <!-- Add data points -->
                        {#each diskHistory as point, index}
                          <circle 
                            cx={`${(index / (diskHistory.length - 1)) * 100}%`} 
                            cy={200 - (point.value * 2)} 
                            r="3" 
                            fill="#4caf50" 
                          />
                        {/each}
                      {/if}
                    </svg>
                    <div class="chart-x-axis">
                      {#each diskHistory.filter((_, i) => i % Math.max(1, Math.floor(diskHistory.length / 5)) === 0) as point}
                        <span>{formatTimestamp(point.timestamp)}</span>
                      {/each}
                    </div>
                  </div>
                </div>
              {/if}
            </div>
          </div>
        </div>
      </section>
    {:else if selectedCategory === 'processes'}
      <section class="dashboard-section">
        <h2 class="section-title">Process Management</h2>
        
        <!-- Search Bar -->
        <div class="search-container">
          <input 
            type="text" 
            placeholder="Search by port (e.g. 8080) or process name" 
            bind:value={processSearchQuery}
            on:input={handleProcessSearch}
          />
          <button class="search-button" on:click={handleProcessSearch}>Search</button>
          <button class="refresh-button" on:click={refreshProcesses} disabled={isLoadingProcesses}>
            {isLoadingProcesses ? 'Loading...' : 'Refresh'}
          </button>
        </div>
        
        <!-- Filter Options -->
        <div class="filter-options">
          <div class="filter-info">
            {#if filterOption === 'all'}
              Showing all processes
            {:else if filterOption === 'withPorts'}
              Showing only processes with ports
            {:else if filterOption === 'sortByCPU'}
              Sorting by CPU usage (highest first)
            {:else if filterOption === 'withPortsAndSortByCPU'}
              Showing only processes with ports & sorting by CPU
            {/if}
            <button class="reset-filter" on:click={() => handleFilterChange('all')}>Reset</button>
          </div>
        </div>
        
        {#if processError}
          <div class="error-message">{processError}</div>
        {/if}
        
        <!-- Process List -->
        <div class="process-list-container">
          <table class="process-table">
            <thead>
              <tr>
                <th>PID</th>
                <th>Name</th>
                <th class="clickable-header {filterOption === 'withPorts' || filterOption === 'withPortsAndSortByCPU' ? 'active-filter' : ''}" on:click={() => handleColumnHeaderClick('ports')}>
                  Ports
                  <span class="filter-icon">🔍</span>
                </th>
                <th class="clickable-header {filterOption === 'sortByCPU' || filterOption === 'withPortsAndSortByCPU' ? 'active-filter' : ''}" on:click={() => handleColumnHeaderClick('cpu')}>
                  CPU
                  <span class="filter-icon">↓</span>
                </th>
                <th>Memory</th>
                <th>User</th>
                <th>Started</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {#if isLoadingProcesses && filteredProcesses.length === 0}
                <tr>
                  <td colspan="8" class="no-data">
                    <div class="loading-spinner"></div>
                    <div>Loading processes...</div>
                  </td>
                </tr>
              {:else if filteredProcesses.length === 0}
                <tr>
                  <td colspan="8" class="no-data">
                    {isSearching ? 'No processes found matching your search' : 'No processes available'}
                  </td>
                </tr>
              {:else}
                {#each filteredProcesses as process}
                  <tr>
                    <td>{process.process.pid}</td>
                    <td class="process-name" title={process.process.commandLine}>
                      {process.process.name}
                    </td>
                    <td>
                      {#if process.ports && process.ports.length > 0}
                        <div class="port-list">
                          {#each process.ports as port}
                            <span class="port-badge {port.state === 'LISTEN' ? 'listening' : ''} {port.isMatch ? 'matching' : ''}">
                              {port.port} ({port.protocol})
                            </span>
                          {/each}
                        </div>
                      {:else}
                        <span class="no-ports">No open ports</span>
                      {/if}
                    </td>
                    <td>{process.process.cpuPercent.toFixed(1)}%</td>
                    <td>{formatMemoryUsage(process.process.memoryUsage)}</td>
                    <td>{process.process.username}</td>
                    <td>{formatProcessTime(process.process.startTime)}</td>
                    <td>
                      <button 
                        class="kill-button" 
                        on:click={() => confirmKillProcess(process.process.pid, process.process.name)}
                      >
                        Terminate
                      </button>
                    </td>
                  </tr>
                {/each}
              {/if}
            </tbody>
          </table>
        </div>
        
        <!-- Kill Confirmation Modal -->
        {#if showKillConfirmation}
          <div class="modal-overlay">
            <div class="modal-content">
              <h3>Confirm Process Termination</h3>
              <p>Are you sure you want to terminate process "{processToKill.name}" (PID: {processToKill.pid})?</p>
              <div class="modal-actions">
                <button class="cancel-button" on:click={() => showKillConfirmation = false}>Cancel</button>
                <button class="confirm-button" on:click={executeKillProcess}>Terminate</button>
              </div>
            </div>
          </div>
        {/if}
      </section>
    {:else if selectedCategory === 'logs'}
      <div class="dashboard-grid">
        <div class="card">
          <h2>System Logs</h2>
          <div class="metric">Coming Soon</div>
          <div class="details">System log viewer will be implemented here</div>
        </div>
        <div class="card">
          <h2>Application Logs</h2>
          <div class="metric">Coming Soon</div>
          <div class="details">Application-specific log viewer</div>
        </div>
        <div class="card">
          <h2>Log Analysis</h2>
          <div class="metric">Coming Soon</div>
          <div class="details">Log analysis and filtering tools</div>
        </div>
      </div>
    {:else if selectedCategory === 'settings'}
      <div class="dashboard-grid">
        <div class="card">
          <h2>General Settings</h2>
          <div class="metric">Coming Soon</div>
          <div class="details">General application settings and preferences</div>
        </div>
        <div class="card">
          <h2>Monitoring Settings</h2>
          <div class="metric">Coming Soon</div>
          <div class="details">Configure monitoring thresholds and alerts</div>
        </div>
        <div class="card">
          <h2>Display Settings</h2>
          <div class="metric">Coming Soon</div>
          <div class="details">Customize the dashboard appearance</div>
        </div>
        <div class="card">
          <h2>Profiling Tools</h2>
          <div class="metric profiling-container">
            <div class="debug-emoji">🐞</div>
            <button class="pprof-button" on:click={openPprofDashboard}>
              Open PProf Dashboard
            </button>
          </div>
          <div class="details">
            Access profiling data and debug information
          </div>
        </div>
        <div class="card">
          <h2>Notification Settings</h2>
          <div class="metric">Coming Soon</div>
          <div class="details">Configure notification preferences and alerts</div>
        </div>
      </div>
    {/if}
  </main>
</div>

<style>
  .app-container {
    display: flex;
    height: 100vh;
    background-color: #f5f5f5;
  }

  .side-nav {
    width: 250px;
    background-color: #1a1a1a;
    color: white;
    padding: 1rem;
    display: flex;
    flex-direction: column;
  }

  .nav-header {
    padding: 1rem 0;
    border-bottom: 1px solid #333;
    margin-bottom: 1rem;
  }

  .nav-header h2 {
    margin: 0;
    font-size: 1.5rem;
    color: #fff;
  }

  .nav-items {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    border: none;
    background: none;
    color: #fff;
    cursor: pointer;
    border-radius: 8px;
    transition: background-color 0.2s;
  }

  .nav-item:hover {
    background-color: #333;
  }

  .nav-item.active {
    background-color: #4a90e2;
  }

  .icon {
    font-size: 1.25rem;
  }

  .main-content {
    flex: 1;
    padding: 2rem;
    overflow-y: auto;
    color: #333;
  }

  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.25rem;
  }

  .card {
    background-color: white;
    border-radius: 12px;
    padding: 1.25rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s, box-shadow 0.2s, background-color 0.3s;
    border-left: 4px solid transparent;
  }

  .card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }

  .card h2 {
    margin: 0 0 0.75rem 0;
    color: #333;
    font-size: 1.1rem;
  }

  .card h3 {
    margin: 0 0 0.75rem 0;
    color: #333;
    font-size: 1.1rem;
  }

  .metric {
    font-size: 1.5rem;
    font-weight: bold;
    color: #4a90e2;
    margin-bottom: 0.5rem;
    text-align: center;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 2.5rem;
  }

  .details {
    font-size: 0.8rem;
    color: #666;
    white-space: pre-line;
  }

  .pprof-button {
    background-color: #4a90e2;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 6px;
    cursor: pointer;
    font-size: 1rem;
    transition: background-color 0.2s;
  }

  .pprof-button:hover {
    background-color: #357abd;
  }

  .debug-emoji {
    font-size: 1.75rem;
    margin-bottom: 0.25rem;
    line-height: 1;
  }

  .profiling-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding-top: 0;
    min-height: auto;
  }

  .dashboard-section {
    margin-bottom: 2rem;
  }

  .section-title {
    font-size: 1.5rem;
    color: #333;
    margin-bottom: 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid #e0e0e0;
  }

  .status-indicator {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .status-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
  }

  .status-dot.pending {
    background-color: #ffa500;
  }

  .status-dot.success {
    background-color: #4caf50;
  }

  .status-dot.error {
    background-color: #f44336;
  }

  .card.status-success {
    border-left-color: #4caf50;
    background-color: rgba(76, 175, 80, 0.05);
  }

  .card.status-warning {
    border-left-color: #ffa500;
    background-color: rgba(255, 165, 0, 0.05);
  }

  .card.status-error {
    border-left-color: #f44336;
    background-color: rgba(244, 67, 54, 0.05);
  }

  .card.status-pending {
    border-left-color: #9e9e9e;
    background-color: rgba(158, 158, 158, 0.05);
  }

  .status-text {
    font-weight: 500;
    margin-right: 0.5rem;
  }

  .status-text:empty {
    display: none;
  }

  .time-range-selector {
    margin-bottom: 1rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .time-range-selector label {
    color: #333;
    font-weight: 500;
  }

  .time-range-selector select {
    padding: 0.5rem;
    border-radius: 4px;
    border: 1px solid #ddd;
    background-color: white;
    font-size: 0.9rem;
    color: #333;
  }

  .refresh-button {
    background-color: #4a90e2;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background-color 0.2s;
  }

  .refresh-button:hover {
    background-color: #357abd;
  }

  .refresh-button:disabled {
    background-color: #cccccc;
    cursor: not-allowed;
  }

  .chart-card {
    grid-column: span 1;
  }

  .chart-container {
    height: 250px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  .no-data {
    text-align: center;
    color: #666;
    font-style: italic;
    padding: 2rem 0;
  }

  .chart {
    display: flex;
    height: 220px;
  }

  .chart-y-axis {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding-right: 0.5rem;
    font-size: 0.7rem;
    color: #666;
    width: 40px;
  }

  .chart-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    position: relative;
  }

  .chart-x-axis {
    display: flex;
    justify-content: space-between;
    font-size: 0.7rem;
    color: #666;
    margin-top: 0.5rem;
  }

  svg {
    background-color: rgba(0, 0, 0, 0.02);
    border-radius: 4px;
  }

  .error-message {
    background-color: #fff3cd;
    color: #856404;
    padding: 0.75rem;
    margin-bottom: 1rem;
    border-radius: 4px;
    border-left: 4px solid #ffeeba;
  }

  option {
    color: #333;
  }

  .search-container {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1rem;
  }

  .search-container input {
    flex: 1;
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 0.9rem;
  }

  .search-button {
    background-color: #4a90e2;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
  }

  .process-list-container {
    overflow-x: auto;
  }

  .process-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 0.9rem;
  }

  .process-table th,
  .process-table td {
    padding: 0.75rem;
    text-align: left;
    border-bottom: 1px solid #eee;
  }

  .process-table th {
    background-color: #f5f5f5;
    font-weight: 600;
    color: #333;
  }

  .process-table tr:hover {
    background-color: rgba(74, 144, 226, 0.05);
  }

  .process-name {
    max-width: 200px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .port-list {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
  }

  .port-badge {
    display: inline-block;
    padding: 0.2rem 0.4rem;
    background-color: #eee;
    border-radius: 4px;
    font-size: 0.8rem;
  }

  .port-badge.listening {
    background-color: #4caf50;
    color: white;
  }

  .port-badge.matching {
    background-color: #ff9800;
    color: white;
    font-weight: bold;
  }

  .port-badge.listening.matching {
    background-color: #ff5722;
  }

  .no-ports {
    color: #999;
    font-style: italic;
  }

  .kill-button {
    background-color: #f44336;
    color: white;
    border: none;
    padding: 0.3rem 0.6rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
  }

  .kill-button:hover {
    background-color: #d32f2f;
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .modal-content {
    background-color: white;
    padding: 1.5rem;
    border-radius: 8px;
    width: 400px;
    max-width: 90%;
  }

  .modal-content h3 {
    margin-top: 0;
    color: #333;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    margin-top: 1.5rem;
  }

  .cancel-button {
    background-color: #ccc;
    color: #333;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
  }

  .confirm-button {
    background-color: #f44336;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
  }

  .no-data {
    text-align: center;
    padding: 2rem;
    color: #666;
    font-style: italic;
  }

  .loading-spinner {
    display: inline-block;
    width: 30px;
    height: 30px;
    border: 3px solid rgba(74, 144, 226, 0.3);
    border-radius: 50%;
    border-top-color: #4a90e2;
    animation: spin 1s ease-in-out infinite;
    margin-bottom: 10px;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .filter-options {
    display: flex;
    flex-wrap: wrap;
    gap: 1rem;
    margin-bottom: 1rem;
    padding: 0.75rem;
    background-color: #f5f5f5;
    border-radius: 4px;
  }

  .filter-info {
    display: flex;
    align-items: center;
    gap: 1rem;
    font-size: 0.9rem;
    color: #555;
  }

  .reset-filter {
    background-color: #e0e0e0;
    border: none;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.8rem;
    transition: background-color 0.2s;
  }

  .reset-filter:hover {
    background-color: #ccc;
  }

  .clickable-header {
    cursor: pointer;
    position: relative;
    transition: background-color 0.2s;
  }

  .clickable-header:hover {
    background-color: #e0e0e0;
  }

  .clickable-header.active-filter {
    background-color: #4a90e2;
    color: white;
  }

  .filter-icon {
    margin-left: 0.25rem;
    font-size: 0.8rem;
    opacity: 0.7;
  }

  .clickable-header:hover .filter-icon {
    opacity: 1;
  }
</style>