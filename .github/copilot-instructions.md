# Project description

This is a lightweight MVC application with a Go backend and a minimalistic frontend with HTML, CSS, and JavaScript. The application is designed to manage a Raspberry Pi device, allowing users to monitor and control various aspects of the device through a web interface.

# Features

- Displaying system information (CPU usage, memory usage, disk space, CPU temperature, etc.).
- Updating the system and installed packages.
- Restarting and shutting down the device.

# Premises

- System information is retrieved from the `/sys` directory, such as `/sys/class/thermal/thermal_zone0/temp` for CPU temperature.
- The application runs `apt update` periodically to keep the package list up to date.
- The web interface communicates with the backend using RESTful APIs.
- The application listens for OS signal events to handle graceful shutdowns and restarts.
- The HTML files are embedded into the Go binary using the `embed` package.


# Web Interface Requirements

- The web interface consists of a single page.
- The top section displays CPU, memory, disk usage, CPU temperature, last update time, and the number of updates available.
- There are buttons for updating the system, restarting, and shutting down the device.
