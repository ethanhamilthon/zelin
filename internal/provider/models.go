package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"zash/internal/metadata"
)

type Pricing = metadata.Pricing

type Model struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	ContextWindow int     `json:"contextWindow"`
	Pricing       Pricing `json:"pricing"`
}

func (c *Client) metadata(ctx context.Context, id string) metadata.Entry {
	if c.Meta == nil {
		return metadata.Entry{}
	}
	entry, _ := c.Meta.Lookup(ctx, id)
	return entry
}

func (c *Client) Models(ctx context.Context) ([]Model, error) {
	c.catalog.Lock()
	defer c.catalog.Unlock()
	if models := c.cached(); models != nil {
		return models, nil
	}
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
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 16<<20)).Decode(&data); err != nil {
		return nil, fmt.Errorf("Invalid model list: %w", err)
	}
	models := make([]Model, 0, len(data.Data))
	for _, raw := range data.Data {
		entry := c.metadata(ctx, raw.ID)
		model := Model{ID: raw.ID, Name: raw.Name, ContextWindow: entry.ContextWindow, Pricing: entry.Pricing}
		if model.Name == "" {
			model.Name = raw.ID
		}
		models = append(models, model)
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
