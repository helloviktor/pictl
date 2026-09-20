package ipc

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type Request struct {
	ID      uint64  `json:"id"`
	Command Command `json:"command"`
}

type Response struct {
	ID     uint64 `json:"id"`
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Client struct {
	nextID      atomic.Uint64
	responses   map[uint64]chan Response
	responsesMu sync.Mutex
	requestCh   chan Request
	conn        net.Conn

	shutdownOnce sync.Once
	shutdownErr  error
	shutdownCh   chan struct{}
	workers      sync.WaitGroup
}

func NewClient(socketPath string) (*Client, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("connect to IPC socket %q: %w", socketPath, err)
	}

	client := &Client{
		responses:  make(map[uint64]chan Response),
		requestCh:  make(chan Request),
		conn:       conn,
		shutdownCh: make(chan struct{}),
	}

	client.workers.Add(2)
	go client.writer()
	go client.reader()

	return client, nil
}

func (c *Client) shutdown() {
	c.shutdownOnce.Do(func() {
		close(c.shutdownCh)
		c.shutdownErr = c.conn.Close()
	})
}

func (c *Client) writer() {
	defer c.workers.Done()
	defer c.shutdown()

	encoder := json.NewEncoder(c.conn)
	for {
		select {
		case <-c.shutdownCh:
			return
		case req := <-c.requestCh:
			if err := encoder.Encode(req); err != nil {
				log.Printf("write request: %v", err)
				return
			}
		}
	}
}

func (c *Client) reader() {
	defer c.workers.Done()
	defer c.shutdown()

	scanner := bufio.NewScanner(c.conn)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		var response Response

		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			log.Printf("read response: %v", err)
			return
		}

		c.responsesMu.Lock()
		ch, ok := c.responses[response.ID]
		delete(c.responses, response.ID)
		c.responsesMu.Unlock()

		if ok {
			ch <- response
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("read response: %v\n", err)
	}
}

func (c *Client) SendCommand(ctx context.Context, command Command) (any, error) {
	if ctx == nil {
		return nil, errors.New("IPC client: nil context")
	}

	id := c.nextID.Add(1)
	response := make(chan Response, 1)

	c.responsesMu.Lock()
	c.responses[id] = response
	c.responsesMu.Unlock()

	cleanup := func() {
		c.responsesMu.Lock()
		delete(c.responses, id)
		c.responsesMu.Unlock()
	}

	select {
	case c.requestCh <- Request{ID: id, Command: command}:
	case <-c.shutdownCh:
		cleanup()
		return nil, errors.New("send IPC request: client is shut down")
	case <-ctx.Done():
		cleanup()
		return nil, fmt.Errorf("send IPC request: %w", ctx.Err())
	}

	select {
	case <-c.shutdownCh:
		cleanup()
		return nil, errors.New("receive IPC response: client is shut down")
	case <-ctx.Done():
		cleanup()
		return nil, fmt.Errorf("receive IPC response: %w", ctx.Err())
	case result := <-response:
		if result.Error != "" {
			return nil, fmt.Errorf("IPC command %q failed: %s", command, result.Error)
		}
		return result.Result, nil
	}
}

// Close closes the underlying connection to the server.
func (c *Client) Close() error {
	c.shutdown()
	c.workers.Wait()
	return c.shutdownErr
}
