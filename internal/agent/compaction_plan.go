package agent

import (
	"encoding/json"
	"log"
	"strings"

	"zelin/internal/provider"
)

type CompactionDiag struct {
	Reason          string
	UtilityModel    string
	Threshold       int
	MainWindow      int
	UtilityWindow   int
	ContextEstimate int64
	Chunks          int
	DurationMs      int64
	Error           string
}

func logCompaction(d CompactionDiag) {
	log.Printf("compaction reason=%s model=%s threshold=%d main_window=%d utility_window=%d estimate=%d chunks=%d duration_ms=%d error=%q",
		d.Reason, d.UtilityModel, d.Threshold, d.MainWindow, d.UtilityWindow, d.ContextEstimate, d.Chunks, d.DurationMs, d.Error)
}

type compactEntry struct {
	Role       string              `json:"role"`
	Content    string              `json:"content,omitempty"`
	Images     int                 `json:"attached_images,omitempty"`
	ToolCalls  []provider.ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string              `json:"tool_call_id,omitempty"`
}

func splitWorkingContext(body []provider.Message, keepBudget int64) ([]provider.Message, []provider.Message) {
	if len(body) <= 2 {
		return body, nil
	}
	if keepBudget <= 0 {
		keepBudget = 8000
	}
	split := -1
	for i := len(body) - 1; i >= 2; i-- {
		if body[i].Role != "user" || strings.HasPrefix(body[i].Content, "Attached image(s) from view_image:") {
			continue
		}
		tailTokens := Estimate(body[i:])
		if split == -1 && tailTokens <= keepBudget*2 {
			split = i
		}
		if tailTokens <= keepBudget && len(body)-i <= 6 {
			split = i
		} else if split != -1 {
			break
		}
	}
	if split < 2 {
		return body, nil
	}
	return body[:split], body[split:]
}

func chunkOlder(older []provider.Message, utilityWindow int) []string {
	if utilityWindow <= 0 {
		utilityWindow = 128000
	}
	budget := max(int64(250), int64(utilityWindow)*60/100)
	maxChars := int(budget * 4)
	var chunks []string
	var current []compactEntry
	var currentTokens int64
	flush := func() {
		if len(current) == 0 {
			return
		}
		data, _ := json.Marshal(current)
		chunks = append(chunks, string(data))
		current, currentTokens = nil, 0
	}
	for _, m := range older {
		text := m.Content
		if len(text) > maxChars {
			text = text[:maxChars] + "\n[truncated for compaction]"
		}
		imgCount := 0
		for _, img := range m.Images {
			if img != "" {
				imgCount++
			}
		}
		item := compactEntry{Role: m.Role, Content: text, Images: imgCount, ToolCalls: m.ToolCalls, ToolCallID: m.ToolCallID}
		raw, _ := json.Marshal(item)
		itemTokens := int64((len(raw) + 3) / 4)
		if len(current) > 0 && currentTokens+itemTokens > budget {
			flush()
		}
		current = append(current, item)
		currentTokens += itemTokens
	}
	flush()
	return chunks
}
