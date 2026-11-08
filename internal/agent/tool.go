package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"zelin/internal/prompts"
	"zelin/internal/provider"
)

var Tool = provider.Tool{
	Name:        "js_repl",
	Description: prompts.JSTool,
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

func (a *Agent) takeImages() *provider.Message {
	images := a.runtime.TakeImages()
	if len(images) == 0 {
		return nil
	}
	lines := make([]string, len(images))
	urls := make([]string, len(images))
	for i, img := range images {
		lines[i] = fmt.Sprintf("- %s [Image #%d]", img.Path, i+1)
		urls[i] = img.DataURL
	}
	return &provider.Message{
		Role:    "user",
		Content: "Attached image(s) from view_image:\n" + strings.Join(lines, "\n"),
		Images:  urls,
	}
}
