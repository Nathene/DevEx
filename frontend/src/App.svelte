<script>
  import { onMount } from 'svelte';
  import { 
    GetCPUInfo, 
    GetCPUDetails,
    GetRAMInfo, 
    GetRAMDetails,
    GetDiskInfo,
    GetDiskDetails
  } from '../wailsjs/go/main/App.js';

  let cpuInfo = "Loading...";
  let ramInfo = "Loading...";
  let diskInfo = "Loading...";

  // Helper function to extract percentage from info string
  function getPercentage(info) {
    const match = info.match(/(\d+\.?\d*)/);
    return match ? parseFloat(match[0]) : 0;
  }

  // Helper function to determine status class based on percentage
  function getStatusClass(percentage) {
    if (percentage >= 80) return 'critical';
    if (percentage >= 60) return 'warning';
    return 'normal';
  }

  onMount(() => {
    const updateInfo = async () => {
      cpuInfo = await GetCPUInfo();
      ramInfo = await GetRAMInfo();
      diskInfo = await GetDiskInfo();
    };

    updateInfo();
    const interval = setInterval(updateInfo, 2000);
    return () => clearInterval(interval);
  });
</script>

<main class="app-container">
  <nav class="navbar">
    <button class="nav-btn">Home</button>
    <button class="nav-btn">About</button>
    <button class="nav-btn">Contact</button>
  </nav>
  <div class="content">
    <h3>System Information</h3>
    <div class="system-info-container">
      <div class="system-info">
        <div class="info-bubble {getStatusClass(getPercentage(cpuInfo))}" 
             title={cpuInfo}>
          {cpuInfo}
        </div>
        <div class="info-bubble {getStatusClass(getPercentage(ramInfo))}" 
             title={ramInfo}>
          {ramInfo}
        </div>
        <div class="info-bubble {getStatusClass(getPercentage(diskInfo))}" 
             title={diskInfo}>
          {diskInfo}
        </div>
      </div>
    </div>
    <h3>Network Information</h3>
    <div class="system-info-container">
      <div class="system-info">
        <div class="info-bubble">Network Info</div>
        <div class="info-bubble">Network Info</div>
        <div class="info-bubble">Network Info</div>
      </div>
    </div>
      <h3>Applications</h3>
      <div class="system-info-container">
        <div class="system-info">
          <div class="info-bubble">App Info</div>
          <div class="info-bubble">App Info</div>
          <div class="info-bubble">App Info</div>
        </div>
  </div>
</main>

<style>
  .info-bubble {
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    text-align: center;
    min-width: 100px;
    flex: 1;
    cursor: help;
    transition: all 0.2s ease;
    color: white;
  }

  .info-bubble:hover {
    transform: scale(1.05);
  }

  .info-bubble.normal {
    background-color: #007bff;
  }

  .info-bubble.warning {
    background-color: #ffa500;
  }

  .info-bubble.critical {
    background-color: #ff4444;
  }


  .app-container {
    display: flex;
    height: 100vh;
    width: 100vw;
    margin: 0;
    padding: 0;
  }

  .navbar {
    width: 200px;
    height: 100vh;
    background-color: #333;
    display: flex;
    flex-direction: column;
    padding: 10px;
    position: fixed;
    left: 0;
    top: 0;
  }

  .nav-btn {
    padding: 10px;
    margin: 5px 0;
    border: none;
    border-radius: 5px;
    background-color: #007bff;
    color: white;
    cursor: pointer;
  }

  .nav-btn:hover {
    background-color: #0056b3;
  }

  .content {
    flex: 1;
    margin-left: 220px;
    padding: 20px;
    overflow-y: auto;
  }

  .system-info-container {
    background-color: #f0f0f0;
    border-radius: 10px;
    padding: 20px;
    margin-bottom: 20px;
    width: calc(100% - 40px);
  }

  .system-info {
    display: flex;
    justify-content: space-around;
    align-items: center;
    gap: 20px;
  }

  .info-bubble {
    background-color: #007bff;
    color: white;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    text-align: center;
    min-width: 100px;
    flex: 1;
  }
</style>