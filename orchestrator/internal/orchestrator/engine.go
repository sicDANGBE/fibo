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

// ConsumeWorkerResults traite les messages entrants de chaque langage (Phase 4)
func (o *Engine) ConsumeWorkerResults(queueName string) {
	o.Mu.Lock()
	ch := o.Channel
	o.Mu.Unlock()

	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Printf("[RMQ] Erreur consommation %s: %v", queueName, err)
		return
	}

	for d := range msgs {
		var res WorkerResult
		if err := json.Unmarshal(d.Body, &res); err != nil {
			continue
		}
		// Envoi immédiat vers l'UI via le hub WebSocket
		o.BroadcastToUI("RESULT", res)
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
