package provider

type Usage struct {
	Input      *int64   `json:"input"`
	Output     *int64   `json:"output"`
	CacheRead  *int64   `json:"cacheRead"`
	CacheWrite *int64   `json:"cacheWrite"`
	Cost       *float64 `json:"cost"`
	Estimated  bool     `json:"estimated"`
}

type rawUsage struct {
	PromptTokens     *int64   `json:"prompt_tokens"`
	CompletionTokens *int64   `json:"completion_tokens"`
	Cost             *float64 `json:"cost"`
	Details          struct {
		Cached     *int64 `json:"cached_tokens"`
		CacheWrite *int64 `json:"cache_write_tokens"`
	} `json:"prompt_tokens_details"`
}

func (r *rawUsage) usage(pricing Pricing) Usage {
	if r == nil {
		return Usage{}
	}
	usage := Usage{Input: r.PromptTokens, Output: r.CompletionTokens, CacheRead: r.Details.Cached, CacheWrite: r.Details.CacheWrite, Cost: r.Cost}
	if usage.Cost == nil {
		usage.Cost = estimate(usage, pricing)
		usage.Estimated = usage.Cost != nil
	}
	return usage
}

func estimate(usage Usage, pricing Pricing) *float64 {
	if usage.Input == nil || usage.Output == nil || pricing.Prompt == nil || pricing.Completion == nil {
		return nil
	}
	read, write := value(usage.CacheRead), value(usage.CacheWrite)
	cost := float64(*usage.Input-read-write)**pricing.Prompt + float64(*usage.Output)**pricing.Completion
	cost += float64(read)*rate(pricing.CacheRead, *pricing.Prompt) + float64(write)*rate(pricing.CacheWrite, *pricing.Prompt)
	return &cost
}

func value(number *int64) int64 {
	if number == nil {
		return 0
	}
	return *number
}

func rate(price *float64, fallback float64) float64 {
	if price == nil {
		return fallback
	}
	return *price
}

type Totals struct {
	Requests   int      `json:"requests"`
	Unpriced   int      `json:"unpriced"`
	Context    *int64   `json:"context"`
	Input      *int64   `json:"input"`
	Output     *int64   `json:"output"`
	CacheRead  *int64   `json:"cacheRead"`
	CacheWrite *int64   `json:"cacheWrite"`
	Cost       *float64 `json:"cost"`
	Estimated  bool     `json:"estimated"`
}

func (t Totals) Add(usage Usage) Totals {
	t.Requests++
	t.Input, t.Output = sum(t.Input, usage.Input), sum(t.Output, usage.Output)
	t.CacheRead, t.CacheWrite = sum(t.CacheRead, usage.CacheRead), sum(t.CacheWrite, usage.CacheWrite)
	t.Context = nil
	if usage.Input != nil {
		context := *usage.Input + value(usage.Output)
		t.Context = &context
	}
	if usage.Cost == nil {
		t.Unpriced++
		return t
	}
	cost := *usage.Cost
	if t.Cost != nil {
		cost += *t.Cost
	}
	t.Cost, t.Estimated = &cost, t.Estimated || usage.Estimated
	return t
}

func sum(total, next *int64) *int64 {
	if next == nil {
		return total
	}
	result := value(total) + *next
	return &result
}
