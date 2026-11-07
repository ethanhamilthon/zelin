package provider

import "testing"

func number[T any](value T) *T { return &value }

func TestUsageEstimate(t *testing.T) {
	pricing := Pricing{Prompt: number(0.01), Completion: number(0.1), CacheRead: number(0.001)}
	usage := (&rawUsage{PromptTokens: number[int64](100), CompletionTokens: number[int64](10)}).usage(pricing)
	if !usage.Estimated || *usage.Cost != 2 {
		t.Fatalf("%+v", usage)
	}
	raw := rawUsage{PromptTokens: number[int64](100), CompletionTokens: number[int64](10)}
	raw.Details.Cached = number[int64](50)
	if cost := *raw.usage(pricing).Cost; cost < 1.549 || cost > 1.551 {
		t.Fatal(cost)
	}
	if usage := (*rawUsage)(nil).usage(pricing); usage.Cost != nil || usage.Input != nil {
		t.Fatalf("%+v", usage)
	}
}

func TestTotals(t *testing.T) {
	totals := Totals{}.Add(Usage{Input: number[int64](100), Output: number[int64](10), Cost: number(0.5)})
	totals = totals.Add(Usage{Input: number[int64](200), Output: number[int64](20), CacheRead: number[int64](80), Cost: number(0.25), Estimated: true})
	totals = totals.Add(Usage{})
	if totals.Requests != 3 || totals.Unpriced != 1 || *totals.Input != 300 || *totals.Output != 30 || *totals.CacheRead != 80 || totals.CacheWrite != nil {
		t.Fatalf("%+v", totals)
	}
	if *totals.Cost != 0.75 || !totals.Estimated || totals.Context != nil {
		t.Fatalf("%+v", totals)
	}
}
