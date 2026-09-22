> [!WARNING]
> **Work in progress:** the project is under active development. Features, APIs, deployment instructions, and service configuration may change.

# pictl - Raspberry Pi Controller

A lightweight Go application and web interface for monitoring and managing a Raspberry Pi device.

The project uses two binaries to keep privileged operations isolated. The `dashboard` web server runs as an unprivileged user and sends requests over a Unix socket to the `pictl` helper, which performs host-level operations.

## Features

- **System Monitoring**: Display CPU usage, memory usage, disk space, and CPU temperature
- **Package Management**: Check for and install available system updates
- **System Control**: Restart or shutdown the device via web interface
- **Privilege Separation**: An unprivileged web server delegates host operations to a root-owned helper over a Unix socket
- **Embedded Assets**: HTML, CSS, and JavaScript are bundled directly into the dashboard binary

## Project Structure

```
pictl/
├── cmd/
│   ├── dashboard/             # Unprivileged HTTP server entry point
│   │   └── main.go
│   └── pictl/                 # Privileged Unix-socket helper entry point
│       └── main.go
├── configs/                   # systemd service and socket unit files
│   ├── dashboard.service
│   ├── pictl.service
│   └── pictl.socket
├── internal/
│   ├── handlers/              # HTTP and static-file handlers
│   ├── host/                  # Linux host metrics and control operations
│   ├── ipc/                   # JSON-over-Unix-socket client and protocol
│   ├── models/                # API response models
│   └── services/              # Dashboard service backed by the IPC client
├── scripts/
│   └── deploy.sh              # arm64 build and rsync deployment script
├── web/
│   ├── embed.go               # Embeds the dashboard assets
│   └── static/                # Single-page HTML, CSS, and JavaScript
├── go.mod                     # Go module definition
├── Makefile                   # Build, test, lint, and deployment targets
└── README.md
```

## Requirements

- Go 1.21 or later
- Linux OS (Debian-based for apt commands)
- `systemd` for the provided production service configuration
- A `pictl` user and group for the provided service configuration

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

This produces both binaries in `out/`:

- `out/dashboard`: HTTP server, defaults to `:8080` and `/run/pictl.sock`
- `out/pictl`: privileged helper, launched with systemd socket activation or `--socket <path>`

Use `GOARCH` to cross-compile and `TAGS` to enable build tags:

```bash
GOARCH=arm64 TAGS=rpi make build
```

### Using Go directly

```bash
go build -o out/dashboard ./cmd/dashboard
go build -o out/pictl ./cmd/pictl
```

## Running

`make run` builds and starts the dashboard. It expects the helper socket at `/run/pictl.sock` to already be available.

```bash
make run
```

### Local development

Start the helper on a temporary socket in one terminal:

```bash
./out/pictl --socket /tmp/pictl.sock
```

Then start the dashboard in another terminal:

```bash
./out/dashboard --addr :8080 --socket /tmp/pictl.sock
```

Open `http://localhost:8080`.

### systemd deployment

The `configs/` directory contains the production units. `pictl.socket` creates `/run/pictl.sock` with group access for `pictl`; `pictl.service` is socket-activated and runs as root; `dashboard.service` runs as the unprivileged `pictl` user.

Install the binaries under `/usr/local/bin`, copy the unit files to `/etc/systemd/system`, reload systemd, and enable the socket and dashboard service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now pictl.socket dashboard.service
```

## Usage

1. Open your browser and navigate to `http://localhost:8080`
2. View real-time system information
3. Use the buttons to update system packages, restart, or shutdown

## API Endpoints

The dashboard exposes the following routes. Metric endpoints return a JSON number; action endpoints return a JSON object with `success` and `message` fields.

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/api/system/cpu-usage` | Current CPU usage percentage |
| `GET` | `/api/system/memory-usage` | Current memory usage percentage |
| `GET` | `/api/system/disk-usage` | Root filesystem usage percentage |
| `GET` | `/api/system/cpu-temperature` | Current CPU temperature in Celsius |
| `GET` | `/api/system/available-updates` | Number of available package updates |
| `POST` | `/api/system/update` | Install available package updates |
| `POST` | `/api/system/restart` | Restart the device |
| `POST` | `/api/system/shutdown` | Shut down the device |

Example action response:

```json
{
  "success": true,
  "message": "System update started"
}
```

## System Information Sources

The application retrieves information from standard Linux system files:

- **CPU Usage**: `/proc/stat`
- **Memory Usage**: `/proc/meminfo`
- **Disk Usage**: Linux filesystem statistics
- **CPU Temperature**: `/sys/class/thermal/thermal_zone0/temp`
- **Available Updates and Installation**: `apt`

## Development

```bash
make test
make fmt
make lint
```

`make deploy` builds arm64 binaries with the `rpi` build tag by default and deploys them with `rsync`. It reads `DEPLOY_HOST`, `DEPLOY_USER`, and optional deployment settings from `.env`.

## Graceful Shutdown

The application listens for OS signals (SIGINT, SIGTERM) to ensure proper cleanup before shutdown.

## License

MIT
