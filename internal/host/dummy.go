//go:build !rpi

package host

// Dummy is a Host implementation that returns stub values without touching the
// host system. For local development and testing purposes only.
type Dummy struct{}

// NewHost creates a new dummy host.
func NewHost() *Dummy {
	return &Dummy{}
}

// CPUUsage returns a stub CPU usage percentage.
func (d *Dummy) CPUUsage() (float64, error) {
	return 12.3, nil
}

// MemoryUsage returns a stub memory usage percentage.
func (d *Dummy) MemoryUsage() (float64, error) {
	return 34.5, nil
}

// DiskUsage returns a stub disk usage percentage.
func (d *Dummy) DiskUsage() (float64, error) {
	return 45.6, nil
}

// CPUTemperature returns a stub CPU temperature in Celsius.
func (d *Dummy) CPUTemperature() (float64, error) {
	return 42.0, nil
}

// AvailableUpdates returns a stub number of available package updates.
func (d *Dummy) AvailableUpdates() (int, error) {
	return 3, nil
}

// ApplyUpdates is a no-op.
func (d *Dummy) ApplyUpdates() error {
	return nil
}

// Restart is a no-op.
func (d *Dummy) Restart() error {
	return nil
}

// Shutdown is a no-op.
func (d *Dummy) Shutdown() error {
	return nil
}
