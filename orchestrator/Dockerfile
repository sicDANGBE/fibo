# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache ca-certificates git
WORKDIR /src

# Cache des dépendances
COPY go.mod go.su* ./
RUN go mod download

# Copie de tout le code pour permettre l'embed
COPY . .

# Build statique (Correction du chemin vers cmd/server)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/orchestrator ./cmd/server/main.go

# --- Stage 2: Image Finale ---
FROM scratch
# Import des certificats pour RabbitMQ TLS [cite: 2025-12-02]
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /

# Le binaire contient l'IHM grâce à go:embed
COPY --from=builder /app/orchestrator .

EXPOSE 8080
ENTRYPOINT ["./orchestrator"]