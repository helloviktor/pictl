package services

import (
	"encoding/json"
	"fmt"
)

// IpcClient is the subset of ipc.Client that RemoteSystemService depends on.
type IpcClient interface {
	SendCommand(command any) (any, error)
	Close() error
}

// RemoteSystemService provides system information and control by delegating
// to the privileged pictl helper process over a Unix socket, instead of
// interacting with the host directly.
type RemoteSystemService struct {
	client IpcClient
}

// NewRemoteSystemService creates a RemoteSystemService backed by the given IpcClient.
func NewRemoteSystemService(client IpcClient) *RemoteSystemService {
	return &RemoteSystemService{client: client}
}

// CPUUsage retrieves the current CPU usage percentage from the pictl helper process
func (s *RemoteSystemService) CPUUsage() (float64, error) {
	var value float64
	if err := s.call("cpu_usage", &value); err != nil {
		return 0, err
	}
	return value, nil
}

// MemoryUsage retrieves the current memory usage percentage from the pictl helper process
func (s *RemoteSystemService) MemoryUsage() (float64, error) {
	var value float64
	if err := s.call("memory_usage", &value); err != nil {
		return 0, err
	}
	return value, nil
}

// DiskUsage retrieves the current disk usage percentage from the pictl helper process
func (s *RemoteSystemService) DiskUsage() (float64, error) {
	var value float64
	if err := s.call("disk_usage", &value); err != nil {
		return 0, err
	}
	return value, nil
}

// CPUTemperature retrieves the current CPU temperature in Celsius from the pictl helper process
func (s *RemoteSystemService) CPUTemperature() (float64, error) {
	var value float64
	if err := s.call("cpu_temperature", &value); err != nil {
		return 0, err
	}
	return value, nil
}

// AvailableUpdates retrieves the number of available package updates from the pictl helper process
func (s *RemoteSystemService) AvailableUpdates() (int, error) {
	var value int
	if err := s.call("available_updates", &value); err != nil {
		return 0, err
	}
	return value, nil
}

// UpdateSystem requests a system update from the pictl helper process
func (s *RemoteSystemService) UpdateSystem() error {
	return s.call("apply_updates", nil)
}

// RestartSystem requests a device restart from the pictl helper process
func (s *RemoteSystemService) RestartSystem() error {
	return s.call("restart", nil)
}

// ShutdownSystem requests a device shutdown from the pictl helper process
func (s *RemoteSystemService) ShutdownSystem() error {
	return s.call("shutdown", nil)
}

// call sends command to the pictl helper process and, if out is non-nil, decodes the result into it
func (s *RemoteSystemService) call(command string, out any) error {
	result, err := s.client.SendCommand(command)
	if err != nil {
		return err
	}
	if out == nil {
		return nil
	}

	// The client decodes results generically as any; round-trip through JSON to
	// populate the caller's concrete type.
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}
	return json.Unmarshal(data, out)
}
