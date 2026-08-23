# Project context

This is a lightweight application with a Go backend and a minimalistic frontend with HTML, CSS, and JavaScript.
The application is designed to manage a Raspberry Pi device, allowing users to monitor and control various aspects of the device through a web interface.

# Features

- Display system information, including CPU usage, memory usage, disk space, CPU temperature.
- Check for available system and package updates.
- Update the system and installed packages.
- Restart or shut down the Raspberry Pi.

# Architecture and constraints

The application must consist of two Go binaries:

- Web application: provides the web interface, RESTful API, and application logic.
- Helper application: performs operations that require elevated privileges, such as restarting or shutting down the Raspberry Pi, and updating installed packages.

## Directory structure

```
project-root/
├── cmd/
│   ├── helper/
│   ├── webapp/
├── internal/
├── web/
│   ├── static/
```

# Security constraints

- The web application must always run as a non-root user.
- The web application must not be granted additional privileges to perform operations that belong in the helper application.
- Privileged operations must be delegated to the helper application.
- The web application must communicate with the helper application through a Unix socket.
- Avoid shell execution where possible; when commands must be executed, use explicit argument lists rather than constructing shell commands from user inputs.

# Technical Details

- System information must be retrieved from the Linux `/sys` filesystem, e.g. CPU temperature is obtained from `/sys/class/thermal/thermal_zone0/temp`.
- The application must periodically run `apt update` to keep the package lists up to date and determine whether package updates are available.
- The web interface must communicate with the Go backend through RESTful APIs.
- The Go backend must listen for relevant operating system signals and perform a graceful shutdown when necessary.
- The HTML files must be embedded into the Go binary using Go's `embed` package.

# Web Interface Requirements

The web interface must consist of a single page.

## System Information

The interface must display:

- CPU usage
- Memory usage
- Disk usage
- CPU temperature
- Last update time
- Number of available updates

## Actions

The interface must provide buttons for:

- Updating the system.
- Restarting the Raspberry Pi.
- Shutting down the Raspberry Pi.


# API Requirements

The backend must provide RESTful API endpoints that allow the web interface to:

- Retrieve the current system information.
- Retrieve the current update status.
- Request a system update.
- Request a device restart.
- Request a device shutdown.

# Development Guidelines

- Preserve the separation between the unprivileged web application and the privileged helper application.
- Keep the frontend minimal and use plain HTML, CSS, and JavaScript unless there is a strong reason to introduce a dependency.
- Keep changes focused on the requested functionality.
- When changing API behavior, update both the backend and frontend as necessary.
- Handle errors explicitly and return useful error responses from API endpoints.
- Add or update tests when changing backend behavior.
