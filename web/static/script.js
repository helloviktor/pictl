// Update system information
function updateSystemInfo() {
    fetch('/api/system/cpu-usage')
        .then(r => r.json())
        .then(cpuUsage => {
            document.getElementById('cpu-usage').textContent = cpuUsage.toFixed(1) + '%';
            document.getElementById('cpu-progress').style.width = cpuUsage + '%';
        })
        .catch(error => console.error('Error fetching CPU usage:', error));

    fetch('/api/system/memory-usage')
        .then(r => r.json())
        .then(memoryUsage => {
            document.getElementById('memory-usage').textContent = memoryUsage.toFixed(1) + '%';
            document.getElementById('memory-progress').style.width = memoryUsage + '%';
        })
        .catch(error => console.error('Error fetching memory usage:', error));

    fetch('/api/system/disk-usage')
        .then(r => r.json())
        .then(diskUsage => {
            document.getElementById('disk-usage').textContent = diskUsage.toFixed(1) + '%';
            document.getElementById('disk-progress').style.width = diskUsage + '%';
        })
        .catch(error => console.error('Error fetching disk usage:', error));

    fetch('/api/system/cpu-temperature')
        .then(r => r.json())
        .then(cpuTemp => {
            document.getElementById('cpu-temp').textContent = cpuTemp.toFixed(1) + '°C';
        })
        .catch(error => console.error('Error fetching CPU temperature:', error));

    fetch('/api/system/available-updates')
        .then(r => r.json())
        .then(updatesAvailable => {
            document.getElementById('updates-available').textContent = updatesAvailable;
        })
        .catch(error => console.error('Error fetching available updates:', error));

    document.getElementById('last-update').textContent = new Date().toLocaleString();
}

// Handle system update
document.getElementById('btn-update').addEventListener('click', async function() {
    if (confirm('Start system update?')) {
        this.disabled = true;
        try {
            const response = await fetch('/api/system/update', { method: 'POST' });
            const data = await response.json();
            alert(data.message);
        } catch (error) {
            alert('Error: ' + error.message);
        } finally {
            this.disabled = false;
        }
    }
});

// Handle system restart
document.getElementById('btn-restart').addEventListener('click', async function() {
    if (confirm('Restart the device? This will take about 30 seconds.')) {
        this.disabled = true;
        try {
            const response = await fetch('/api/system/restart', { method: 'POST' });
            const data = await response.json();
            alert(data.message);
        } catch (error) {
            alert('Error: ' + error.message);
        } finally {
            this.disabled = false;
        }
    }
});

// Handle system shutdown
document.getElementById('btn-shutdown').addEventListener('click', async function() {
    if (confirm('Shutdown the device? Make sure to power it back on via hardware.')) {
        this.disabled = true;
        try {
            const response = await fetch('/api/system/shutdown', { method: 'POST' });
            const data = await response.json();
            alert(data.message);
        } catch (error) {
            alert('Error: ' + error.message);
        } finally {
            this.disabled = false;
        }
    }
});

// Update system info every 5 seconds
updateSystemInfo();
setInterval(updateSystemInfo, 5000);
