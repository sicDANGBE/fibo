package main

import (
	"fibo-orchestrateur/internal/api"
	"fibo-orchestrateur/internal/orchestrator"
	"os"
)

func main() {
	// 1. Initialisation du Hub WebSocket
	hub := api.NewHub()
	go hub.Run()

	// 2. Initialisation de l'Engine (RabbitMQ + Orchestration)
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://bench_leader:qsd65f4c98dc7fd9s87ga6fsd5g4zsdrf9g879dfs7g@192.168.1.12:5672/benchmarks"
	}

	// NewEngine lance la boucle de reconnexion en interne
	orch := orchestrator.NewEngine(amqpURL, hub)

	// 3. Setup et Lancement du serveur Web
	r := api.SetupRouter(orch, hub)
	r.Run(":8080")
}
