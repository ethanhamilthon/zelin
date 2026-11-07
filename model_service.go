package main

import (
	"context"
	"os"
	"slices"
	"strconv"

	"zash/internal/provider"
	"zash/internal/sessions"
)

type ModelInfo struct {
	Model         string           `json:"model"`
	Name          string           `json:"name"`
	Effort        string           `json:"effort"`
	ContextWindow int              `json:"contextWindow"`
	Efforts       []string         `json:"efforts"`
	DefaultEffort string           `json:"defaultEffort"`
	Pricing       provider.Pricing `json:"pricing"`
}

func (s *AgentService) Info(ctx context.Context) ModelInfo {
	id, effort := s.client.Selection()
	info := ModelInfo{Model: id, Name: id, Effort: effort, Efforts: []string{}}
	if model, err := s.client.Lookup(ctx, id); err == nil {
		info.Name, info.ContextWindow, info.Efforts = model.Name, model.ContextWindow, model.Efforts
		info.DefaultEffort, info.Pricing = model.DefaultEffort, model.Pricing
	}
	if size, err := strconv.Atoi(os.Getenv("OPENAI_CONTEXT_WINDOW")); err == nil && size > 0 {
		info.ContextWindow = size
	}
	return info
}

func (s *AgentService) Models(ctx context.Context) ([]provider.Model, error) {
	return s.client.Models(ctx)
}

func (s *AgentService) Configure(ctx context.Context, id, effort string) (ModelInfo, error) {
	model, err := s.client.Lookup(ctx, id)
	if err != nil {
		return ModelInfo{}, err
	}
	if !slices.Contains(model.Efforts, effort) {
		effort = ""
	}
	if err := s.store.SetSetting("model", id); err != nil {
		return ModelInfo{}, err
	}
	if err := s.store.SetSetting("effort", effort); err != nil {
		return ModelInfo{}, err
	}
	s.client.Select(id, effort)
	return s.Info(ctx), nil
}

func restoreSelection(store *sessions.Store, client *provider.Client) error {
	model, err := store.Setting("model")
	if err != nil || model == "" {
		return err
	}
	effort, err := store.Setting("effort")
	client.Select(model, effort)
	return err
}
