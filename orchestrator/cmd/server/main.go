package main

import (
	"context"
	"fibo-orchestrateur/internal/api"
	"fibo-orchestrateur/internal/orchestrator"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. Initialisation du Hub WebSocket
	hub := api.NewHub()
	go hub.Run()

	// 2. Configuration de l'URL RabbitMQ via .env
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://bench_leader:qsd65f4c98dc7fd9s87ga6fsd5g4zsdrf9g879dfs7g@192.168.1.12:5672/benchmarks"
	}

	// 3. Lancement de l'Engine (RabbitMQ + Orchestration)
	// NewEngine lance déjà la boucle de reconnexion en interne
	orch := orchestrator.NewEngine(amqpURL, hub)

	// 4. Configuration du serveur HTTP avec Graceful Shutdown
	router := api.SetupRouter(orch, hub)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Lancement du serveur dans une goroutine pour ne pas bloquer le thread principal
	go func() {
		log.Println("[API] Serveur démarré sur :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[CRITICAL] Erreur serveur Web: %v", err)
		}
	}()

	// 5. Gestion des Signaux d'Arrêt (SIGINT, SIGTERM)
	// Indispensable pour k3s lors d'un redéploiement
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Blocage jusqu'à réception d'un signal
	sig := <-quit
	log.Printf("[MAIN] Signal %v reçu. Début de la procédure d'arrêt propre...", sig)

	// 6. Procédure de fermeture Graceful (Timeout de 5 secondes)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Arrêt du serveur HTTP (ne prend plus de nouvelles requêtes)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("[ERROR] Arrêt forcé du serveur HTTP:", err)
	}

	// Fermeture des ressources critiques
	log.Println("[RMQ] Fermeture des connexions RabbitMQ...")
	if orch.Conn != nil {
		orch.Conn.Close()
	}

	log.Println("[MAIN] Orchestrateur arrêté proprement. Bye!")
}
