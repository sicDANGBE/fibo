# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache protoc protobuf-dev git ca-certificates

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

# 2. Génération du code gRPC
COPY proto/ ./proto/
# On crée le dossier de destination
RUN mkdir -p internal/pb 

# Commande protoc alignée sur le module "fibo-worker" 
RUN protoc --proto_path=./proto \
    --go_out=. --go_opt=module=fibo-worker \
    --go-grpc_out=. --go-grpc_opt=module=fibo-worker \
    ./proto/sync.proto

# 3. Copie du code source et Build
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/worker ./cmd/worker/main.go

# --- Stage 2: Image Finale (Scratch pour ton cluster k3s) ---
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /
COPY --from=builder /app/worker .

EXPOSE 8081 50051
ENTRYPOINT ["./worker"]