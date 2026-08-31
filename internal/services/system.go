package services

import "github.com/pictl/pictl/internal/host"

// SystemService provides system information and control
type SystemService struct {
	host host.Host
}

// NewSystemService creates a new SystemService backed by the given Host
func NewSystemService(h host.Host) *SystemService {
	return &SystemService{host: h}
}

// CPUUsage returns the current CPU usage percentage
func (s *SystemService) CPUUsage() (float64, error) {
	return s.host.CPUUsage()
}

// MemoryUsage returns the current memory usage percentage
func (s *SystemService) MemoryUsage() (float64, error) {
	return s.host.MemoryUsage()
}

// DiskUsage returns the current disk usage percentage for the root partition
func (s *SystemService) DiskUsage() (float64, error) {
	return s.host.DiskUsage()
}

// CPUTemperature returns the current CPU temperature in Celsius
func (s *SystemService) CPUTemperature() (float64, error) {
	return s.host.CPUTemperature()
}

// AvailableUpdates returns the number of available package updates
func (s *SystemService) AvailableUpdates() (int, error) {
	return s.host.AvailableUpdates()
}

// UpdateSystem installs available package updates
func (s *SystemService) UpdateSystem() error {
	return s.host.ApplyUpdates()
}

// RestartSystem restarts the device
func (s *SystemService) RestartSystem() error {
	return s.host.Restart()
}

// ShutdownSystem shuts down the device
func (s *SystemService) ShutdownSystem() error {
	return s.host.Shutdown()
}
