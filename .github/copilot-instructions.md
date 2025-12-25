# Copilot / Agent Instructions for Fibonnaci Benchmark

Quick, actionable notes to get an AI coding agent productive in this repository.

1) Big-picture architecture
- **Orchestrator (Go)**: `orchestrator/` runs an HTTP server + WebSocket hub and owns RabbitMQ orchestration. Entry: `orchestrator/cmd/server/main.go`.
- **Workers**: language-specific worker folders at `worker-go/`, `worker-node/`, `worker-python/`, `worker-rust/`. Each worker connects to RabbitMQ and participates in the benchmark.
- **Message bus**: RabbitMQ is the central integration point. Key code: `orchestrator/internal/orchestrator/rabbitmq.go`.
- **UI**: static files and client code under `orchestrator/internal/api/web/` and WebSocket hub in `orchestrator/internal/api/hub.go`.
- **Compose / orchestration**: `compose.yml` defines services, environment variables and network `fibo-benchmark-network`.

2) Core dataflows & conventions (must-follow)
- Worker registration: workers announce themselves to queue `isReady` (see `ListenForWorkers()` in `rabbitmq.go`).
- Per-worker result queues follow the pattern `results_{workerID}`. These are created as RabbitMQ streams (queue arg `x-queue-type: stream`). See the `QueueDeclare` call in `rabbitmq.go`.
- Admin sync: a durable fanout exchange `fibo_admin_exchange` is declared for synchronous broadcast across services.
- Durability: queues/exchanges are declared `durable: true` in orchestrator — when adding or changing queues follow the same durability/stream semantics to avoid PRECONDITION_FAILED errors.
- Concurrency: orchestrator uses `Engine.Mu` to guard `Channel`/`Workers`. Follow locking patterns when mutating shared Engine state.

3) Important files to reference
- RabbitMQ orchestration: `orchestrator/internal/orchestrator/rabbitmq.go`
- Engine and types: `orchestrator/internal/orchestrator/types.go`
- WebSocket hub & routes: `orchestrator/internal/api/hub.go`, `orchestrator/internal/api/routes.go` and `orchestrator/internal/api/web/`
- Docker compose: `compose.yml` (service names, env var names like `AMQP_URL_*`, `REDIS_URL`)
- Worker examples: `worker-go/`, `worker-node/`, `worker-python/`, `worker-rust/` (follow folder layout and Dockerfiles)
- Proto definitions: `proto/sync.proto` (check if type generation is needed before cross-language changes)

4) Developer workflows (how to build/run/debug)
- Full stack (recommended): from repo root
  - Build & start: `docker compose up -d --force-recreate --build --remove-orphans`
  - Tail orchestrator logs: `docker compose logs -f orchestrator`
  - Start a single service: `docker compose up -d orchestrator` or `docker compose up -d fibo-go`
- Local Go iteration (without Docker): set `AMQP_URL` locally then run:
  - `go run ./orchestrator/cmd/server` (or `go build ./...` then run binary)
  - The orchestrator listens on `:8080` and exposes health at `/health`.
- Worker iteration:
  - Go worker: build/run under `worker-go/` (`go run ./cmd/worker`)
  - Node worker: `node worker-node/index.js`
  - Python worker: `python3 worker-python/main.py`
  - Rust worker: `cargo build` / `cargo run` in `worker-rust/`

5) Environment variables and defaults
- `AMQP_URL` (or variants in `compose.yml`: `AMQP_URL_LEADER`, `AMQP_URL_GO`, `AMQP_URL_NODE`, `AMQP_URL_PYTHON`, `AMQP_URL_RUST`) — orchestrator falls back to an embedded URL in `orchestrator/cmd/server/main.go` if unset.
- `REDIS_URL` referenced in compose for caching/coordination.

6) When changing message schemas or queues
- Update `proto/sync.proto` first if the change affects message formats used by multiple languages; regenerate stubs if the repo uses codegen (no generator found in repo—confirm with maintainers).
- Keep queue names and durability consistent with `rabbitmq.go` (`isReady`, `results_{id}`, `fibo_admin_exchange`). Streams require declaring `x-queue-type: stream`.

7) Code patterns the agent should follow
- Keep resilient reconnect logic intact: `Engine.InitRabbitMQ` → `handleReconnect` is responsible for reconnection and re-declaring infra.
- Use `Engine.Mu` when accessing `Engine.Channel` or `Engine.Workers`.
- When adding a worker implementation, mirror the `results_{id}` creation and support the `isReady` registration flow.
- Prefer the existing logging style (`log.Printf` with tags like `[RMQ]`, `[SYNC]`, `[WORKER]`).

8) Quick examples (copyable)
- Declare a stream results queue (Go):
```go
args := amqp.Table{"x-queue-type":"stream"}
_, _ = ch.QueueDeclare("results_<id>", true, false, false, false, args)
```
- Broadcast to UI via hub:
```go
hub.BroadcastMessage(map[string]interface{"type": "WORKER_JOIN", "data": reg})
```

9) What the agent should ask the maintainers (if unclear)
- Are stubs generated from `proto/sync.proto` in CI or manually? Where is generation configured?
- Which worker languages are currently considered canonical for benchmarks (the compose file comments indicate some are commented out)?
- Any cluster-specific assumptions (k3s, network addresses) that must be preserved when tweaking reconnect logic?

If anything above is unclear or you want more detail in a specific area (proto/codegen, worker onboarding, Docker/CI), tell me which section to expand. Thanks!
