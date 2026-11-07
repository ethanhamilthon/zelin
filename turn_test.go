package main

import (
	"testing"
	"time"

	"zash/internal/provider"
)

func TestTurnReportsOnlyThisRun(t *testing.T) {
	first, total, context := 0.25, 0.75, int64(1200)
	before := provider.Totals{Requests: 2, Cost: &first}
	after := provider.Totals{Requests: 4, Cost: &total, Context: &context}
	turn := (&AgentService{}).turn(time.Now().Add(-3*time.Second), before, after)
	if turn.Cost == nil || *turn.Cost != 0.5 || *turn.Context != 1200 || turn.Duration < 3000 || turn.Finished == "" {
		t.Fatalf("unexpected turn %+v", turn)
	}
	unpriced := provider.Totals{Requests: 3, Unpriced: 1, Cost: &first}
	if turn := (&AgentService{}).turn(time.Now(), before, unpriced); turn.Cost != nil {
		t.Fatalf("cost should be unknown, got %v", *turn.Cost)
	}
}
