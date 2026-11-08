package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const (
	OpenRouterURL = "https://openrouter.ai/api/v1/models"
	LiteLLMURL    = "https://raw.githubusercontent.com/BerriAI/litellm/main/model_prices_and_context_window.json"
)

func fetch(ctx context.Context, client *http.Client, url string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned HTTP %d", url, res.StatusCode)
	}
	return json.NewDecoder(io.LimitReader(res.Body, 64<<20)).Decode(target)
}

func openRouter(ctx context.Context, client *http.Client, url string, entries map[string]Entry) error {
	var data struct {
		Data []struct {
			ID                  string            `json:"id"`
			ContextLength       int               `json:"context_length"`
			Pricing             map[string]string `json:"pricing"`
			SupportedParameters []string          `json:"supported_parameters"`
		} `json:"data"`
	}
	if err := fetch(ctx, client, url, &data); err != nil {
		return err
	}
	for _, model := range data.Data {
		key := Key(model.ID)
		if _, exists := entries[key]; exists || strings.Contains(model.ID, ":") {
			continue
		}
		p := model.Pricing
		var reasoning *bool
		if len(model.SupportedParameters) > 0 {
			supports := slices.Contains(model.SupportedParameters, "reasoning") || slices.Contains(model.SupportedParameters, "include_reasoning") || slices.Contains(model.SupportedParameters, "reasoning_effort")
			reasoning = &supports
		}
		entries[key] = Entry{model.ContextLength, Pricing{parse(p["prompt"]), parse(p["completion"]), parse(p["input_cache_read"]), parse(p["input_cache_write"])}, reasoning}
	}
	return nil
}

func liteLLM(ctx context.Context, client *http.Client, url string, entries map[string]Entry) error {
	var data map[string]json.RawMessage
	if err := fetch(ctx, client, url, &data); err != nil {
		return err
	}
	for _, direct := range []bool{true, false} {
		for _, id := range slices.Sorted(maps.Keys(data)) {
			raw := data[id]
			var model struct {
				Input      int      `json:"max_input_tokens"`
				Prompt     *float64 `json:"input_cost_per_token"`
				Completion *float64 `json:"output_cost_per_token"`
				CacheRead  *float64 `json:"cache_read_input_token_cost"`
				CacheWrite *float64 `json:"cache_creation_input_token_cost"`
				Reasoning  *bool    `json:"supports_reasoning"`
			}
			key := Key(id)
			if _, exists := entries[key]; exists || direct == strings.Contains(id, "/") || json.Unmarshal(raw, &model) != nil || model.Input == 0 {
				continue
			}
			entries[key] = Entry{model.Input, Pricing{model.Prompt, model.Completion, model.CacheRead, model.CacheWrite}, model.Reasoning}
		}
	}
	return nil
}

func parse(value string) *float64 {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil || number < 0 {
		return nil
	}
	return &number
}
