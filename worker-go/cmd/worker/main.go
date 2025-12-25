package main

import (
	"fibo-worker/internal/api"
	"fibo-worker/internal/worker"
	"os"
)

func main() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://bench_go:9q8s7d9qs87dqs654dq6s54d6qs54dqs321dqs2d1qs98d7qs9d8q7@192.168.1.12:5672/benchmarks"
	}

	// 1. Lancement du moteur Worker (Async)
	engine := worker.NewEngine(amqpURL)
	go engine.Start()

	// 2. Lancement de l'API de santé et IO (Sync)
	r := api.SetupRouter()
	r.Run(":8081")
}
