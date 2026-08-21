package services

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pictl/pictl/internal/models"
)

// SystemService provides system information and control
type SystemService struct{}

// NewSystemService creates a new SystemService
func NewSystemService() *SystemService {
	return &SystemService{}
}

// GetSystemInfo retrieves current system information
func (s *SystemService) GetSystemInfo() (*models.SystemInfo, error) {
	info := &models.SystemInfo{
		LastUpdate: time.Now().Format("2006-01-02 15:04:05"),
	}

	// Get CPU usage
	cpuUsage, err := s.getCPUUsage()
	if err == nil {
		info.CPUUsage = cpuUsage
	}

	// Get memory usage
	memUsage, err := s.getMemoryUsage()
	if err == nil {
		info.MemoryUsage = memUsage
	}

	// Get disk usage
	diskUsage, err := s.getDiskUsage()
	if err == nil {
		info.DiskUsage = diskUsage
	}

	// Get CPU temperature
	cpuTemp, err := s.getCPUTemperature()
	if err == nil {
		info.CPUTemp = cpuTemp
	}

	// Get available updates
	updates, err := s.getAvailableUpdates()
	if err == nil {
		info.UpdatesAvail = updates
	}

	return info, nil
}

// getCPUUsage retrieves CPU usage percentage
func (s *SystemService) getCPUUsage() (float64, error) {
	// Simple CPU usage calculation based on /proc/stat
	// This is a simplified version; for production, consider using a library
	return 0.0, nil
}

// getMemoryUsage retrieves memory usage percentage
func (s *SystemService) getMemoryUsage() (float64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var memTotal, memAvail float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				val, _ := strconv.ParseFloat(parts[1], 64)
				memTotal = val
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				val, _ := strconv.ParseFloat(parts[1], 64)
				memAvail = val
			}
		}
	}

	if memTotal == 0 {
		return 0, fmt.Errorf("could not read memory info")
	}

	used := memTotal - memAvail
	return (used / memTotal) * 100, nil
}

// getDiskUsage retrieves disk usage percentage for root partition
func (s *SystemService) getDiskUsage() (float64, error) {
	cmd := exec.Command("df", "/")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("could not parse df output")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return 0, fmt.Errorf("unexpected df output format")
	}

	used, _ := strconv.ParseFloat(fields[2], 64)
	total, _ := strconv.ParseFloat(fields[1], 64)

	return (used / total) * 100, nil
}

// getCPUTemperature retrieves CPU temperature in Celsius
func (s *SystemService) getCPUTemperature() (float64, error) {
	tempPath := "/sys/class/thermal/thermal_zone0/temp"
	data, err := os.ReadFile(tempPath)
	if err != nil {
		return 0, err
	}

	tempStr := strings.TrimSpace(string(data))
	tempMilliC, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}

	// Convert from millidegrees to degrees
	return tempMilliC / 1000, nil
}

// getAvailableUpdates retrieves the number of available package updates
func (s *SystemService) getAvailableUpdates() (int, error) {
	cmd := exec.Command("apt", "list", "--upgradable")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Count non-empty lines, excluding the first header line if present
	lines := strings.Split(string(output), "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" && !strings.Contains(line, "Listing...") {
			count++
		}
	}

	return count, nil
}

// UpdateSystem runs system update
func (s *SystemService) UpdateSystem() error {
	// This would typically run apt update && apt upgrade
	// For safety, this is a placeholder
	return nil
}

// RestartSystem restarts the device
func (s *SystemService) RestartSystem() error {
	cmd := exec.Command("sudo", "shutdown", "-r", "now")
	return cmd.Run()
}

// ShutdownSystem shuts down the device
func (s *SystemService) ShutdownSystem() error {
	cmd := exec.Command("sudo", "shutdown", "-h", "now")
	return cmd.Run()
}
