package agent

import "zash/internal/provider"

type Execution struct {
	ID     string `json:"id"`
	Kind   string `json:"kind,omitempty"`
	Code   string `json:"code"`
	Output string `json:"output"`
	State  string `json:"state"`
}

type Event struct {
	Session   string            `json:"session,omitempty"`
	Kind      string            `json:"kind"`
	Activity  string            `json:"activity,omitempty"`
	Execution *Execution        `json:"execution,omitempty"`
	Usage     *provider.Totals  `json:"usage,omitempty"`
	Message   *provider.Message `json:"message,omitempty"`
}
