// Package host provides access to host system metrics and control operations.
package host

// Host provides system metrics and control operations for the device pictl manages.
type Host interface {
	// CPUUsage returns the current CPU usage percentage.
	CPUUsage() (float64, error)
	// MemoryUsage returns the current memory usage percentage.
	MemoryUsage() (float64, error)
	// DiskUsage returns the current disk usage percentage for the root partition.
	DiskUsage() (float64, error)
	// CPUTemperature returns the current CPU temperature in Celsius.
	CPUTemperature() (float64, error)
	// AvailableUpdates returns the number of available package updates.
	AvailableUpdates() (int, error)
	// ApplyUpdates installs available package updates.
	ApplyUpdates() error
	// Restart restarts the device.
	Restart() error
	// Shutdown shuts down the device.
	Shutdown() error
}
