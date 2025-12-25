package orchestrator

type WorkerRegistration struct {
	ID       string `json:"id"`
	Language string `json:"language"`
	// Utilise "-" pour que ce ne soit pas envoyé à l'IHM, mais reste accessible en Go
	LastSeen int64 `json:"_"`
}

type AdminTask struct {
	TaskID  string                 `json:"task_id"`
	Handler string                 `json:"handler"`
	StartAt int64                  `json:"start_at"`
	Params  map[string]interface{} `json:"params"`
}

type WorkerResult struct {
	WorkerID  string      `json:"worker_id"`
	TaskID    string      `json:"task_id"`
	Handler   string      `json:"handler"`
	Index     int         `json:"index"`
	Metadata  interface{} `json:"metadata"`
	Timestamp int64       `json:"timestamp"`
}
