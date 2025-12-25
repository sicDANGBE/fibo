package worker

import (
	"crypto/sha256"
	"fmt"
	"os"
)

type WorkerRegistration struct {
	ID       string `json:"id"`
	Language string `json:"language"`
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

// GenerateID crée l'identifiant unique SHA256 du worker au démarrage
func GenerateID() string {
	hostname, _ := os.Hostname()
	data := fmt.Sprintf("%s-%d", hostname, os.Getpid())
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}
