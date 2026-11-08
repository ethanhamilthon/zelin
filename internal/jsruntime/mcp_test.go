package jsruntime

import (
	"context"
	"testing"
)

func TestMCPBindingWithoutServers(t *testing.T) {
	r := NewIn(t.TempDir())
	if out := r.Run(context.Background(), "mcp.servers()"); out != "[]" {
		t.Fatal(out)
	}
	if out := r.Run(context.Background(), "mcp.tools('none').error"); out != `"Unknown MCP server \"none\""` {
		t.Fatal(out)
	}
}
