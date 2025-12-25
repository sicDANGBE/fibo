package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Engine struct {
	ID      string
	Conn    *amqp.Connection
	Channel *amqp.Channel
	AMQPURL string
	Mu      sync.Mutex
}

func NewEngine(url string) *Engine {
	return &Engine{
		ID:      GenerateID(),
		AMQPURL: url,
	}
}

func (e *Engine) Start() {
	for {
		log.Printf("[WORKER] Tentative de connexion RMQ sur %s", e.AMQPURL)
		conn, err := amqp.Dial(e.AMQPURL)
		if err != nil {
			log.Printf("[ERROR] Échec connexion: %v. Re-tentative dans 15s...", err)
			time.Sleep(15 * time.Second)
			continue
		}

		e.Conn = conn
		ch, err := conn.Channel()
		if err != nil {
			log.Printf("[ERROR] Impossible d'ouvrir un canal: %v", err)
			conn.Close()
			time.Sleep(15 * time.Second)
			continue
		}

		e.Mu.Lock()
		e.Channel = ch
		e.Mu.Unlock()

		// Configuration de l'infrastructure minimale pour le worker
		if err := e.setupInfra(); err != nil {
			log.Printf("[ERROR] Erreur setup infra: %v", err)
			conn.Close()
			time.Sleep(15 * time.Second)
			continue
		}

		// Phase 1 : Signalement de présence
		e.register()
		// Démarrage du heartbeat
		go e.startHeartbeat()

		// Phase 2 : Écoute des ordres (Exchange Fanout)
		go e.listenTasks()

		// Surveillance de la santé de la connexion
		closeChan := make(chan *amqp.Error)
		e.Conn.NotifyClose(closeChan)

		log.Println("[WORKER] Connecté et prêt à recevoir des tâches.")

		err = <-closeChan
		log.Printf("[WARN] Connexion perdue: %v. Relance de la boucle...", err)
	}
}

func (e *Engine) setupInfra() error {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	// Déclaration de la queue de présence (doit être identique à l'orchestrateur)
	_, err := e.Channel.QueueDeclare(
		"isReady",
		true, // Durable: true (match orchestrateur)
		false,
		false,
		false,
		nil,
	)
	return err
}

func (e *Engine) register() {
	e.Mu.Lock()
	ch := e.Channel
	defer e.Mu.Unlock()

	// 1. Créer la file d'attente pour la confirmation (Auto-delete pour le nettoyage)
	ackQueue := fmt.Sprintf("ack_%s", e.ID)
	_, err := ch.QueueDeclare(ackQueue, false, true, false, false, nil)
	if err != nil {
		log.Fatalf("[CRITICAL] Impossible de créer la file ACK: %v", err)
	}

	// 2. Envoyer la demande d'enregistrement
	reg, _ := json.Marshal(WorkerRegistration{ID: e.ID, Language: "go"})
	ch.Publish("", "isReady", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        reg,
	})

	// 3. Bloquer jusqu'à réception du signal "READY" de l'orchestrateur
	msgs, _ := ch.Consume(ackQueue, "", true, false, false, false, nil)
	log.Println("[WAIT] En attente de la validation de l'orchestrateur...")

	for d := range msgs {
		if string(d.Body) == "READY" {
			log.Printf("[PHASE 1] Enregistrement validé. Stream prêt.")
			break
		}
	}

	// Nettoyage de la file ACK
	ch.QueueDelete(ackQueue, false, false, false)

	if err != nil {
		log.Printf("[ERROR] Échec de l'enregistrement: %v", err)
	} else {
		log.Printf("[PHASE 1] Worker enregistré avec l'ID: %s", e.ID)
	}
}

func (e *Engine) startHeartbeat() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		health := map[string]interface{}{
			"worker_id": e.ID,
			"status":    "active",
			"ram":       m.Alloc / 1024 / 1024,
			"cpu":       runtime.NumGoroutine(),
			"timestamp": time.Now().Unix(),
			// Simuler I/O pour l'instant (à lier à /proc/net/dev plus tard)
			"net_io":  "2.4MB/s",
			"disk_io": "150KB/s",
		}

		body, _ := json.Marshal(health)
		e.Mu.Lock()
		if e.Channel != nil {
			e.Channel.Publish("", "worker_health", false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		}
		e.Mu.Unlock()
	}
}
