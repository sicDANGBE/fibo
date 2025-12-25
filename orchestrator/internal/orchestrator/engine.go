package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Engine struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Workers map[string]WorkerRegistration
	Mu      sync.Mutex
	Hub     UIHub
}

type UIHub interface {
	BroadcastMessage(msg interface{})
}

func (o *Engine) safeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[CRITICAL] Panic capturé dans l'orchestrateur: %v", r)
				// Ici, on pourrait envoyer une alerte à ton Sentry ou Loki
			}
		}()
		fn()
	}()
}

func NewEngine(amqpURL string, hub UIHub) *Engine {
	e := &Engine{
		Workers: make(map[string]WorkerRegistration),
		Hub:     hub,
	}
	e.InitRabbitMQ(amqpURL)
	return e
}

// StartTask diffuse l'ordre de calcul avec synchronisation temporelle (Phase 3)
func (o *Engine) StartTask(handler string, params map[string]interface{}) {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if o.Channel == nil {
		log.Println("[WARN] Abandon : Orchestrateur non connecté à RabbitMQ")
		return
	}

	task := AdminTask{
		TaskID:  fmt.Sprintf("T-%d", time.Now().Unix()),
		Handler: handler,
		StartAt: time.Now().Add(5 * time.Second).Unix(), // Barrière T+5s pour tous les workers
		Params:  params,
	}

	body, _ := json.Marshal(task)
	err := o.Channel.Publish(
		"fibo_admin_exchange",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("[RMQ] Échec envoi Task: %v", err)
	}
}

// ConsumeWorkerResults traite les messages entrants de chaque langage

func (o *Engine) ConsumeWorkerResults(queueName string) {
	// 1. On crée un canal DÉDIÉ pour ce flux Stream spécifique
	ch, err := o.Conn.Channel()
	if err != nil {
		log.Printf("[ERROR] Impossible de créer un canal pour %s: %v", queueName, err)
		return
	}
	defer ch.Close() // Fermeture propre si le worker disparaît

	// 2. Configuration obligatoire pour Stream (Prefetch > 0)
	if err := ch.Qos(100, 0, false); err != nil {
		log.Printf("[RMQ] Erreur QoS pour %s: %v", queueName, err)
		return
	}

	// 3. Consommation avec offset 'next'
	msgs, err := ch.Consume(
		queueName, "",
		false, // auto-ack: false obligatoire
		false, false, false,
		amqp.Table{"x-stream-offset": "next"},
	)
	if err != nil {
		log.Printf("[RMQ] Erreur Consume %s: %v", queueName, err)
		return
	}

	for d := range msgs {
		var res WorkerResult
		if err := json.Unmarshal(d.Body, &res); err == nil {
			o.BroadcastToUI("RESULT", res)
			d.Ack(false) // Ack manuel sur canal dédié
		}
	}
}

func (o *Engine) BroadcastToUI(msgType string, data interface{}) {
	if o.Hub != nil {
		o.Hub.BroadcastMessage(map[string]interface{}{
			"type": msgType,
			"data": data,
		})
	}
}

func (o *Engine) StartGarbageCollector() {
	ticker := time.NewTicker(30 * time.Second)
	o.safeGo(func() {
		for range ticker.C {
			o.Mu.Lock()
			now := time.Now().Unix()
			for id, worker := range o.Workers {
				if now-worker.LastSeen > 60 {
					log.Printf("[GC] Worker %s inactif. Nettoyage...", id)

					// 1. Suppression de la queue durable sur RabbitMQ
					if o.Channel != nil {
						queueName := "results_" + id
						_, err := o.Channel.QueueDelete(queueName, false, false, false)
						if err != nil {
							log.Printf("[GC] Erreur suppression queue %s: %v", queueName, err)
						}
					}

					// 2. Nettoyage mémoire et UI
					delete(o.Workers, id)
					o.BroadcastToUI("WORKER_LEAVE", map[string]string{"id": id})
				}
			}
			o.Mu.Unlock()
		}
	})
}
