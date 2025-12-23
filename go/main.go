package main

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	// L'import doit correspondre au nom dans go.mod + le dossier de destination
	pb "fibonacci/pb"
)

type server struct {
	pb.UnimplementedBarrierServer
	mu         sync.Mutex
	cond       *sync.Cond
	readyCount int
}

// WaitToStart bloque les clients jusqu'à ce que les 2 (Node + Python) soient là
func (s *server) WaitToStart(ctx context.Context, in *pb.Empty) (*pb.StartSignal, error) {
	s.mu.Lock()
	s.readyCount++
	fmt.Printf("[gRPC] Client connecté. Total : %d/2\n", s.readyCount)
	if s.readyCount >= 2 {
		s.cond.Broadcast()
	} else {
		s.cond.Wait()
	}
	s.mu.Unlock()
	return &pb.StartSignal{Message: "Signal de départ reçu !"}, nil
}

func runFibo() {
	for run := 1; run <= 10; run++ {
		fmt.Printf("[GO] --- Démarrage Run %d/10 ---\n", run)
		a, b := big.NewInt(0), big.NewInt(1)
		start := time.Now()

		for i := 0; i <= 400000; i++ {
			a.Add(a, b)
			a, b = b, a

			if i%10000 == 0 && i > 0 {
				fmt.Printf("[GO] Run %d | %d itérations | Temps écoulé: %v\n", run, i, time.Since(start))
			}
		}
		fmt.Printf("[GO] Run %d terminé en %v\n", run, time.Since(start))
	}
}

func main() {
	// 1. Initialisation du serveur gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	s := &server{}
	s.cond = sync.NewCond(&s.mu)
	grpcServer := grpc.NewServer()
	pb.RegisterBarrierServer(grpcServer, s)

	// 2. Lancement du serveur en arrière-plan
	go func() {
		fmt.Println("[GO] Serveur gRPC à l'écoute sur :50051")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// 3. Attente de la barrière de synchronisation pour le worker GO lui-même
	fmt.Println("[GO] En attente de Node.js et Python...")
	s.mu.Lock()
	for s.readyCount < 2 {
		s.cond.Wait()
	}
	s.mu.Unlock()

	// 4. Lancement du calcul
	fmt.Println("[GO] Signal reçu. Lancement des calculs !")
	runFibo()

	// On laisse un peu de temps pour que les autres finissent avant de couper le container
	time.Sleep(30 * time.Second)
}
