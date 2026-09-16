package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"sync"
	"testing"
)

func startMockServer(t *testing.T, handler func(req Request) Response) (string, func()) {
	t.Helper()

	socketPath := filepath.Join(t.TempDir(), "test.sock")
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("failed to listen on unix socket: %v", err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		scanner := bufio.NewScanner(conn)
		encoder := json.NewEncoder(conn)

		for scanner.Scan() {
			var req Request
			if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
				continue
			}

			resp := handler(req)
			_ = encoder.Encode(resp)
		}
	}()

	cleanup := func() {
		listener.Close()
		<-done
	}

	return socketPath, cleanup
}

func TestNewClientError(t *testing.T) {
	nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.sock")
	_, err := NewClient(nonExistentPath)
	if err == nil {
		t.Fatal("expected error connecting to non-existent socket, got nil")
	}
}

func TestSendCommand(t *testing.T) {
	socketPath, cleanup := startMockServer(t, func(req Request) Response {
		return Response{
			Id:     req.Id,
			Result: fmt.Sprintf("echo:%v", req.Command),
		}
	})
	defer cleanup()

	client, err := NewClient(socketPath)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	result, err := client.SendCommand("ping")
	if err != nil {
		t.Fatalf("SendCommand returned an error: %v", err)
	}
	if result != "echo:ping" {
		t.Fatalf("expected 'echo:ping', got %v", result)
	}
}

func TestSendCommandConcurrent(t *testing.T) {
	socketPath, cleanup := startMockServer(t, func(req Request) Response {
		return Response{
			Id:     req.Id,
			Result: req.Command,
		}
	})
	defer cleanup()

	client, err := NewClient(socketPath)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	const numGoroutines = 10
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(val Command) {
			defer wg.Done()
			res, err := client.SendCommand(val)
			if err != nil {
				t.Errorf("SendCommand returned an error: %v", err)
				return
			}
			if res != string(val) {
				t.Errorf("expected response %v, got %v", val, res)
			}
		}(Command(fmt.Sprintf("cmd-%d", i)))
	}

	wg.Wait()
}
