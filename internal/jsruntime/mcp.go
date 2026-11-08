package jsruntime

import (
	"context"
	"encoding/json"

	"zelin/internal/mcp"
)

var MCP = &mcp.Manager{Servers: func() []mcp.Server { return nil }}

func (r *Runtime) bindMCP(ctx context.Context) {
	r.vm.Set("mcp", map[string]any{
		"servers": func() any {
			names := []string{}
			for _, server := range MCP.Servers() {
				names = append(names, server.Name)
			}
			return names
		},
		"tools": func(server string) any {
			tools, err := MCP.Tools(ctx, server, r.cwd)
			if err != nil {
				return failure(err)
			}
			return plain(tools)
		},
		"call": func(server, tool string, args map[string]any) any {
			text, err := MCP.Call(ctx, server, r.cwd, tool, args)
			if err != nil {
				return failure(err)
			}
			return text
		},
	})
}

func plain(value any) any {
	var result any
	data, _ := json.Marshal(value)
	json.Unmarshal(data, &result)
	return result
}
