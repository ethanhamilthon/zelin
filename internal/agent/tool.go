package agent

import (
	"context"
	"encoding/json"

	"zash/internal/provider"
)

const toolDescription = `Execute synchronous JavaScript in a persistent Goja REPL.
Host APIs:
- bash(cmd, timeout): runs a shell command, returns {stdout, stderr, exitCode, timedOut}. Timeout is in seconds, default 30, maximum 120.
- fs.read(path, {offset, limit}): returns file text with line numbers (default 2000 lines from line 1).
- fs.write(path, text): creates or overwrites a file, creating parent directories.
- fs.edit(path, oldText, newText): replaces exactly one occurrence of oldText.
- fs.patch(path, diff): applies a unified diff with @@ hunks; nothing is written if any hunk fails.
- fs.grep(pattern, path, {glob}): regex search, returns ["file:line:text"], skips .git and node_modules.
fs functions return {error} on failure; write, edit and patch return {ok: true} on success. Paths may start with ~.
No async/await, promises, console, require, or other host APIs.
Use the last expression as the tool output.
Variables persist across turns; reuse variables or use var to avoid redeclaration errors.
Commands run with the user's permissions in the app's working directory.`

var Tool = provider.Tool{
	Name:        "js_repl",
	Description: toolDescription,
	Parameters: map[string]any{
		"type":                 "object",
		"properties":           map[string]any{"code": map[string]string{"type": "string", "description": "JavaScript source to execute"}},
		"required":             []string{"code"},
		"additionalProperties": false,
	},
}

func (a *Agent) call(ctx context.Context, call provider.ToolCall, emit func(Event)) string {
	if call.Function.Name != Tool.Name {
		return "Unknown tool; use js_repl only"
	}
	var args struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
		return "Invalid arguments; expected a JSON object with a code string"
	}
	execution := Execution{ID: call.ID, Code: args.Code, State: "running"}
	emit(Event{Kind: "activity", Activity: "Running JavaScript"})
	emit(Event{Kind: "execution", Execution: &execution})
	result := a.runtime.Run(ctx, args.Code)
	completed := Execution{ID: call.ID, Code: args.Code, Output: result, State: "complete"}
	emit(Event{Kind: "execution", Execution: &completed})
	return result
}
