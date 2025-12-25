package worker

import (
	"encoding/json"
	"log"
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
	defer e.Mu.Unlock()

	reg, _ := json.Marshal(WorkerRegistration{ID: e.ID, Language: "go"})
	err := e.Channel.Publish(
		"",
		"isReady",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        reg,
		},
	)
	if err != nil {
		log.Printf("[ERROR] Échec de l'enregistrement: %v", err)
	} else {
		log.Printf("[PHASE 1] Worker enregistré avec l'ID: %s", e.ID)
	}
}
