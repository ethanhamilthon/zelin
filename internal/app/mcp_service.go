package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"zelin/internal/jsruntime"
	"zelin/internal/mcp"
	"zelin/internal/sessions"
)

type MCPService struct {
	store *sessions.Store
}

func (s *MCPService) List() []mcp.Server {
	var servers []mcp.Server
	if raw, err := s.store.Setting("mcp"); err == nil && raw != "" {
		if err := json.Unmarshal([]byte(raw), &servers); err != nil {
			log.Printf("mcp settings: %v", err)
		}
	}
	if servers == nil {
		servers = []mcp.Server{}
	}
	return servers
}

func (s *MCPService) save(servers []mcp.Server) error {
	data, err := json.Marshal(servers)
	if err != nil {
		return err
	}
	if err := s.store.SetSetting("mcp", string(data)); err != nil {
		return err
	}
	jsruntime.MCP.Reset()
	return nil
}

func (s *MCPService) Add(server mcp.Server) ([]mcp.Server, error) {
	server.Name, server.Command = strings.TrimSpace(server.Name), strings.TrimSpace(server.Command)
	if server.Name == "" || server.Command == "" {
		return nil, fmt.Errorf("Name and command are required")
	}
	servers := s.List()
	for _, existing := range servers {
		if existing.Name == server.Name {
			return nil, fmt.Errorf("A server named %s already exists", server.Name)
		}
	}
	servers = append(servers, server)
	return servers, s.save(servers)
}

func (s *MCPService) Remove(name string) ([]mcp.Server, error) {
	kept := []mcp.Server{}
	for _, server := range s.List() {
		if server.Name != name {
			kept = append(kept, server)
		}
	}
	return kept, s.save(kept)
}

func (s *MCPService) Test(name string) ([]mcp.Tool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	dir, _ := os.UserHomeDir()
	return jsruntime.MCP.Tools(ctx, name, dir)
}

func (s *MCPService) attach() {
	jsruntime.MCP.Servers = s.List
}
