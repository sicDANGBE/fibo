# Export de projet

_Généré le 2025-12-23T00:40:11+01:00_

## compose.yml

```yaml
services:
  fibo-go:
    build:
      context: .
      dockerfile: go/Dockerfile
    container_name: fibo-go
    deploy:
      resources:
        limits:
          cpus: '5.0'
          memory: 5G
    networks: [fibo-net]

  fibo-python:
    build:
      context: .
      dockerfile: python/Dockerfile
    depends_on: [fibo-go]
    deploy:
      resources:
        limits:
          cpus: '5.0'
          memory: 5G
    networks: [fibo-net]

  fibo-node:
    build:
      context: .
      dockerfile: node/Dockerfile
    depends_on: [fibo-go]
    deploy:
      resources:
        limits:
          cpus: '5.0'
          memory: 5G
    networks: [fibo-net]

networks:
  fibo-net:
    driver: bridge
```

## go/Dockerfile

```text
FROM golang:1.25-alpine AS builder

# Installation de protoc et plugins
RUN apk add --no-cache protoc protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /src

# Copie des dépendances
COPY go/go.mod ./
# On force le téléchargement des dépendances grpc
RUN go mod download

# Création du dossier pb et génération
RUN mkdir -p pb
COPY proto/sync.proto ./proto/
RUN protoc --proto_path=./proto --go_out=. --go-grpc_out=. ./proto/sync.proto

# Copie du main et build
COPY go/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/fibo main.go

# Image finale légère
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/fibo .
# On expose le port gRPC
EXPOSE 50051
ENTRYPOINT ["./fibo"]
```

## go/go.mod

```text
module fibonacci

go 1.25.1

require google.golang.org/grpc v1.77.0

require (
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

```

## go/go.sum

```text
github.com/go-logr/logr v1.4.3 h1:CjnDlHq8ikf6E492q6eKboGOC0T8CDaOvkHCIg8idEI=
github.com/go-logr/logr v1.4.3/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
go.opentelemetry.io/auto/sdk v1.2.1 h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=
go.opentelemetry.io/auto/sdk v1.2.1/go.mod h1:KRTj+aOaElaLi+wW1kO/DZRXwkF4C5xPbEe3ZiIhN7Y=
go.opentelemetry.io/otel v1.38.0 h1:RkfdswUDRimDg0m2Az18RKOsnI8UDzppJAtj01/Ymk8=
go.opentelemetry.io/otel v1.38.0/go.mod h1:zcmtmQ1+YmQM9wrNsTGV/q/uyusom3P8RxwExxkZhjM=
go.opentelemetry.io/otel/metric v1.38.0 h1:Kl6lzIYGAh5M159u9NgiRkmoMKjvbsKtYRwgfrA6WpA=
go.opentelemetry.io/otel/metric v1.38.0/go.mod h1:kB5n/QoRM8YwmUahxvI3bO34eVtQf2i4utNVLr9gEmI=
go.opentelemetry.io/otel/sdk v1.38.0 h1:l48sr5YbNf2hpCUj/FoGhW9yDkl+Ma+LrVl8qaM5b+E=
go.opentelemetry.io/otel/sdk v1.38.0/go.mod h1:ghmNdGlVemJI3+ZB5iDEuk4bWA3GkTpW+DOoZMYBVVg=
go.opentelemetry.io/otel/sdk/metric v1.38.0 h1:aSH66iL0aZqo//xXzQLYozmWrXxyFkBJ6qT5wthqPoM=
go.opentelemetry.io/otel/sdk/metric v1.38.0/go.mod h1:dg9PBnW9XdQ1Hd6ZnRz689CbtrUp0wMMs9iPcgT9EZA=
go.opentelemetry.io/otel/trace v1.38.0 h1:Fxk5bKrDZJUH+AMyyIXGcFAPah0oRcT+LuNtJrmcNLE=
go.opentelemetry.io/otel/trace v1.38.0/go.mod h1:j1P9ivuFsTceSWe1oY+EeW3sc+Pp42sO++GHkg4wwhs=
golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 h1:6/3JGEh1C88g7m+qzzTbl3A0FtsLguXieqofVLU/JAo=
golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82/go.mod h1:Q9BGdFy1y4nkUwiLvT5qtyhAnEHgnQ/zd8PfU6nc210=
golang.org/x/sys v0.37.0 h1:fdNQudmxPjkdUTPnLn5mdQv7Zwvbvpaxqs831goi9kQ=
golang.org/x/sys v0.37.0/go.mod h1:OgkHotnGiDImocRcuBABYBEXf8A9a87e/uXjp9XT3ks=
golang.org/x/text v0.30.0 h1:yznKA/E9zq54KzlzBEAWn1NXSQ8DIp/NYMy88xJjl4k=
golang.org/x/text v0.30.0/go.mod h1:yDdHFIX9t+tORqspjENWgzaCVXgk0yYnYuSZ8UzzBVM=
gonum.org/v1/gonum v0.16.0 h1:5+ul4Swaf3ESvrOnidPp4GZbzf0mxVQpDCYUQE7OJfk=
gonum.org/v1/gonum v0.16.0/go.mod h1:fef3am4MQ93R2HHpKnLk4/Tbh/s0+wqD5nfa6Pnwy4E=
google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 h1:M1rk8KBnUsBDg1oPGHNCxG4vc1f49epmTO7xscSajMk=
google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8/go.mod h1:7i2o+ce6H/6BluujYR+kqX3GKH+dChPTQU19wjRPiGk=
google.golang.org/grpc v1.77.0 h1:wVVY6/8cGA6vvffn+wWK5ToddbgdU3d8MNENr4evgXM=
google.golang.org/grpc v1.77.0/go.mod h1:z0BY1iVj0q8E1uSQCjL9cppRj+gnZjzDnzV0dHhrNig=
google.golang.org/protobuf v1.36.10 h1:AYd7cD/uASjIL6Q9LiTjz8JLcrh/88q5UObnmY3aOOE=
google.golang.org/protobuf v1.36.10/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=

```

