package sessions

type Turn struct {
	Finished string   `json:"finished"`
	Duration int64    `json:"duration"`
	Cost     *float64 `json:"cost"`
	Context  *int64   `json:"context"`
	Window   int      `json:"window"`
}
