package main

import (
	"encoding/json"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/pictl/pictl/internal/host"
	"github.com/pictl/pictl/internal/ipc"
)

func TestHandleConnStopsWhenClientCloses(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer clientConn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		handleConn(serverConn, &host.Dummy{})
	}()

	if err := json.NewEncoder(clientConn).Encode(ipc.Request{Id: 1, Command: "cpu_usage"}); err != nil {
		t.Fatalf("encode request: %v", err)
	}

	if err := clientConn.Close(); err != nil {
		t.Fatalf("close client: %v", err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("handleConn did not return after client disconnect")
	}
}

func TestIsSystemdSocketActivated(t *testing.T) {
	// Backup env vars
	origFds := os.Getenv("LISTEN_FDS")
	origPid := os.Getenv("LISTEN_PID")
	defer func() {
		os.Setenv("LISTEN_FDS", origFds)
		os.Setenv("LISTEN_PID", origPid)
	}()

	t.Run("absent LISTEN_FDS", func(t *testing.T) {
		os.Unsetenv("LISTEN_FDS")
		os.Unsetenv("LISTEN_PID")
		if isSystemdSocketActivated() {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("invalid LISTEN_FDS", func(t *testing.T) {
		os.Setenv("LISTEN_FDS", "0")
		os.Unsetenv("LISTEN_PID")
		if isSystemdSocketActivated() {
			t.Errorf("expected false, got true")
		}
	})

	t.Run("valid LISTEN_FDS without LISTEN_PID", func(t *testing.T) {
		os.Setenv("LISTEN_FDS", "1")
		os.Unsetenv("LISTEN_PID")
		if !isSystemdSocketActivated() {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("valid LISTEN_FDS with matching LISTEN_PID", func(t *testing.T) {
		os.Setenv("LISTEN_FDS", "1")
		os.Setenv("LISTEN_PID", strconv.Itoa(os.Getpid()))
		if !isSystemdSocketActivated() {
			t.Errorf("expected true, got false")
		}
	})

	t.Run("valid LISTEN_FDS with mismatched LISTEN_PID", func(t *testing.T) {
		os.Setenv("LISTEN_FDS", "1")
		os.Setenv("LISTEN_PID", "99999999")
		if isSystemdSocketActivated() {
			t.Errorf("expected false, got true")
		}
	})
}