## go/main.go

```go
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

```

## go/pb/sync.pb.go

```go

```

## go/pb/sync_grpc.pb.go

```go

```

## main.go

```go
package main

import (
	"fmt"
	"math/big"
)

// FibGenerator produit la suite de Fibonacci via un channel.
// S'arrête si la RAM estimée des deux derniers nombres dépasse limitBytes.
func FibGenerator(limitBytes uint64) <-chan *big.Int {
	ch := make(chan *big.Int)

	go func() {
		defer close(ch)
		// Initialisation : F(0)=0, F(1)=1
		a := big.NewInt(0)
		b := big.NewInt(1)

		for {
			// Calcul de la taille approximative en RAM des deux termes
			// big.Int stocke les données dans un slice de 'Word' (uint sur 64 bits)
			// On compte environ 8 octets par Word + le overhead de la structure.
			sizeA := uint64(len(a.Bits())) * 8
			sizeB := uint64(len(b.Bits())) * 8

			if sizeA+sizeB > limitBytes {
				fmt.Printf("\n[Limite de %d Go atteinte]\n", limitBytes/1024/1024/1024)
				return
			}

			// On envoie une copie pour éviter les effets de bord si l'appelant modifie la valeur
			val := new(big.Int).Set(a)
			ch <- val

			// Fibonacci : a, b = b, a+b
			// On utilise Add pour additionner b à a, puis on swap.
			a.Add(a, b)
			a, b = b, a
		}
	}()

	return ch
}

func main() {
	const maxRAM = 5 * 1024 * 1024 * 1024 // 5 Go
	gen := FibGenerator(maxRAM)

	count := 0
	for f := range gen {
		count++
		// Pour l'exemple, on affiche tous les 100 000 termes
		// car l'affichage console est très lent pour de gros nombres.
		if count%100000 == 0 {
			fmt.Printf("Terme n°%d calculé (Taille actuelle : ~%d MB)\n", count, len(f.Bits())*8/1024/1024)
		}
	}
}

```

