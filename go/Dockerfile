FROM golang:1.25-alpine AS builder

RUN apk add --no-cache protoc protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /src

# 1. Préparation des modules
COPY go/go.mod ./
RUN go mod download

# 2. Génération gRPC
# On copie le proto
COPY proto/sync.proto ./proto/

# On génère le code en disant à protoc que la racine du module est '.'
# Cela va créer automatiquement le dossier /src/pb/
RUN mkdir -p pb
RUN protoc --proto_path=./proto \
    --go_out=. --go_opt=module=fibonacci \
    --go-grpc_out=. --go-grpc_opt=module=fibonacci \
    ./proto/sync.proto

# 3. Build du binaire
COPY go/main.go .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/fibo main.go

# --- Image finale ---
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/fibo .
EXPOSE 50051
ENTRYPOINT ["./fibo"]