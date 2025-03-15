<script>
  import { GetCPUInfo, GetCPUDetails, GetRAMInfo, GetRAMDetails, GetDiskInfo, GetDiskDetails, GetDockerStatus, GetDockerMetrics, GetNetworkStatus } from '../wailsjs/go/main/App';
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
    { id: 'processes', name: 'Processes', icon: '⚙️' },
    { id: 'logs', name: 'Logs', icon: '📝' },
    { id: 'settings', name: 'Settings', icon: '⚙️' }
  ];

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

  onMount(() => {
    updateMetrics();
    const interval = setInterval(updateMetrics, 2000);
    return () => clearInterval(interval);
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
    {:else if selectedCategory === 'processes'}
      <div class="dashboard-grid">
        <div class="card">
          <h2>Process List</h2>
          <div class="metric">Coming Soon</div>
          <div class="details">Process monitoring will be implemented in a future update</div>
        </div>
        <div class="card">
          <h2>Process Details</h2>
          <div class="metric">Coming Soon</div>
          <div class="details">Detailed process information will be available here</div>
        </div>
        <div class="card">
          <h2>Process Stats</h2>
          <div class="metric">Coming Soon</div>
          <div class="details">Process statistics and performance metrics</div>
        </div>
      </div>
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
</style>