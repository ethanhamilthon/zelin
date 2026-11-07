package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Metadata struct {
	Model         string `json:"model"`
	ContextWindow int    `json:"contextWindow"`
}

func (c *Client) Metadata(ctx context.Context) Metadata {
	info := Metadata{Model: c.Model}
	if size, err := strconv.Atoi(os.Getenv("OPENAI_CONTEXT_WINDOW")); err == nil && size > 0 {
		info.ContextWindow = size
		return info
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	url := strings.TrimSuffix(c.URL, "/chat/completions") + "/models"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return info
	}
	req.Header.Set("Authorization", "Bearer "+c.Key)
	res, err := c.HTTP.Do(req)
	if err != nil {
		return info
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return info
	}
	var models struct {
		Data []struct {
			ID            string `json:"id"`
			ContextWindow int    `json:"context_window"`
			ContextLength int    `json:"context_length"`
			MaxModelLen   int    `json:"max_model_len"`
			TopProvider   struct {
				ContextLength int `json:"context_length"`
			} `json:"top_provider"`
		} `json:"data"`
	}
	if json.NewDecoder(io.LimitReader(res.Body, 2<<20)).Decode(&models) != nil {
		return info
	}
	for _, model := range models.Data {
		if model.ID != c.Model {
			continue
		}
		for _, size := range []int{model.ContextWindow, model.ContextLength, model.MaxModelLen, model.TopProvider.ContextLength} {
			if size > 0 {
				info.ContextWindow = size
				return info
			}
		}
	}
	return info
}
