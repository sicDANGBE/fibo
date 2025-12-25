package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"time"

	"github.com/streadway/amqp"
)

func (e *Engine) listenTasks() {
	e.Mu.Lock()
	ch := e.Channel
	e.Mu.Unlock()

	// Déclaration de l'exchange de diffusion
	ch.ExchangeDeclare("fibo_admin_exchange", "fanout", true, false, false, false, nil)

	// Création d'une queue temporaire pour ce worker spécifique
	q, _ := ch.QueueDeclare("", false, true, true, false, nil)
	ch.QueueBind(q.Name, "", "fibo_admin_exchange", false, nil)

	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	for d := range msgs {
		var task AdminTask
		if err := json.Unmarshal(d.Body, &task); err != nil {
			continue
		}

		// Synchronisation temporelle (Phase 3)
		wait := time.Until(time.Unix(task.StartAt, 0))
		if wait > 0 {
			log.Printf("[SYNC] Attente de %v pour démarrage synchrone...", wait)
			time.Sleep(wait)
		}

		log.Printf("[TASK] Démarrage du calcul: %s (ID: %s)", task.Handler, task.TaskID)
		e.runHandler(task)

		// Nettoyage mémoire après calcul intensif pour ton cluster IA
		runtime.GC()
		log.Println("[CLEAN] Ressources libérées après tâche.")
	}
}

func (e *Engine) runHandler(task AdminTask) {
	switch task.Handler {
	case "fibonacci":
		e.handleFibo(task)
	default:
		log.Printf("[WARN] Handler inconnu: %s", task.Handler)
	}
}

func (e *Engine) handleFibo(task AdminTask) {
	limit := 400000
	if val, ok := task.Params["limit"].(float64); ok {
		limit = int(val)
	}

	a, b := big.NewInt(0), big.NewInt(1)
	resQueue := "results_" + e.ID

	// On prépare les stats système une seule fois ou périodiquement
	var m runtime.MemStats

	for i := 0; i <= limit; i++ {
		a.Add(a, b)
		a, b = b, a

		// OPTIMISATION : On ne lit la RAM que toutes les 1000 itérations
		// pour ne pas ralentir le calcul pur, mais on ENVOIE chaque message.
		if i%1000 == 0 {
			runtime.ReadMemStats(&m)
		}

		res := WorkerResult{
			WorkerID:  e.ID,
			TaskID:    task.TaskID,
			Handler:   "fibonacci",
			Index:     i,
			Timestamp: time.Now().UnixMilli(),
			Metadata: map[string]interface{}{
				"cpu":     runtime.NumGoroutine(),
				"ram":     m.Alloc / 1024 / 1024,
				"net_io":  fmt.Sprintf("%.2f KB", float64(len(a.Bits())*8)/1024.0),
				"disk_io": "0.1 MB/s",
			},
		}

		body, _ := json.Marshal(res)

		// PUBLICATION À CHAQUE ITÉRATION
		e.Channel.Publish(
			"",
			resQueue,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
	}
	log.Printf("[FINISH] %d messages envoyés.", limit)
}
