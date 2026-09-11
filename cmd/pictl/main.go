package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/pictl/pictl/internal/host"
	"github.com/pictl/pictl/internal/ipc"
	"github.com/pictl/pictl/internal/models"
)

func main() {
	socketPath := flag.String("socket", "", "path to a Unix socket to listen on for commands")
	flag.Usage = func() {
		fmt.Println("Usage: pictl -socket <path>")
		fmt.Println("\nListens on a Unix socket and executes commands received over it, one JSON request per line.")
		fmt.Println("With no -socket flag, expects to be launched via systemd socket activation.")
	}
	flag.Parse()

	if *socketPath == "" && !isSystemdSocketActivated() {
		flag.Usage()
		os.Exit(1)
	}

	h, err := host.NewHost()
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	if *socketPath != "" {
		if err := serveSocket(*socketPath, h); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := serveSystemd(h); err != nil {
		log.Fatal(err)
	}
}

// isSystemdSocketActivated returns true if the process was launched via systemd socket activation
func isSystemdSocketActivated() bool {
	fds, err := strconv.Atoi(os.Getenv("LISTEN_FDS"))
	if err != nil || fds < 1 {
		return false
	}
	if pid := os.Getenv("LISTEN_PID"); pid != "" && pid != strconv.Itoa(os.Getpid()) {
		return false
	}
	return true
}

// serveSystemd listens on the socket passed by systemd via file descriptor 3
func serveSystemd(h host.Host) error {
	file := os.NewFile(uintptr(3), "systemd-socket")
	if file == nil {
		return errors.New("failed to create file from fd 3")
	}
	defer file.Close()

	listener, err := net.FileListener(file)
	if err != nil {
		return fmt.Errorf("create listener from systemd socket: %w", err)
	}
	defer listener.Close()

	log.Println("Listening on systemd socket")

	return serveListener(listener, h)
}

// executeCommand runs command against h and returns its result
func executeCommand(command string, h host.Host) (any, error) {
	switch command {
	case "cpu_usage":
		return h.CPUUsage()
	case "memory_usage":
		return h.MemoryUsage()
	case "disk_usage":
		return h.DiskUsage()
	case "cpu_temperature":
		return h.CPUTemperature()
	case "available_updates":
		return h.AvailableUpdates()
	case "apply_updates":
		if err := h.ApplyUpdates(); err != nil {
			return nil, err
		}
		return models.UpdateResponse{Success: true, Message: "System update started"}, nil
	case "restart":
		if err := h.Restart(); err != nil {
			return nil, err
		}
		return models.UpdateResponse{Success: true, Message: "System restart initiated"}, nil
	case "shutdown":
		if err := h.Shutdown(); err != nil {
			return nil, err
		}
		return models.UpdateResponse{Success: true, Message: "System shutdown initiated"}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}

// serveSocket listens on a Unix socket at socketPath and executes commands received over connections
func serveSocket(socketPath string, h host.Host) error {
	// Remove any stale socket file left over from a previous run
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove stale socket: %w", err)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("listen on socket: %w", err)
	}
	defer listener.Close()
	defer os.Remove(socketPath)

	log.Printf("Listening on socket %s\n", socketPath)

	return serveListener(listener, h)
}

// serveListener accepts connections on listener and executes commands
func serveListener(listener net.Listener, h host.Host) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			// A single failed accept shouldn't bring down the server; log and keep serving.
			log.Printf("accept connection: %v\n", err)
			continue
		}

		go handleConn(conn, h)
	}
}

// handleConn reads newline-terminated JSON requests from conn, executes them, and writes back JSON responses
func handleConn(conn net.Conn, h host.Host) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	encoder := json.NewEncoder(conn)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var req ipc.Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			encoder.Encode(ipc.Response{Error: fmt.Sprintf("invalid request: %v", err)})
			continue
		}

		command, ok := req.Command.(string)
		if !ok {
			encoder.Encode(ipc.Response{Id: req.Id, Error: "command must be a string"})
			continue
		}

		result, err := executeCommand(command, h)
		resp := ipc.Response{Id: req.Id, Result: result}
		if err != nil {
			resp.Error = err.Error()
		}
		encoder.Encode(resp)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("read request: %v\n", err)
	}
}
