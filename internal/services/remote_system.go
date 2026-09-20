package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/helloviktor/pictl/internal/ipc"
)

// IpcClient is the subset of ipc.Client that RemoteSystemService depends on.
type IpcClient interface {
	SendCommand(ctx context.Context, command ipc.Command) (any, error)
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
	if err := s.call(ipc.CommandCPUUsage, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// MemoryUsage retrieves the current memory usage percentage from the pictl helper process
func (s *RemoteSystemService) MemoryUsage() (float64, error) {
	var value float64
	if err := s.call(ipc.CommandMemoryUsage, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// DiskUsage retrieves the current disk usage percentage from the pictl helper process
func (s *RemoteSystemService) DiskUsage() (float64, error) {
	var value float64
	if err := s.call(ipc.CommandDiskUsage, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// CPUTemperature retrieves the current CPU temperature in Celsius from the pictl helper process
func (s *RemoteSystemService) CPUTemperature() (float64, error) {
	var value float64
	if err := s.call(ipc.CommandCPUTemperature, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// AvailableUpdates retrieves the number of available package updates from the pictl helper process
func (s *RemoteSystemService) AvailableUpdates() (int, error) {
	var value int
	if err := s.call(ipc.CommandAvailableUpdates, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// UpdateSystem requests a system update from the pictl helper process
func (s *RemoteSystemService) UpdateSystem() error {
	return s.call(ipc.CommandApplyUpdates, nil)
}

// RestartSystem requests a device restart from the pictl helper process
func (s *RemoteSystemService) RestartSystem() error {
	return s.call(ipc.CommandRestart, nil)
}

// ShutdownSystem requests a device shutdown from the pictl helper process
func (s *RemoteSystemService) ShutdownSystem() error {
	return s.call(ipc.CommandShutdown, nil)
}

// call sends command to the pictl helper process and, if out is non-nil, decodes the result into it
func (s *RemoteSystemService) call(command ipc.Command, out any) error {
	result, err := s.client.SendCommand(context.Background(), command)
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
