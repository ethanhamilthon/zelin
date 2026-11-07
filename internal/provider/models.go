package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

var gatewayEfforts = []string{"max", "xhigh", "high", "medium", "low", "minimal", "none"}

type Pricing struct {
	Prompt     *float64 `json:"prompt"`
	Completion *float64 `json:"completion"`
	CacheRead  *float64 `json:"cacheRead"`
	CacheWrite *float64 `json:"cacheWrite"`
}

type Model struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	ContextWindow int      `json:"contextWindow"`
	Pricing       Pricing  `json:"pricing"`
	Efforts       []string `json:"efforts"`
	DefaultEffort string   `json:"defaultEffort"`
}

type rawModel struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	ContextWindow int               `json:"context_window"`
	ContextLength int               `json:"context_length"`
	MaxModelLen   int               `json:"max_model_len"`
	Pricing       map[string]string `json:"pricing"`
	TopProvider   struct {
		ContextLength int `json:"context_length"`
	} `json:"top_provider"`
	Reasoning *struct {
		SupportedEfforts json.RawMessage `json:"supported_efforts"`
		DefaultEffort    string          `json:"default_effort"`
		Mandatory        bool            `json:"mandatory"`
	} `json:"reasoning"`
}

func (r rawModel) model() Model {
	model := Model{ID: r.ID, Name: r.Name, Efforts: []string{}}
	if model.Name == "" {
		model.Name = r.ID
	}
	for _, size := range []int{r.ContextWindow, r.ContextLength, r.MaxModelLen, r.TopProvider.ContextLength} {
		if size > 0 {
			model.ContextWindow = size
			break
		}
	}
	model.Pricing = Pricing{price(r.Pricing["prompt"]), price(r.Pricing["completion"]), price(r.Pricing["input_cache_read"]), price(r.Pricing["input_cache_write"])}
	if r.Reasoning == nil {
		return model
	}
	var efforts []string
	switch raw := string(r.Reasoning.SupportedEfforts); raw {
	case "":
	case "null":
		efforts = gatewayEfforts
	default:
		json.Unmarshal(r.Reasoning.SupportedEfforts, &efforts)
	}
	for _, effort := range efforts {
		if !(r.Reasoning.Mandatory && effort == "none") {
			model.Efforts = append(model.Efforts, effort)
		}
	}
	if slices.Contains(model.Efforts, r.Reasoning.DefaultEffort) {
		model.DefaultEffort = r.Reasoning.DefaultEffort
	}
	return model
}

func price(value string) *float64 {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil || number < 0 {
		return nil
	}
	return &number
}

func (c *Client) Models(ctx context.Context) ([]Model, error) {
	c.catalog.Lock()
	defer c.catalog.Unlock()
	if models := c.cached(); models != nil {
		return models, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	url := strings.TrimSuffix(c.URL, "/chat/completions") + "/models"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Key)
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not load models: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Model list returned HTTP %d", res.StatusCode)
	}
	var data struct {
		Data []rawModel `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 16<<20)).Decode(&data); err != nil {
		return nil, fmt.Errorf("Invalid model list: %w", err)
	}
	models := make([]Model, 0, len(data.Data))
	for _, raw := range data.Data {
		models = append(models, raw.model())
	}
	c.mu.Lock()
	c.models = models
	c.mu.Unlock()
	return models, nil
}

func (c *Client) cached() []Model {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.models
}

func (c *Client) Lookup(ctx context.Context, id string) (Model, error) {
	models, err := c.Models(ctx)
	if err != nil {
		return Model{}, err
	}
	for _, model := range models {
		if model.ID == id {
			return model, nil
		}
	}
	return Model{}, fmt.Errorf("Unknown model %s", id)
}
