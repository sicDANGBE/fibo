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

// handleReconnect assure la survie de la connexion sur ton cluster k3s
func (o *Engine) handleReconnect(url string) {
	for {
		log.Printf("[RMQ] Tentative de connexion à %s...", url)
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Printf("[RMQ] Échec Dial: %v. Re-tentative dans 10s...", err)
			time.Sleep(10 * time.Second)
			continue
		}

		o.Conn = conn
		ch, err := conn.Channel()
		if err != nil {
			log.Printf("[RMQ] Échec Canal: %v", err)
			conn.Close()
			time.Sleep(10 * time.Second)
			continue
		}

		o.Mu.Lock()
		o.Channel = ch
		o.Mu.Unlock()

		if err := o.setupInfrastructure(); err != nil {
			log.Printf("[RMQ] Erreur Infrastructure: %v", err)
			conn.Close()
			time.Sleep(10 * time.Second)
			continue
		}

		log.Println("[RMQ] Connecté et infrastructure prête.")

		// --- Surveillance double : Connexion + Canal ---
		notifyConnClose := make(chan *amqp.Error)
		o.Conn.NotifyClose(notifyConnClose)

		notifyChanClose := make(chan *amqp.Error)
		o.Channel.NotifyClose(notifyChanClose)

		// Lancement des écoutes
		go o.ListenForWorkers()
		go o.ListenForHealth()

		// Blocage intelligent
		select {
		case err := <-notifyConnClose:
			log.Printf("[RMQ] Connexion perdue (vhost/server): %v", err)
		case err := <-notifyChanClose:
			log.Printf("[RMQ] Canal perdu (queue/protocol error): %v", err)
		}

		// Nettoyage avant reconnexion
		o.Mu.Lock()
		o.Channel = nil
		o.Mu.Unlock()
		o.Conn.Close()

		log.Println("[RMQ] Nettoyage effectué. Re-tentative de connexion...")
		time.Sleep(5 * time.Second)
	}
}

func (o *Engine) setupInfrastructure() error {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	// Exchange Fanout pour la synchronisation synchrone (Phase 3)
	err := o.Channel.ExchangeDeclare(
		"fibo_admin_exchange",
		"fanout",
		true, // Durable pour la persistance dans k3s
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
	// Déclaration de la queue de santé (Classic Durable)
	_, err = o.Channel.QueueDeclare("worker_health", true, false, false, false, nil)
	return err
}

func (o *Engine) ListenForHealth() {
	o.Mu.Lock()
	ch := o.Channel
	o.Mu.Unlock()

	// Sécurité supplémentaire : Nil check
	if ch == nil {
		log.Println("[ERROR] ListenForHealth abandonné : canal nil")
		return
	}

	msgs, err := ch.Consume("worker_health", "", true, false, false, false, nil)
	if err != nil {
		log.Printf("[RMQ] Erreur Consume health: %v", err)
		return
	}

	for d := range msgs {
		var healthData struct {
			WorkerID string `json:"worker_id"`
			RAM      int    `json:"ram"`
			CPU      int    `json:"cpu"`
		}
		if err := json.Unmarshal(d.Body, &healthData); err == nil {
			o.Mu.Lock()
			if w, ok := o.Workers[healthData.WorkerID]; ok {
				w.LastSeen = time.Now().Unix() // Rafraîchissement
				o.Workers[healthData.WorkerID] = w
			}
			o.Mu.Unlock()
			o.BroadcastToUI("HEALTH_UPDATE", healthData)
		}
	}
}

func (o *Engine) ListenForWorkers() {
	o.Mu.Lock()
	if o.Channel == nil {
		o.Mu.Unlock()
		return
	}
	ch := o.Channel // On extrait le canal pour ne pas garder le lock durant le Consume
	o.Mu.Unlock()

	msgs, err := ch.Consume("isReady", "", true, false, false, false, nil)
	if err != nil {
		log.Printf("[RMQ] Erreur Consume: %v", err)
		return
	}

	for d := range msgs {
		var reg WorkerRegistration
		if err := json.Unmarshal(d.Body, &reg); err != nil {
			continue
		}

		// --- SCOPE DE VERROUILLAGE ATOMIQUE ---
		o.Mu.Lock()
		reg.LastSeen = time.Now().Unix() // Initialisation du timestamp pour le GC
		o.Workers[reg.ID] = reg
		o.Mu.Unlock() // ON LIBÈRE IMMÉDIATEMENT

		// Création de la queue Stream sur un canal dédié pour éviter l'erreur 503
		resQueue := fmt.Sprintf("results_%s", reg.ID)
		ackQueue := fmt.Sprintf("ack_%s", reg.ID) // File de réponse pour le worker

		// Utilisation d'un canal temporaire pour la déclaration (DevOps Best Practice)
		tmpCh, err := o.Conn.Channel()
		if err != nil {
			log.Printf("[RMQ] Erreur canal temporaire: %v", err)
			continue
		}

		_, err = tmpCh.QueueDeclare(
			resQueue,
			true, // Durable obligatoire pour les Streams
			false,
			false,
			false,
			amqp.Table{"x-queue-type": "stream", "x-max-age": "1h"}, // 7 jours en ms
		)
		if err != nil {
			log.Printf("[ERROR] Échec création Stream %s: %v", resQueue, err)
			tmpCh.Close()
			continue
		}

		// 2. Confirmation au worker que tout est prêt
		err = tmpCh.Publish(
			"",
			ackQueue,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("READY"),
			},
		)

		log.Printf("[SYNC] Queue %s créée et Worker %s notifié.", resQueue, reg.ID)
		tmpCh.Close()

		if err == nil {
			o.safeGo(func() {
				o.ConsumeWorkerResults(resQueue)
			})
			o.BroadcastToUI("WORKER_JOIN", reg) //
			log.Printf("[SYNC] Worker %s (%s) prêt.", reg.ID, reg.Language)
		}
	}
}
