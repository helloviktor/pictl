package services

import (
	"time"

	"github.com/pictl/pictl/internal/host"
	"github.com/pictl/pictl/internal/models"
)

// SystemService provides system information and control by interacting with the host directly.
type SystemService struct {
	host host.Host
}

// NewSystemService creates a new SystemService backed by the given Host
func NewSystemService(h host.Host) *SystemService {
	return &SystemService{host: h}
}

// SystemInfo gathers current system metrics from the host
func (s *SystemService) SystemInfo() (models.SystemInfo, error) {
	cpuUsage, err := s.host.CPUUsage()
	if err != nil {
		return models.SystemInfo{}, err
	}

	memoryUsage, err := s.host.MemoryUsage()
	if err != nil {
		return models.SystemInfo{}, err
	}

	diskUsage, err := s.host.DiskUsage()
	if err != nil {
		return models.SystemInfo{}, err
	}

	cpuTemp, err := s.host.CPUTemperature()
	if err != nil {
		return models.SystemInfo{}, err
	}

	updatesAvail, err := s.host.AvailableUpdates()
	if err != nil {
		return models.SystemInfo{}, err
	}

	return models.SystemInfo{
		CPUUsage:     cpuUsage,
		MemoryUsage:  memoryUsage,
		DiskUsage:    diskUsage,
		CPUTemp:      cpuTemp,
		LastUpdate:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatesAvail: updatesAvail,
	}, nil
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
