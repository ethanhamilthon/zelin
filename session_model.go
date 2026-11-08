package main

import (
	"context"
	"slices"

	"zash/internal/provider"
)

type sessionClient struct {
	service *AgentService
	id      string
}

func (c sessionClient) Complete(ctx context.Context, request provider.Request) (provider.Message, provider.Usage, error) {
	if request.Model == "" {
		request.Model, request.Effort = c.service.selection(c.id)
	}
	return c.service.client.Complete(ctx, request)
}

func (s *AgentService) selection(id string) (string, string) {
	model, effort := s.client.Selection()
	if id == "" || s.store == nil {
		return model, effort
	}
	if saved, savedEffort, err := s.store.Model(id); err == nil && saved != "" {
		return saved, savedEffort
	}
	return model, effort
}

func (s *AgentService) pin(id string) error {
	if saved, _, err := s.store.Model(id); err != nil || saved != "" || s.client == nil {
		return err
	}
	model, effort := s.client.Selection()
	return s.store.SetModel(id, model, effort)
}

func (s *AgentService) setDefault(model, effort string) error {
	old, oldEffort := s.client.Selection()
	if err := s.store.PinModels(old, oldEffort); err != nil {
		return err
	}
	if err := s.store.SetSetting("model", model); err != nil {
		return err
	}
	if err := s.store.SetSetting("effort", effort); err != nil {
		return err
	}
	s.client.Select(model, effort)
	return nil
}

func (s *AgentService) Configure(ctx context.Context, id, model, effort string) (ModelInfo, error) {
	if _, err := s.client.Lookup(ctx, model); err != nil {
		return ModelInfo{}, err
	}
	if !slices.Contains(provider.Efforts, effort) {
		effort = provider.DefaultEffort
	}
	session, err := s.store.Get(id)
	if err != nil {
		return ModelInfo{}, err
	}
	pinned, _, err := s.store.Model(id)
	if err != nil {
		return ModelInfo{}, err
	}
	if len(session.Messages) == 0 && pinned == "" {
		err = s.setDefault(model, effort)
	} else {
		err = s.store.SetModel(id, model, effort)
	}
	if err != nil {
		return ModelInfo{}, err
	}
	return s.Info(ctx, id), nil
}