## node/Dockerfile

```text
FROM node:20-alpine

WORKDIR /app

# Installation des dépendances
RUN npm install @grpc/grpc-js @grpc/proto-loader

# Copie du proto et du script
COPY ../proto/sync.proto .
COPY index.js .

CMD ["node", "index.js"]
```

## node/index.js

```javascript
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const packageDefinition = protoLoader.loadSync('sync.proto');
const syncProto = grpc.loadPackageDefinition(packageDefinition).sync;

function runFibo() {
    for (let run = 1; run <= 10; run++) {
        let a = 0n, b = 1n;
        const start = Date.now();
        for (let i = 0; i <= 400000; i++) {
            [a, b] = [b, a + b];
            if (i % 10000 === 0 && i > 0) {
                console.log(`[NODE] Run ${run} - ${i} iters - Temps: ${(Date.now() - start)/1000}s`);
            }
        }
    }
}

const client = new syncProto.Barrier('fibo-go:50051', grpc.credentials.createInsecure());
console.log("Node prêt, en attente du signal...");
client.waitToStart({}, (err) => {
    if (err) console.error(err);
    else runFibo();
});
```

## project_export.log

```text
[2025-12-23 00:40:11] Source  : .
[2025-12-23 00:40:11] Sortie  : project_export.md
[2025-12-23 00:40:11] Fichiers trouvés (avant filtre): 15
[2025-12-23 00:40:11] Fichiers à concaténer (après filtre): 14 (exclus auto:1 dir:0 file:0)
[2025-12-23 00:40:11] Concatène [1] compose.yml (size=697)
[2025-12-23 00:40:11] Concatène [2] go/Dockerfile (size=802)
[2025-12-23 00:40:11] Concatène [3] go/go.mod (size=365)
[2025-12-23 00:40:11] Concatène [4] go/go.sum (size=3182)
[2025-12-23 00:40:11] Concatène [5] go/main.go (size=2057)
[2025-12-23 00:40:11] Concatène [6] go/pb/sync.pb.go (size=0)
[2025-12-23 00:40:11] Concatène [7] go/pb/sync_grpc.pb.go (size=0)
[2025-12-23 00:40:11] Concatène [8] main.go (size=1476)
[2025-12-23 00:40:11] Concatène [9] node/Dockerfile (size=215)
[2025-12-23 00:40:11] Concatène [10] node/index.js (size=836)

```

## proto/sync.proto

```text
syntax = "proto3";

package sync;
option go_package = "./pb";

service Barrier {
  // Les clients appellent ceci et attendent la réponse pour démarrer
  rpc WaitToStart (Empty) returns (StartSignal);
}

message Empty {}
message StartSignal {
  string message = 1;
}
```

## python/Dockerfile

```text
FROM python:3.12-slim

WORKDIR /app

# Installation des dépendances gRPC
RUN pip install --no-cache-dir grpcio grpcio-tools

# Copie du proto
COPY ../proto/sync.proto .

# Génération du code Python
RUN python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. sync.proto

# Copie du script
COPY main.py .

CMD ["python", "main.py"]
```

## python/main.py

```python
import grpc
import time
from sync_pb2 import Empty
from sync_pb2_grpc import BarrierStub

def run_fibo():
    for run in range(1, 11):
        a, b = 0, 1
        start_time = time.time()
        for i in range(400001):
            a, b = b, a + b
            if i % 10000 == 0 and i > 0:
                print(f"[PYTHON] Run {run} - {i} iters - Temps: {time.time() - start_time:.4f}s")

if __name__ == "__main__":
    with grpc.insecure_channel('fibo-go:50051') as channel:
        stub = BarrierStub(channel)
        print("Python prêt, en attente du signal...")
        stub.WaitToStart(Empty())
        run_fibo()
```

