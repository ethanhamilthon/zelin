package main

import (
	"context"
	"time"

	"zelin/internal/provider"
	"zelin/internal/sessions"
)

func (s *AgentService) turn(id string, started time.Time, before, after provider.Totals) *sessions.Turn {
	turn := &sessions.Turn{Finished: time.Now().UTC().Format(time.RFC3339), Duration: time.Since(started).Milliseconds(), Context: after.Context}
	priced := after.Requests - before.Requests - (after.Unpriced - before.Unpriced)
	if priced > 0 && after.Cost != nil {
		cost := *after.Cost
		if before.Cost != nil {
			cost -= *before.Cost
		}
		turn.Cost = &cost
	}
	if s.client != nil {
		turn.Window = s.Info(context.Background(), id).ContextWindow
	}
	return turn
}
