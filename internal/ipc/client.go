package ipc

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type Request struct {
	Id      uint64 `json:"id"`
	Command any    `json:"command"`
}

type Response struct {
	Id     uint64 `json:"id"`
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Client struct {
	nextId    atomic.Uint64
	responses map[uint64]chan Response
	lock      sync.Mutex
	reqCh     chan Request
	conn      net.Conn
}

func NewClient(socketPath string) (*Client, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}

	client := &Client{
		responses: make(map[uint64]chan Response),
		reqCh:     make(chan Request),
		conn:      conn,
	}

	go client.writer()
	go client.reader()

	return client, nil
}

func (c *Client) writer() {
	encoder := json.NewEncoder(c.conn)
	for req := range c.reqCh {
		encoder.Encode(req)
	}
}

func (c *Client) reader() {
	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		var response Response

		if err := json.Unmarshal(scanner.Bytes(), &response); err != nil {
			log.Printf("read response: %v", err)
			continue
		}

		c.lock.Lock()
		ch := c.responses[response.Id]
		delete(c.responses, response.Id)
		c.lock.Unlock()

		ch <- response
	}
}

func (c *Client) SendCommand(command any) (any, error) {
	id := c.nextId.Add(1)
	response := make(chan Response)

	c.lock.Lock()
	c.responses[id] = response
	c.lock.Unlock()

	c.reqCh <- Request{
		Id:      id,
		Command: command,
	}

	result := <-response

	if result.Error != "" {
		return nil, errors.New(result.Error)
	}

	return result.Result, nil
}

// Close closes the underlying connection to the server.
func (c *Client) Close() error {
	return c.conn.Close()
}
