package services

import (
	"encoding/json"
	"fmt"

	"github.com/pictl/pictl/internal/ipc"
	"github.com/pictl/pictl/internal/models"
)

// RemoteSystemService provides system information and control by delegating
// to the privileged pictl helper process over a Unix socket, instead of
// interacting with the host directly.
type RemoteSystemService struct {
	client *ipc.Client
}

// NewRemoteSystemService connects to the pictl helper socket at socketPath.
func NewRemoteSystemService(socketPath string) (*RemoteSystemService, error) {
	client, err := ipc.NewClient(socketPath)
	if err != nil {
		return nil, fmt.Errorf("connect to pictl socket: %w", err)
	}
	return &RemoteSystemService{client: client}, nil
}

// Close closes the underlying connection to the pictl helper process.
func (s *RemoteSystemService) Close() error {
	return s.client.Close()
}

// SystemInfo retrieves current system information from the pictl helper process
func (s *RemoteSystemService) SystemInfo() (models.SystemInfo, error) {
	var info models.SystemInfo
	if err := s.call("info", &info); err != nil {
		return models.SystemInfo{}, err
	}
	return info, nil
}

// UpdateSystem requests a system update from the pictl helper process
func (s *RemoteSystemService) UpdateSystem() error {
	return s.call("update", nil)
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
