package tracer

type Event struct {
	Timestamp int64                  `json:"ts,omitempty"`
	Name      string                 `json:"name,omitempty"`
	PID       int                    `json:"pid,omitempty"`
	TID       int                    `json:"tid,omitempty"`
	Phase     string                 `json:"ph"`
	Args      map[string]interface{} `json:"args,omitempty"`
}
