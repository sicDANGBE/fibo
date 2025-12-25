package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// InitRabbitMQ amorce la boucle de connexion résiliente
func (o *Engine) InitRabbitMQ(url string) {
	go o.handleReconnect(url)
}

// handleReconnect assure la survie de la connexion sur ton cluster k3s [cite: 2025-12-20]
func (o *Engine) handleReconnect(url string) {
	for {
		log.Printf("[RMQ] Tentative de connexion à %s...", url)

		conn, err := amqp.Dial(url)
		if err != nil {
			log.Printf("[RMQ] Connexion impossible: %v. Nouvel essai dans 15s...", err)
			time.Sleep(15 * time.Second)
			continue
		}

		o.Conn = conn
		ch, err := conn.Channel()
		if err != nil {
			log.Printf("[RMQ] Canal impossible: %v", err)
			conn.Close()
			time.Sleep(15 * time.Second)
			continue
		}

		o.Mu.Lock()
		o.Channel = ch
		o.Mu.Unlock()

		// Configuration des exchanges et queues de base
		if err := o.setupInfrastructure(); err != nil {
			log.Printf("[RMQ] Erreur setup infrastructure: %v. Reconnexion dans 15s...", err)
			conn.Close()
			time.Sleep(15 * time.Second)
			continue
		}

		log.Println("[RMQ] Connecté et infrastructure prête.")

		// Notification de fermeture pour déclencher la reconnexion automatique
		closeChan := make(chan *amqp.Error)
		o.Conn.NotifyClose(closeChan)

		// Lancement de l'écoute des workers (Phase 1)
		go o.ListenForWorkers()

		// Blocage jusqu'à la déconnexion
		err = <-closeChan
		log.Printf("[RMQ] Rupture de flux: %v. Relance de la procédure de résilience...", err)
	}
}

func (o *Engine) setupInfrastructure() error {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	// Exchange Fanout pour la synchronisation synchrone (Phase 3)
	err := o.Channel.ExchangeDeclare(
		"fibo_admin_exchange",
		"fanout",
		true, // Durable pour la persistance dans k3s [cite: 2025-12-20]
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Queue isReady pour Phase 1
	// Changement : durable: true pour éviter l'erreur PRECONDITION_FAILED (406)
	_, err = o.Channel.QueueDeclare(
		"isReady",
		true, // Durable : doit correspondre à la queue existante sur ton RMQ
		false,
		false,
		false,
		nil,
	)
	return err
}

func (o *Engine) ListenForWorkers() {
	o.Mu.Lock()
	if o.Channel == nil {
		o.Mu.Unlock()
		return
	}
	ch := o.Channel
	o.Mu.Unlock()

	msgs, err := ch.Consume("isReady", "", true, false, false, false, nil)
	if err != nil {
		log.Printf("[RMQ] Erreur Consume isReady: %v", err)
		return
	}

	for d := range msgs {
		var reg WorkerRegistration
		if err := json.Unmarshal(d.Body, &reg); err != nil {
			log.Printf("[RMQ] Erreur décodage: %v", err)
			continue
		}

		o.Mu.Lock()
		o.Workers[reg.ID] = reg

		// Phase 2 : Queue de résultats unique
		resQueue := fmt.Sprintf("results_%s", reg.ID)
		args := amqp.Table{
			"x-queue-type": "stream",
		}
		_, err := o.Channel.QueueDeclare(
			resQueue,
			true,  // Durable : obligatoire pour les streams
			false, // Auto-delete : non supporté par les streams
			false,
			false,
			args,
		)

		if err == nil {
			go o.ConsumeWorkerResults(resQueue)
			log.Printf("[SYNC] Worker %s (%s) enregistré via %s", reg.ID, reg.Language, resQueue)
			o.BroadcastToUI("WORKER_JOIN", reg)
		} else {
			log.Printf("[RMQ] Erreur création queue %s: %v", resQueue, err)
		}
		o.Mu.Unlock()
	}
}
