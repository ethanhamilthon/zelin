package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"zelin/internal/home"
	"zelin/internal/provider"
	"zelin/internal/sessions"
)

const DefaultBase = "https://openrouter.ai/api/v1"

type SetupService struct {
	client *provider.Client
	store  *sessions.Store
	agent  *AgentService
}

type SetupCheck struct {
	Model  string `json:"model"`
	Reply  string `json:"reply"`
	Error  string `json:"error"`
	Millis int64  `json:"millis"`
}

func (s *SetupService) Required() (bool, error) {
	done, err := s.store.Setting("setup")
	if err != nil || done == "done" {
		return false, err
	}
	model, err := s.store.Setting("model")
	return model == "", err
}

func (s *SetupService) Connect(ctx context.Context, base, key string) ([]provider.Model, error) {
	base, key = strings.TrimSpace(base), strings.TrimSpace(key)
	if base == "" {
		base = DefaultBase
	}
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		return nil, fmt.Errorf("Provider URL must start with http:// or https://")
	}
	if key == "" {
		return nil, fmt.Errorf("API key is required")
	}
	s.client.Connect(base, key)
	models, err := s.client.Models(ctx)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, fmt.Errorf("Provider returned no models")
	}
	return models, home.SaveEnv(map[string]string{"OPENAI_BASE_URL": strings.TrimRight(base, "/"), "OPENAI_API_KEY": key})
}

func (s *SetupService) Check(ctx context.Context, model string) SetupCheck {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	started := time.Now()
	reply, _, err := s.client.Complete(ctx, provider.Request{Model: model, Messages: []provider.Message{{Role: "user", Content: "Reply with the single word OK."}}})
	check := SetupCheck{Model: model, Millis: time.Since(started).Milliseconds()}
	if err != nil {
		check.Error = err.Error()
	} else if check.Reply = strings.TrimSpace(reply.Content); check.Reply == "" {
		check.Error = "Model returned an empty reply"
	}
	return check
}

func (s *SetupService) Finish(ctx context.Context, model, titleModel string) (ModelInfo, error) {
	if _, err := s.client.Lookup(ctx, model); err != nil {
		return ModelInfo{}, err
	}
	if err := s.agent.setDefault(model, provider.DefaultEffort); err != nil {
		return ModelInfo{}, err
	}
	if _, err := s.agent.SetTitleModel(ctx, titleModel); err != nil {
		return ModelInfo{}, err
	}
	if err := s.store.SetSetting("setup", "done"); err != nil {
		return ModelInfo{}, err
	}
	return s.agent.Info(ctx, ""), nil
}
