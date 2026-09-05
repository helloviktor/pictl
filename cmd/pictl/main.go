package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
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
	"github.com/pictl/pictl/internal/services"
)

func main() {
	socketPath := flag.String("socket", "", "path to a Unix socket to listen on for commands")
	flag.Usage = func() {
		fmt.Println("Usage: pictl [-socket <path>] <info|update|restart|shutdown>")
		fmt.Println("\nCommands:")
		fmt.Println("  info      print current system information as JSON")
		fmt.Println("  update    update installed system packages")
		fmt.Println("  restart   restart the device")
		fmt.Println("  shutdown  shut down the device")
		fmt.Println("\nFlags:")
		fmt.Println("  -socket   listen on a Unix socket and execute commands received over it, one per line")
	}
	flag.Parse()

	service := services.NewSystemService(host.NewHost())

	if *socketPath != "" {
		if err := serveSocket(*socketPath, service); err != nil {
			log.Fatal(err)
		}
		return
	}

	if flag.NArg() == 0 {
		if isSystemdSocketActivated() {
			if err := serveSystemd(service); err != nil {
				log.Fatal(err)
			}
			return
		}
		flag.Usage()
		return
	}

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	if err := runCommand(flag.Arg(0), service, os.Stdout); err != nil {
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
func serveSystemd(service *services.SystemService) error {
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

	return serveListener(listener, service)
}

// runCommand executes a single command, writing any result to out
func runCommand(command string, service *services.SystemService, out io.Writer) error {
	result, err := executeCommand(command, service)
	if err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(result)
}

// executeCommand runs command against service and returns its result
func executeCommand(command string, service *services.SystemService) (any, error) {
	switch command {
	case "info":
		return service.SystemInfo()
	case "update":
		if err := service.UpdateSystem(); err != nil {
			return nil, err
		}
		return models.UpdateResponse{Success: true, Message: "System update started"}, nil
	case "restart":
		if err := service.RestartSystem(); err != nil {
			return nil, err
		}
		return models.UpdateResponse{Success: true, Message: "System restart initiated"}, nil
	case "shutdown":
		if err := service.ShutdownSystem(); err != nil {
			return nil, err
		}
		return models.UpdateResponse{Success: true, Message: "System shutdown initiated"}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}

// serveSocket listens on a Unix socket at socketPath and executes commands received over connections
func serveSocket(socketPath string, service *services.SystemService) error {
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

	return serveListener(listener, service)
}

// serveListener accepts connections on listener and executes commands
func serveListener(listener net.Listener, service *services.SystemService) error {
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

		go handleConn(conn, service)
	}
}

// handleConn reads newline-terminated JSON requests from conn, executes them, and writes back JSON responses
func handleConn(conn net.Conn, service *services.SystemService) {
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

		result, err := executeCommand(command, service)
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
