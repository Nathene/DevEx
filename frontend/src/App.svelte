<script>
  import { GetCPUInfo, GetCPUDetails, GetRAMInfo, GetRAMDetails, GetDiskInfo, GetDiskDetails, GetDockerStatus, GetDockerMetrics, GetNetworkStatus, GetCPUHistory, GetRAMHistory, GetDiskHistory, GetAllProcesses, SearchProcessesByPort, KillProcess, FormatProcessBytes, GetTopMemoryProcesses, GetTopCPUProcesses, GetTopDiskProcesses } from '../wailsjs/go/main/App';
  import { BrowserOpenURL } from '../wailsjs/runtime/runtime';
  import { onMount } from 'svelte';
  import { GetAllServers, StartServer, StopServer, GetAllDatabases, ConnectDatabase, DisconnectDatabase, SendAPIRequest, GetAllGitRepos, RefreshGitRepo, RefreshAllGitRepos, AddGitRepo, RemoveGitRepo, GetGitRepoChanges, OpenInVSCode, OpenFolderPicker } from '../wailsjs/go/main/App';
  import * as models from '../wailsjs/go/models';

  // Debug log to check if models is imported correctly
  console.log("Models import:", models);
  console.log("GitRepoInfo available:", models.devtools && models.devtools.GitRepoInfo);

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
    { id: 'devtools', name: 'Dev Tools', icon: '🛠️' },
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

  // Add state variables for DevTools
  let servers = [];
  let databases = [];
  let gitRepos = [];
  let apiRequest = {
    url: 'https://jsonplaceholder.typicode.com/posts/1',
    method: 'GET',
    headers: { 'Accept': 'application/json' },
    body: '',
    timeout: 30
  };
  let apiResponse = null;
  let isLoadingServers = false;
  let isLoadingDatabases = false;
  let isLoadingGitRepos = false;
  let isSendingRequest = false;
  let headersJson = '';
  let showAddRepoModal = false;
  let showRemoveRepoModal = false;
  let repoToRemove = null;
  let newRepo = {
    name: '',
    path: '',
    description: '',
    url: ''
  };

  // Add state for processes modals
  let showProcessesModal = false;
  let topProcesses = [];
  let currentMetricType = ''; // 'memory', 'cpu', or 'disk'
  let modalTitle = '';

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

  // Add function to load DevTools data
  async function loadDevToolsData() {
    if (selectedCategory === 'devtools') {
      await Promise.all([
        loadServers(),
        loadDatabases(),
        loadGitRepos()
      ]);
    }
  }

  // Add function to load servers
  async function loadServers() {
    isLoadingServers = true;
    try {
      servers = await GetAllServers();
    } catch (error) {
      console.error('Error loading servers:', error);
    } finally {
      isLoadingServers = false;
    }
  }

  // Add function to start a server
  async function handleStartServer(serverId) {
    try {
      const updatedServer = await StartServer(serverId);
      servers = servers.map(server => 
        server.id === serverId ? updatedServer : server
      );
    } catch (error) {
      console.error('Error starting server:', error);
    }
  }

  // Add function to stop a server
  async function handleStopServer(serverId) {
    try {
      const updatedServer = await StopServer(serverId);
      servers = servers.map(server => 
        server.id === serverId ? updatedServer : server
      );
    } catch (error) {
      console.error('Error stopping server:', error);
    }
  }

  // Add function to load databases
  async function loadDatabases() {
    isLoadingDatabases = true;
    try {
      databases = await GetAllDatabases();
    } catch (error) {
      console.error('Error loading databases:', error);
    } finally {
      isLoadingDatabases = false;
    }
  }

  // Add function to connect to a database
  async function handleConnectDatabase(databaseId) {
    try {
      const updatedDatabase = await ConnectDatabase(databaseId);
      databases = databases.map(db => 
        db.id === databaseId ? updatedDatabase : db
      );
    } catch (error) {
      console.error('Error connecting to database:', error);
    }
  }

  // Add function to disconnect from a database
  async function handleDisconnectDatabase(databaseId) {
    try {
      const updatedDatabase = await DisconnectDatabase(databaseId);
      databases = databases.map(db => 
        db.id === databaseId ? updatedDatabase : db
      );
    } catch (error) {
      console.error('Error disconnecting from database:', error);
    }
  }

  // Add function to send API request
  async function handleSendAPIRequest() {
    isSendingRequest = true;
    try {
      // Parse headers from JSON string
      try {
        apiRequest.headers = JSON.parse(headersJson);
      } catch (error) {
        console.error('Error parsing headers JSON:', error);
      }
      
      apiResponse = await SendAPIRequest(apiRequest);
    } catch (error) {
      console.error('Error sending API request:', error);
      apiResponse = {
        statusCode: 0,
        status: 'Error',
        headers: {},
        body: '',
        duration: 0,
        error: `Error sending request: ${error.message || error}`
      };
    } finally {
      isSendingRequest = false;
    }
  }

  // Add function to load Git repositories
  async function loadGitRepos() {
    isLoadingGitRepos = true;
    try {
      gitRepos = await GetAllGitRepos();
    } catch (error) {
      console.error('Error loading Git repositories:', error);
    } finally {
      isLoadingGitRepos = false;
    }
  }

  // Add function to refresh all Git repositories
  async function handleRefreshAllRepos() {
    // Save current scroll position
    const scrollPosition = window.scrollY;
    
    // Add a refreshing class to each repository card instead of setting isLoadingGitRepos
    const repoCards = document.querySelectorAll('.devtools-card');
    repoCards.forEach(card => {
      card.classList.add('refreshing');
    });
    
    try {
      // Get updated repositories without changing the loading state
      const updatedRepos = await RefreshAllGitRepos();
      
      // Update repositories without full page reload effect
      gitRepos = updatedRepos;
    } catch (error) {
      console.error('Error refreshing all Git repositories:', error);
    } finally {
      // Remove refreshing class
      setTimeout(() => {
        const updatedRepoCards = document.querySelectorAll('.devtools-card');
        updatedRepoCards.forEach(card => {
          card.classList.remove('refreshing');
        });
        
        // Restore scroll position
        window.scrollTo({
          top: scrollPosition,
          behavior: 'auto'
        });
      }, 100);
    }
  }

  // Add function to refresh a Git repository
  async function handleRefreshRepo(repoId) {
    try {
      const updatedRepo = await RefreshGitRepo(repoId);
      gitRepos = gitRepos.map(repo => 
        repo.id === repoId ? updatedRepo : repo
      );
    } catch (error) {
      console.error('Error refreshing Git repository:', error);
    }
  }

  // Add function to add a new Git repository
  async function handleAddRepo() {
    if (!newRepo.name || !newRepo.path) {
      alert("Repository name and path are required");
      return;
    }

    try {
      // Create a new repository using the GitRepoInfo model
      const repoData = {
        id: '',
        name: newRepo.name,
        path: newRepo.path,
        branch: 'main',
        status: 'unknown',
        lastCommit: '',
        lastCommitBy: '',
        lastUpdated: new Date().toISOString(),
        changes: 0,
        description: newRepo.description,
        url: newRepo.url || '' // URL is now optional
      };
      
      // Use the createFrom method to ensure proper structure
      const repoToAdd = models.devtools.GitRepoInfo.createFrom(repoData);

      await AddGitRepo(repoToAdd);
      
      // Reset form fields
      newRepo.name = '';
      newRepo.path = '';
      newRepo.description = '';
      newRepo.url = '';
      
      // Close the modal
      showAddRepoModal = false;
      
      // Refresh the repo list
      await loadGitRepos();
    } catch (error) {
      console.error("Error adding repository:", error);
      alert("Failed to add repository: " + error.message);
    }
  }

  // Add function to remove a Git repository
  async function handleRemoveRepo(repo) {
    console.log("Attempting to remove repository:", repo);
    // Show the custom confirmation modal instead of using confirm()
    repoToRemove = repo;
    showRemoveRepoModal = true;
  }

  // Function to confirm repository removal
  async function confirmRemoveRepo() {
    if (!repoToRemove) return;
    
    try {
      console.log("Confirmed removal of repository with ID:", repoToRemove.id);
      await RemoveGitRepo(repoToRemove.id);
      console.log("Repository removed successfully, refreshing list");
      
      // Close the modal
      showRemoveRepoModal = false;
      repoToRemove = null;
      
      // Refresh the repo list
      await loadGitRepos();
    } catch (error) {
      console.error('Error removing Git repository:', error);
      alert(`Error removing repository: ${error.message || 'Unknown error'}`);
    }
  }

  // Function to cancel repository removal
  function cancelRemoveRepo() {
    console.log("Repository removal cancelled by user");
    showRemoveRepoModal = false;
    repoToRemove = null;
  }

  // Reset the add repository form
  function resetAddRepoForm() {
    newRepo = {
      name: '',
      path: '',
      description: '',
      url: ''
    };
  }

  // Function to open repository in VS Code
  async function handleOpenInVSCode(path) {
    console.log("Opening repository in VS Code:", path);
    try {
      await OpenInVSCode(path);
      console.log("Repository opened in VS Code successfully");
    } catch (error) {
      console.error("Error opening repository in VS Code:", error);
      alert(`Error opening repository: ${error.message || 'Unknown error'}`);
    }
  }

  // Function to open folder picker
  async function handleBrowseForFolder() {
    console.log("Opening folder picker");
    try {
      const selectedPath = await OpenFolderPicker();
      if (selectedPath) {
        console.log("Selected path:", selectedPath);
        newRepo.path = selectedPath;
      }
    } catch (error) {
      console.error("Error opening folder picker:", error);
      alert(`Error selecting folder: ${error.message || 'Unknown error'}`);
    }
  }

  // Add function to fetch top memory processes
  async function fetchTopMemoryProcesses() {
    try {
      isLoadingProcesses = true;
      currentMetricType = 'memory';
      modalTitle = 'Top Memory Consuming Processes';
      topProcesses = await GetTopMemoryProcesses();
      showProcessesModal = true;
    } catch (error) {
      console.error("Error fetching top memory processes:", error);
    } finally {
      isLoadingProcesses = false;
    }
  }

  // Add function to fetch top CPU processes
  async function fetchTopCPUProcesses() {
    try {
      isLoadingProcesses = true;
      currentMetricType = 'cpu';
      modalTitle = 'Top CPU Consuming Processes';
      topProcesses = await GetTopCPUProcesses();
      showProcessesModal = true;
    } catch (error) {
      console.error("Error fetching top CPU processes:", error);
    } finally {
      isLoadingProcesses = false;
    }
  }

  // Add function to fetch top disk processes
  async function fetchTopDiskProcesses() {
    try {
      isLoadingProcesses = true;
      currentMetricType = 'disk';
      modalTitle = 'Top Disk Consuming Processes';
      topProcesses = await GetTopDiskProcesses();
      showProcessesModal = true;
    } catch (error) {
      console.error("Error fetching top disk processes:", error);
    } finally {
      isLoadingProcesses = false;
    }
  }

  // Function to get the percentage of a resource used by a process
  function getResourcePercentage(process, type) {
    if (type === 'cpu') {
      return process.process.cpuPercent.toFixed(1) + '%';
    } else if (type === 'memory') {
      if (ramDetails.includes("Total:") && ramDetails.split("\n").length > 1 && 
          ramDetails.split("\n")[1].includes("Total:") && 
          ramDetails.split("\n")[1].split("Total:")[1].trim().split(" ").length > 1) {
        const totalRAM = parseFloat(ramDetails.split("\n")[1].split("Total:")[1].trim().split(" ")[0]) * 1024 * 1024 * 1024;
        return ((process.process.memoryUsage / totalRAM) * 100).toFixed(1) + '%';
      }
      return 'N/A';
    } else if (type === 'disk') {
      // For disk, we're using memory as a proxy, so show memory percentage
      if (ramDetails.includes("Total:") && ramDetails.split("\n").length > 1 && 
          ramDetails.split("\n")[1].includes("Total:") && 
          ramDetails.split("\n")[1].split("Total:")[1].trim().split(" ").length > 1) {
        const totalRAM = parseFloat(ramDetails.split("\n")[1].split("Total:")[1].trim().split(" ")[0]) * 1024 * 1024 * 1024;
        return ((process.process.memoryUsage / totalRAM) * 100).toFixed(1) + '%';
      }
      return 'N/A';
    }
    return 'N/A';
  }

  onMount(() => {
    updateMetrics();
    updateHistoryData();
    
    // Load repositories on startup
    loadGitRepos();
    
    // Delay process loading to avoid initial UI freeze
    setTimeout(updateProcesses, 1000);
    
    const interval = setInterval(updateMetrics, 2000);
    const historyInterval = setInterval(updateHistoryData, 30000); // Update history every 30 seconds
    const processInterval = setInterval(updateProcesses, 30000); // Update processes every 30 seconds instead of 10
    
    // Initialize headersJson
    headersJson = JSON.stringify(apiRequest.headers, null, 2);
    
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
          <div class="card status-{getStatusClass(cpuInfo)}" on:click={fetchTopCPUProcesses} style="cursor: pointer;">
            <h3>CPU Usage</h3>
            <div class="metric status-indicator">
              <span class="status-dot {getStatusClass(cpuInfo)}"></span>
              {cpuInfo}
            </div>
            <div class="details">
              <span class="status-text">{getStatusText(cpuInfo)}</span>
              {cpuDetails}
              <div class="click-hint">(Click to see top CPU processes)</div>
            </div>
          </div>
          <div class="card status-{getStatusClass(ramInfo)}" on:click={fetchTopMemoryProcesses} style="cursor: pointer;">
            <h3>Memory Usage</h3>
            <div class="metric status-indicator">
              <span class="status-dot {getStatusClass(ramInfo)}"></span>
              {ramInfo}
            </div>
            <div class="details">
              <span class="status-text">{getStatusText(ramInfo)}</span>
              {ramDetails}
              <div class="click-hint">(Click to see top memory processes)</div>
            </div>
          </div>
          <div class="card status-{getStatusClass(diskInfo)}" on:click={fetchTopDiskProcesses} style="cursor: pointer;">
            <h3>Disk Usage</h3>
            <div class="metric status-indicator">
              <span class="status-dot {getStatusClass(diskInfo)}"></span>
              {diskInfo}
            </div>
            <div class="details">
              <span class="status-text">{getStatusText(diskInfo)}</span>
              {diskDetails}
              <div class="click-hint">(Click to see top disk processes)</div>
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
    {:else if selectedCategory === 'devtools'}
      <section class="content-section">
        <h1>Developer Tools</h1>
        
        <!-- Development Servers Section -->
        <div class="devtools-section">
          <h2>Development Servers</h2>
          <div class="devtools-grid">
            {#if isLoadingServers}
              <div class="loading">Loading servers...</div>
            {:else if servers.length === 0}
              <div class="empty-state">No development servers found.</div>
            {:else}
              {#each servers as server}
                <div class="devtools-card">
                  <div class="devtools-card-header">
                    <h3>{server.name}</h3>
                    <span class="badge {server.status === 'running' ? 'success' : 'pending'}">{server.status}</span>
                  </div>
                  <div class="devtools-card-content">
                    <p>{server.description || 'No description'}</p>
                    <div class="server-details">
                      <div><strong>Type:</strong> {server.type}</div>
                      <div><strong>Port:</strong> {server.port}</div>
                      <div><strong>Path:</strong> {server.path}</div>
                      {#if server.status === 'running'}
                        <div><strong>PID:</strong> {server.pid}</div>
                        <div><strong>Started:</strong> {new Date(server.startTime).toLocaleString()}</div>
                      {/if}
                    </div>
                  </div>
                  <div class="devtools-card-actions">
                    {#if server.status === 'running'}
                      <button class="btn btn-danger" on:click={() => handleStopServer(server.id)}>Stop Server</button>
                      <a href={server.url} target="_blank" class="btn btn-secondary">Open URL</a>
                    {:else}
                      <button class="btn btn-primary" on:click={() => handleStartServer(server.id)}>Start Server</button>
                    {/if}
                  </div>
                </div>
              {/each}
            {/if}
          </div>
        </div>
        
        <!-- Database Connections Section -->
        <div class="devtools-section">
          <h2>Database Connections</h2>
          <div class="devtools-grid">
            {#if isLoadingDatabases}
              <div class="loading">Loading databases...</div>
            {:else if databases.length === 0}
              <div class="empty-state">No database connections found.</div>
            {:else}
              {#each databases as db}
                <div class="devtools-card">
                  <div class="devtools-card-header">
                    <h3>{db.name}</h3>
                    <span class="badge {db.status === 'connected' ? 'success' : 'pending'}">{db.status}</span>
                  </div>
                  <div class="devtools-card-content">
                    <p>{db.description || 'No description'}</p>
                    <div class="db-details">
                      <div><strong>Type:</strong> {db.type}</div>
                      <div><strong>Host:</strong> {db.host}</div>
                      <div><strong>Port:</strong> {db.port}</div>
                      <div><strong>Database:</strong> {db.database}</div>
                      <div><strong>Username:</strong> {db.username}</div>
                      {#if db.status === 'connected'}
                        <div><strong>Connected at:</strong> {new Date(db.connectedAt).toLocaleString()}</div>
                      {/if}
                    </div>
                  </div>
                  <div class="devtools-card-actions">
                    {#if db.status === 'connected'}
                      <button class="btn btn-danger" on:click={() => handleDisconnectDatabase(db.id)}>Disconnect</button>
                    {:else}
                      <button class="btn btn-primary" on:click={() => handleConnectDatabase(db.id)}>Connect</button>
                    {/if}
                  </div>
                </div>
              {/each}
            {/if}
          </div>
        </div>
        
        <!-- API Endpoint Tester Section -->
        <div class="devtools-section">
          <h2>API Endpoint Tester</h2>
          <div class="api-tester">
            <div class="request-section">
              <div class="form-group">
                <label for="api-url">URL</label>
                <input type="text" id="api-url" bind:value={apiRequest.url} placeholder="https://api.example.com/endpoint" />
              </div>
              <div class="form-row">
                <div class="form-group">
                  <label for="api-method">Method</label>
                  <select id="api-method" bind:value={apiRequest.method}>
                    <option value="GET">GET</option>
                    <option value="POST">POST</option>
                    <option value="PUT">PUT</option>
                    <option value="DELETE">DELETE</option>
                    <option value="PATCH">PATCH</option>
                  </select>
                </div>
                <div class="form-group">
                  <label for="api-timeout">Timeout (seconds)</label>
                  <input type="number" id="api-timeout" bind:value={apiRequest.timeout} min="1" max="60" />
                </div>
              </div>
              <div class="form-group">
                <label for="api-headers">Headers (JSON)</label>
                <textarea id="api-headers" bind:value={headersJson} rows="3"></textarea>
              </div>
              <div class="form-group">
                <label for="api-body">Request Body</label>
                <textarea id="api-body" bind:value={apiRequest.body} rows="5" placeholder="Request body (JSON, XML, etc.)"></textarea>
              </div>
              <div class="form-actions">
                <button class="btn btn-primary" on:click={handleSendAPIRequest} disabled={isSendingRequest}>
                  {isSendingRequest ? 'Sending...' : 'Send Request'}
                </button>
              </div>
            </div>
            <div class="response-section">
              <h3>Response</h3>
              {#if apiResponse}
                <div class="response-header">
                  <span class="status-code {apiResponse.statusCode >= 200 && apiResponse.statusCode < 300 ? 'success' : apiResponse.statusCode >= 400 ? 'error' : 'warning'}">
                    {apiResponse.statusCode} {apiResponse.status}
                  </span>
                  <span class="duration">{apiResponse.duration} ms</span>
                </div>
                {#if apiResponse.error}
                  <div class="response-error">{apiResponse.error}</div>
                {:else}
                  <div class="response-headers">
                    <h4>Headers</h4>
                    <pre style="text-align: left;">{JSON.stringify(apiResponse.headers, null, 2)}</pre>
                  </div>
                  <div class="response-body">
                    <pre style="text-align: left; white-space: pre-wrap; font-family: monospace;">{`${apiResponse.body}`}</pre>
                  </div>
                {/if}
              {:else}
                <div class="empty-response">Send a request to see the response</div>
              {/if}
            </div>
          </div>
        </div>
        
        <!-- Git Repositories Section -->
        <div class="devtools-section">
          <div class="section-header">
            <h2>Git Repositories</h2>
            <div>
              <button class="btn btn-secondary" on:click={handleRefreshAllRepos}>Refresh All</button>
              <button class="btn btn-primary" on:click={() => showAddRepoModal = true}>Add Repository</button>
            </div>
          </div>
          <div class="devtools-grid">
            {#if isLoadingGitRepos}
              <div class="loading">Loading repositories...</div>
            {:else if gitRepos.length === 0}
              <div class="empty-state">
                <p>No Git repositories found.</p>
                <button class="btn btn-primary" on:click={() => showAddRepoModal = true}>Add Repository</button>
              </div>
            {:else}
              {#each gitRepos as repo}
                <div class="devtools-card">
                  <div class="devtools-card-header">
                    <h3>{repo.name}</h3>
                    <span class="badge {repo.status === 'clean' ? 'success' : repo.status === 'modified' ? 'warning' : 'pending'}">{repo.status}</span>
                  </div>
                  <div class="devtools-card-content">
                    <p>{repo.description || 'No description'}</p>
                    <div class="repo-details">
                      <div><strong>Branch:</strong> {repo.branch}</div>
                      <div><strong>Path:</strong> {repo.path}</div>
                      <div><strong>Last Commit:</strong> {repo.lastCommit}</div>
                      <div><strong>By:</strong> {repo.lastCommitBy}</div>
                      <div><strong>Last Updated:</strong> {new Date(repo.lastUpdated).toLocaleString()}</div>
                      {#if repo.changes > 0}
                        <div><strong>Changes:</strong> <span class="changes-count">{repo.changes}</span></div>
                      {/if}
                    </div>
                  </div>
                  <div class="devtools-card-actions">
                    <button class="btn btn-primary" on:click={() => handleRefreshRepo(repo.id)}>Refresh</button>
                    <button class="btn btn-secondary" on:click={() => handleOpenInVSCode(repo.path)}>Open</button>
                    <button class="btn btn-danger" on:click={() => handleRemoveRepo(repo)}>Remove</button>
                  </div>
                </div>
              {/each}
            {/if}
          </div>
        </div>
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

<!-- Add Repository Modal -->
{#if showAddRepoModal}
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h3>Add Git Repository</h3>
        <button class="close-btn" on:click={() => showAddRepoModal = false}>&times;</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label for="repo-name">Repository Name*</label>
          <input type="text" id="repo-name" bind:value={newRepo.name} placeholder="My Project" required />
        </div>
        <div class="form-group">
          <label for="repo-path">Repository Path*</label>
          <div class="input-with-button">
            <input type="text" id="repo-path" bind:value={newRepo.path} placeholder="Path to repository" required />
            <button class="btn btn-secondary" on:click={handleBrowseForFolder}>Browse...</button>
          </div>
          <small>Path to the Git repository on your local machine.</small>
        </div>
        <div class="form-group">
          <label for="repo-description">Description</label>
          <input type="text" id="repo-description" bind:value={newRepo.description} placeholder="Description of the repository" />
        </div>
        <div class="form-group">
          <label for="repo-url">Repository URL <small>(Optional)</small></label>
          <input type="text" id="repo-url" bind:value={newRepo.url} placeholder="https://github.com/user/repo" />
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-secondary" on:click={() => showAddRepoModal = false}>Cancel</button>
        <button class="btn btn-primary" on:click={handleAddRepo}>Add Repository</button>
      </div>
    </div>
  </div>
{/if}

{#if showRemoveRepoModal && repoToRemove}
  <div class="modal-overlay">
    <div class="modal">
      <div class="modal-header">
        <h3>Confirm Repository Removal</h3>
        <button class="close-btn" on:click={cancelRemoveRepo}>&times;</button>
      </div>
      <div class="modal-body">
        <p class="confirmation-text">Are you sure you want to remove the repository "{repoToRemove.name}"?</p>
        <p class="confirmation-text">This action cannot be undone.</p>
      </div>
      <div class="modal-footer">
        <button class="btn btn-secondary" on:click={cancelRemoveRepo}>Cancel</button>
        <button class="btn btn-danger" on:click={confirmRemoveRepo}>Remove Repository</button>
      </div>
    </div>
  </div>
{/if}

<!-- Processes Modal -->
{#if showProcessesModal}
  <div class="modal-overlay">
    <div class="modal-content processes-modal">
      <div class="modal-header">
        <h3>{modalTitle}</h3>
        <button class="close-button" on:click={() => showProcessesModal = false}>×</button>
      </div>
      
      {#if isLoadingProcesses}
        <div class="loading-container">
          <div class="loading-spinner"></div>
          <p>Loading processes...</p>
        </div>
      {:else if topProcesses.length === 0}
        <p>No process information available</p>
      {:else}
        <div class="table-container">
          <table class="process-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>{currentMetricType === 'memory' ? 'Memory' : currentMetricType === 'cpu' ? 'CPU' : 'Disk'}</th>
                <th>{currentMetricType === 'memory' ? 'Memory %' : currentMetricType === 'cpu' ? 'CPU %' : 'Disk %'}</th>
                <th>PID</th>
                {#if currentMetricType !== 'cpu'}
                  <th>CPU %</th>
                {/if}
                {#if currentMetricType !== 'memory'}
                  <th>Memory</th>
                {/if}
                <th>User</th>
              </tr>
            </thead>
            <tbody>
              {#each topProcesses as process}
                <tr>
                  <td class="process-name" title={process.process.commandLine}>
                    {process.process.name}
                  </td>
                  <td>
                    {currentMetricType === 'memory' || currentMetricType === 'disk' 
                      ? formatMemoryUsage(process.process.memoryUsage) 
                      : process.process.cpuPercent.toFixed(1) + '%'}
                  </td>
                  <td>{getResourcePercentage(process, currentMetricType)}</td>
                  <td>{process.process.pid}</td>
                  {#if currentMetricType !== 'cpu'}
                    <td>{process.process.cpuPercent.toFixed(1)}%</td>
                  {/if}
                  {#if currentMetricType !== 'memory'}
                    <td>{formatMemoryUsage(process.process.memoryUsage)}</td>
                  {/if}
                  <td>{process.process.username}</td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        <div class="modal-footer">
          <p class="process-note">Click on a process name to see its full command line</p>
          {#if currentMetricType === 'disk'}
            <p class="process-note warning">Note: Disk usage is approximated by memory consumption as actual disk I/O monitoring requires additional tools</p>
          {/if}
        </div>
      {/if}
    </div>
  </div>
{/if}

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
    color: #333; /* Add explicit text color */
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

  /* DevTools Styles */
  .devtools-section {
    margin-bottom: 2rem;
  }
  
  .devtools-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
  }
  
  .devtools-card {
    background-color: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  
  .devtools-card-header {
    padding: 1rem;
    background-color: #f5f5f5;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .devtools-card-header h3 {
    margin: 0;
    font-size: 1.1rem;
  }
  
  .devtools-card-content {
    padding: 1rem;
    flex-grow: 1;
  }
  
  .devtools-card-actions {
    padding: 1rem;
    background-color: #f5f5f5;
    display: flex;
    gap: 0.5rem;
  }
  
  .server-details, .db-details, .repo-details {
    margin-top: 1rem;
    font-size: 0.9rem;
  }
  
  .server-details > div, .db-details > div, .repo-details > div {
    margin-bottom: 0.5rem;
  }
  
  .badge {
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.8rem;
    font-weight: bold;
  }
  
  .badge.success {
    background-color: #4caf50;
    color: white;
  }
  
  .badge.warning {
    background-color: #ff9800;
    color: white;
  }
  
  .badge.error {
    background-color: #f44336;
    color: white;
  }
  
  .badge.pending {
    background-color: #9e9e9e;
    color: white;
  }
  
  .loading, .empty-state {
    grid-column: 1 / -1;
    padding: 2rem;
    text-align: center;
    background-color: #f5f5f5;
    border-radius: 8px;
  }
  
  .api-tester {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }
  
  .request-section, .response-section {
    background-color: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    padding: 1rem;
  }
  
  .form-group {
    margin-bottom: 1rem;
  }
  
  .form-row {
    display: flex;
    gap: 1rem;
  }
  
  .form-row .form-group {
    flex: 1;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: bold;
  }
  
  input, select, textarea {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-family: inherit;
    font-size: inherit;
  }
  
  .form-actions {
    margin-top: 1rem;
  }
  
  .response-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    padding: 0.5rem;
    background-color: #f5f5f5;
    border-radius: 4px;
  }
  
  .status-code {
    font-weight: bold;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
  }
  
  .duration {
    font-size: 0.9rem;
    color: #666;
  }
  
  .response-error {
    padding: 1rem;
    background-color: #ffebee;
    color: #f44336;
    border-radius: 4px;
    margin-bottom: 1rem;
  }
  
  .response-headers, .response-body {
    margin-top: 1rem;
  }
  
  .response-headers h4, .response-body h4 {
    margin-top: 0;
    margin-bottom: 0.5rem;
  }
  
  .empty-response {
    padding: 2rem;
    text-align: center;
    color: #666;
    background-color: #f5f5f5;
    border-radius: 4px;
  }
  
  .changes-count {
    background-color: #ff9800;
    color: white;
    padding: 0.1rem 0.4rem;
    border-radius: 10px;
    font-size: 0.8rem;
  }
  
  .btn {
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-weight: bold;
    text-decoration: none;
    display: inline-block;
    text-align: center;
  }
  
  .btn-primary {
    background-color: #2196f3;
    color: white;
  }
  
  .btn-primary:hover {
    background-color: #1976d2;
  }
  
  .btn-secondary {
    background-color: #9e9e9e;
    color: white;
  }
  
  .btn-secondary:hover {
    background-color: #757575;
  }
  
  .btn-danger {
    background-color: #f44336;
    color: white;
  }
  
  .btn-danger:hover {
    background-color: #d32f2f;
  }
  
  .btn:disabled {
    background-color: #ddd;
    color: #888;
    cursor: not-allowed;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
  }
  
  .modal {
    background-color: #fff;
    border-radius: 8px;
    width: 500px;
    max-width: 90%;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
  }
  
  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    border-bottom: 1px solid #eee;
  }
  
  .modal-header h3 {
    margin: 0;
    color: #333;
  }
  
  .modal-body {
    padding: 1rem;
    max-height: 60vh;
    overflow-y: auto;
  }
  
  .confirmation-text {
    color: #333;
    font-size: 1rem;
    margin-bottom: 0.5rem;
  }
  
  .modal-footer {
    padding: 1rem;
    border-top: 1px solid #eee;
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
  }
  
  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #666;
  }
  
  .close-btn:hover {
    color: #000;
  }
  
  .empty-state {
    text-align: center;
    padding: 2rem;
    background-color: #f9f9f9;
    border-radius: 8px;
    margin: 1rem 0;
  }
  
  .empty-state p {
    margin-bottom: 1rem;
    color: #666;
  }
  
  small {
    display: block;
    margin-top: 0.25rem;
    color: #666;
    font-size: 0.8rem;
  }

  .input-with-button {
    display: flex;
    gap: 0.5rem;
  }

  .input-with-button input {
    flex: 1;
  }

  .click-hint {
    font-size: 0.8rem;
    color: #666;
    margin-top: 5px;
    font-style: italic;
  }

  .memory-processes-modal {
    width: 80%;
    max-width: 900px;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #eee;
    padding-bottom: 10px;
    margin-bottom: 15px;
  }

  .close-button {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #666;
  }

  .close-button:hover {
    color: #f44336;
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 30px;
  }

  .modal-footer {
    margin-top: 15px;
    border-top: 1px solid #eee;
    padding-top: 10px;
    text-align: center;
  }

  .memory-note {
    font-size: 0.8rem;
    color: #666;
    font-style: italic;
  }

  .processes-modal {
    width: 80%;
    max-width: 900px;
    color: #333; /* Add explicit text color */
    max-height: 80vh; /* Limit height to 80% of viewport height */
    display: flex;
    flex-direction: column;
  }

  .table-container {
    overflow-y: auto; /* Enable vertical scrolling */
    max-height: 60vh; /* Limit height to 60% of viewport height */
    margin-bottom: 10px;
  }

  .processes-modal table {
    width: 100%;
    margin-bottom: 10px;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    border-bottom: 1px solid #eee;
    padding-bottom: 10px;
    margin-bottom: 15px;
  }

  .close-button {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #666;
  }

  .close-button:hover {
    color: #f44336;
  }

  .loading-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 30px;
  }

  .modal-footer {
    margin-top: 15px;
    border-top: 1px solid #eee;
    padding-top: 10px;
    text-align: center;
  }

  .process-note {
    font-size: 0.8rem;
    color: #666;
    font-style: italic;
    margin: 5px 0;
  }

  .process-note.warning {
    color: #f44336;
  }
  
  /* Add styles for refreshing state */
  .devtools-card.refreshing {
    position: relative;
    pointer-events: none;
    opacity: 0.7;
  }
  
  .devtools-card.refreshing::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1;
  }
  
  .devtools-card.refreshing::before {
    content: "";
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 30px;
    height: 30px;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #3498db;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    z-index: 2;
  }
  
  @keyframes spin {
    0% { transform: translate(-50%, -50%) rotate(0deg); }
    100% { transform: translate(-50%, -50%) rotate(360deg); }
  }
</style>