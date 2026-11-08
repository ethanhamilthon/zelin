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
	Pricing       provider.Pricing `json:"pricing"`
	TitleModel    string           `json:"titleModel"`
}

func (s *AgentService) Info(ctx context.Context, session string) ModelInfo {
	id, effort := s.selection(session)
	info := ModelInfo{Model: id, Name: id, Effort: effort, Efforts: provider.Efforts}
	if model, err := s.client.Lookup(ctx, id); err == nil {
		info.Name, info.ContextWindow, info.Pricing = model.Name, model.ContextWindow, model.Pricing
	}
	if size, err := strconv.Atoi(os.Getenv("OPENAI_CONTEXT_WINDOW")); err == nil && size > 0 {
		info.ContextWindow = size
	}
	info.TitleModel, _ = s.store.Setting("titleModel")
	return info
}

func (s *AgentService) Models(ctx context.Context) ([]provider.Model, error) {
	return s.client.Models(ctx)
}

func (s *AgentService) SetTitleModel(ctx context.Context, id string) (ModelInfo, error) {
	if _, err := s.client.Lookup(ctx, id); err != nil {
		return ModelInfo{}, err
	}
	if err := s.store.SetSetting("titleModel", id); err != nil {
		return ModelInfo{}, err
	}
	return s.Info(ctx, ""), nil
}

func restoreSelection(store *sessions.Store, client *provider.Client) error {
	model, err := store.Setting("model")
	if err != nil || model == "" {
		return err
	}
	effort, err := store.Setting("effort")
	if !slices.Contains(provider.Efforts, effort) {
		effort = provider.DefaultEffort
	}
	client.Select(model, effort)
	return err
}
