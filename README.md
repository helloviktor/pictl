# pictl - Raspberry Pi Controller

A lightweight MVC application with a Go backend and minimalistic web interface for managing Raspberry Pi devices. Monitor system metrics and control your device remotely.

## Features

- **System Monitoring**: Real-time display of CPU usage, memory usage, disk space, and CPU temperature
- **Package Management**: Check available system updates
- **System Control**: Restart or shutdown the device via web interface
- **Responsive Design**: Modern web interface that works on all devices
- **Embedded Assets**: HTML, CSS, and JavaScript bundled directly into the Go binary

## Project Structure

```
pictl/
├── cmd/                       # Command-line tools
├── internal/
│   ├── handlers/             # HTTP request handlers
│   ├── models/               # Data structures
│   └── services/             # Business logic
├── web/
│   └── static/               # Static assets (HTML, CSS, JS)
├── main.go                   # Application entry point
├── go.mod                    # Go module definition
├── Makefile                  # Build and run targets
└── README.md                 # This file
```

## Requirements

- Go 1.21 or later
- Linux OS (Debian-based for apt commands)
- Administrative privileges for system control (restart/shutdown)

## Installation

### Clone the repository

```bash
cd /workspaces/pictl
```

### Install dependencies

```bash
go mod download
```

## Building

### Using Make

```bash
make build
```

### Using Go directly

```bash
go build -o pictl .
```

## Running

### Using Make

```bash
make run
```

### Using the binary directly

```bash
./pictl --addr :8080
```

The application will start on `http://localhost:8080`

## Usage

1. Open your browser and navigate to `http://localhost:8080`
2. View real-time system information
3. Use the buttons to update system packages, restart, or shutdown

## API Endpoints

### GET `/api/system/info`
Retrieves current system information (CPU, memory, disk, temperature, updates available)

**Response:**
```json
{
  "cpu_usage": 25.5,
  "memory_usage": 45.2,
  "disk_usage": 62.8,
  "cpu_temp": 52.3,
  "last_update": "2024-01-15 14:30:45",
  "updates_available": 3
}
```

### POST `/api/system/update`
Starts system package update

**Response:**
```json
{
  "success": true,
  "message": "System update started"
}
```

### POST `/api/system/restart`
Initiates system restart

**Response:**
```json
{
  "success": true,
  "message": "System restart initiated"
}
```

### POST `/api/system/shutdown`
Initiates system shutdown

**Response:**
```json
{
  "success": true,
  "message": "System shutdown initiated"
}
```

## System Information Sources

The application retrieves information from standard Linux system files:

- **CPU Usage**: `/proc/stat`
- **Memory Usage**: `/proc/meminfo`
- **Disk Usage**: `df` command
- **CPU Temperature**: `/sys/class/thermal/thermal_zone0/temp`
- **Available Updates**: `apt list --upgradable`

## Graceful Shutdown

The application listens for OS signals (SIGINT, SIGTERM) to ensure proper cleanup before shutdown.

## License

MIT
