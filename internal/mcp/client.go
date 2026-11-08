package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
)

type Server struct {
	Name    string            `json:"name"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

type Client struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	mu      sync.Mutex
	write   sync.Mutex
	pending map[int64]chan response
	next    atomic.Int64
	done    chan struct{}
	err     error
}

type response struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func Start(server Server, dir string) (*Client, error) {
	cmd := exec.Command(server.Command, server.Args...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	for key, value := range server.Env {
		cmd.Env = append(cmd.Env, key+"="+value)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Could not start MCP server %s: %w", server.Name, err)
	}
	c := &Client{cmd: cmd, stdin: stdin, pending: map[int64]chan response{}, done: make(chan struct{})}
	go c.read(stdout)
	return c, nil
}

func (c *Client) read(stdout io.Reader) {
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 64<<10), 32<<20)
	for scanner.Scan() {
		var message struct {
			ID *int64 `json:"id"`
			response
		}
		if json.Unmarshal(scanner.Bytes(), &message) != nil || message.ID == nil {
			continue
		}
		c.mu.Lock()
		reply := c.pending[*message.ID]
		delete(c.pending, *message.ID)
		c.mu.Unlock()
		if reply != nil {
			reply <- message.response
		}
	}
	c.mu.Lock()
	c.err = fmt.Errorf("MCP server exited")
	c.mu.Unlock()
	close(c.done)
}

func (c *Client) Call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	id := c.next.Add(1)
	reply := make(chan response, 1)
	c.mu.Lock()
	if c.err != nil {
		c.mu.Unlock()
		return nil, c.err
	}
	c.pending[id] = reply
	c.mu.Unlock()
	if err := c.send(map[string]any{"jsonrpc": "2.0", "id": id, "method": method, "params": params}); err != nil {
		return nil, err
	}
	select {
	case r := <-reply:
		if r.Error != nil {
			return nil, fmt.Errorf("MCP error %d: %s", r.Error.Code, r.Error.Message)
		}
		return r.Result, nil
	case <-c.done:
		return nil, fmt.Errorf("MCP server exited")
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, ctx.Err()
	}
}

func (c *Client) Notify(method string) error {
	return c.send(map[string]any{"jsonrpc": "2.0", "method": method})
}

func (c *Client) send(message any) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	c.write.Lock()
	defer c.write.Unlock()
	_, err = c.stdin.Write(append(data, '\n'))
	return err
}

func (c *Client) Close() {
	c.stdin.Close()
	if c.cmd.Process != nil {
		c.cmd.Process.Kill()
	}
	c.cmd.Wait()
}
