package metadata

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"
)

type Pricing struct {
	Prompt     *float64 `json:"prompt"`
	Completion *float64 `json:"completion"`
	CacheRead  *float64 `json:"cacheRead"`
	CacheWrite *float64 `json:"cacheWrite"`
}

type Entry struct {
	ContextWindow int     `json:"contextWindow"`
	Pricing       Pricing `json:"pricing"`
	Reasoning     *bool   `json:"reasoning,omitempty"`
}

type Catalog struct {
	Path    string
	Sources [2]string
	HTTP    *http.Client
	TTL     time.Duration
	mu      sync.Mutex
	entries map[string]Entry
	tried   time.Time
}

func New(path string) *Catalog {
	return &Catalog{Path: path, Sources: [2]string{OpenRouterURL, LiteLLMURL}, HTTP: &http.Client{Timeout: 30 * time.Second}, TTL: 24 * time.Hour}
}

func (c *Catalog) Lookup(ctx context.Context, id string) (Entry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) == 0 && time.Since(c.tried) > 5*time.Minute {
		c.entries, c.tried = c.load(ctx), time.Now()
	}
	entry, ok := c.entries[Key(id)]
	return entry, ok
}

func (c *Catalog) load(ctx context.Context) map[string]Entry {
	cached := map[string]Entry{}
	info, statErr := os.Stat(c.Path)
	if data, err := os.ReadFile(c.Path); err == nil {
		json.Unmarshal(data, &cached)
	}
	if statErr == nil && time.Since(info.ModTime()) < c.TTL && len(cached) > 0 {
		return cached
	}
	fresh := map[string]Entry{}
	err := errors.Join(openRouter(ctx, c.HTTP, c.Sources[0], fresh), liteLLM(ctx, c.HTTP, c.Sources[1], fresh))
	if len(fresh) == 0 || (err != nil && len(cached) > len(fresh)) {
		return cached
	}
	if data, err := json.Marshal(fresh); err == nil {
		os.WriteFile(c.Path, data, 0600)
	}
	return fresh
}
