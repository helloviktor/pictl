package main

import (
	"os"
	"strconv"
	"testing"
)

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
