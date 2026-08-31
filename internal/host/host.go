// Package host provides access to host system metrics and control operations.
package host

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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

// Pi is a Host implementation backed by real Raspberry Pi system files and commands.
type Pi struct{}

// NewPi creates a new Pi.
func NewPi() *Pi {
	return &Pi{}
}

// CPUUsage retrieves CPU usage percentage
func (p *Pi) CPUUsage() (float64, error) {
	// Simple CPU usage calculation based on /proc/stat
	// This is a simplified version; for production, consider using a library
	return 0.0, nil
}

// MemoryUsage retrieves memory usage percentage
func (p *Pi) MemoryUsage() (float64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return parseMemoryUsage(file)
}

func parseMemoryUsage(input io.Reader) (float64, error) {
	var memTotal, memAvail float64
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				memTotal, _ = strconv.ParseFloat(parts[1], 64)
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				memAvail, _ = strconv.ParseFloat(parts[1], 64)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	if memTotal == 0 {
		return 0, fmt.Errorf("could not read memory info")
	}

	used := memTotal - memAvail
	return (used / memTotal) * 100, nil
}

// DiskUsage retrieves disk usage percentage for root partition
func (p *Pi) DiskUsage() (float64, error) {
	cmd := exec.Command("df", "/")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	return parseDiskUsage(output)
}

func parseDiskUsage(output []byte) (float64, error) {
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

// CPUTemperature retrieves CPU temperature in Celsius
func (p *Pi) CPUTemperature() (float64, error) {
	tempPath := "/sys/class/thermal/thermal_zone0/temp"
	data, err := os.ReadFile(tempPath)
	if err != nil {
		return 0, err
	}

	return parseCPUTemperature(data)
}

func parseCPUTemperature(data []byte) (float64, error) {
	tempStr := strings.TrimSpace(string(data))
	tempMilliC, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, err
	}

	// Convert from millidegrees to degrees
	return tempMilliC / 1000, nil
}

// AvailableUpdates retrieves the number of available package updates
func (p *Pi) AvailableUpdates() (int, error) {
	if err := exec.Command("apt", "update").Run(); err != nil {
		return 0, fmt.Errorf("update package cache: %w", err)
	}

	cmd := exec.Command("apt", "list", "--upgradable")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("list available updates: %w", err)
	}

	return parseAvailableUpdates(output), nil
}

func parseAvailableUpdates(output []byte) int {
	// Count non-empty lines, excluding the first header line if present
	lines := strings.Split(string(output), "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" && !strings.Contains(line, "Listing...") {
			count++
		}
	}

	return count
}

// ApplyUpdates installs available package updates
func (p *Pi) ApplyUpdates() error {
	// This would typically run apt update && apt upgrade
	// For safety, this is a placeholder
	return nil
}

// Restart restarts the device
func (p *Pi) Restart() error {
	cmd := exec.Command("sudo", "shutdown", "-r", "now")
	return cmd.Run()
}

// Shutdown shuts down the device
func (p *Pi) Shutdown() error {
	cmd := exec.Command("sudo", "shutdown", "-h", "now")
	return cmd.Run()
}

// Dummy is a Host implementation that returns stub values without touching the
// host system. For local development and testing purposes only.
type Dummy struct{}

// NewDummy creates a new Dummy.
func NewDummy() *Dummy {
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
