// Update system information
async function updateSystemInfo() {
    try {
        const response = await fetch('/api/system/info');
        const data = await response.json();

        // Update UI with system info
        document.getElementById('cpu-usage').textContent = data.cpu_usage.toFixed(1) + '%';
        document.getElementById('cpu-progress').style.width = data.cpu_usage + '%';

        document.getElementById('memory-usage').textContent = data.memory_usage.toFixed(1) + '%';
        document.getElementById('memory-progress').style.width = data.memory_usage + '%';

        document.getElementById('disk-usage').textContent = data.disk_usage.toFixed(1) + '%';
        document.getElementById('disk-progress').style.width = data.disk_usage + '%';

        document.getElementById('cpu-temp').textContent = data.cpu_temp.toFixed(1) + '°C';

        document.getElementById('last-update').textContent = data.last_update;
        document.getElementById('updates-available').textContent = data.updates_available;
    } catch (error) {
        console.error('Error fetching system info:', error);
    }
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
