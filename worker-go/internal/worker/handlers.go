package worker

import (
	"encoding/json"
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

	// On s'assure que la queue de résultats existe (Phase 2 & 4)
	e.Mu.Lock()
	e.Channel.QueueDeclare(resQueue, false, true, false, false, nil)
	e.Mu.Unlock()

	for i := 0; i <= limit; i++ {
		a.Add(a, b)
		a, b = b, a

		if i%10000 == 0 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			res := WorkerResult{
				WorkerID:  e.ID,
				TaskID:    task.TaskID,
				Handler:   "fibonacci",
				Index:     i,
				Timestamp: time.Now().UnixMilli(),
				Metadata: map[string]interface{}{
					"value": a.String()[:10] + "...", // Troncature pour éviter les gros messages
					"cpu":   runtime.NumGoroutine(),  // Approximation du CPU utilisé
					"ram":   m.Alloc / 1024 / 1024,   // RAM en MB
				},
			}
			body, _ := json.Marshal(res)

			e.Mu.Lock()
			e.Channel.Publish("", resQueue, false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
			e.Mu.Unlock()
		}
	}
	log.Printf("[FINISH] Tâche %s terminée.", task.TaskID)
}
