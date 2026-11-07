package agent

type Execution struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Output string `json:"output"`
	State  string `json:"state"`
}

type Event struct {
	Kind      string     `json:"kind"`
	Activity  string     `json:"activity,omitempty"`
	Execution *Execution `json:"execution,omitempty"`
}
