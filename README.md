# 📊 Log Analyzer 

A high-concurrency log ingestion and analytics system built in Go.  
This project simulates real-world backend challenges like message queue ingestion, batch aggregation, caching, and system observability — all built with **Clean Architecture** and **modular services**.


## ✨ Features

- Modular service design (API, Ingestor, Batch Job)
- Kafka-based log ingestion (async + scalable)
- PostgreSQL log storage & denormalized summary tables
- RESTful APIs with Gin + Redis caching
- Batch job for daily aggregation & keyword stats
- Prometheus metrics / pprof profiling / graceful shutdown
- CI/CD with GitHub Actions & containerized environment

---

## 🗂 Project Structure

```bash
log-analyzer/
├── cmd/
│   ├── api-server/        # REST API server (dashboard & stats)
│   ├── log-ingestor/      # Kafka consumer to persist logs to DB
│   └── batch-job/         # Daily aggregation job
├── internal/
│   ├── domain/            # Entity models & interface definitions
│   ├── usecase/           # Application logic (analyze, query, etc.)
│   ├── infrastructure/    # DB, Redis, Kafka adapters
│   ├── interface/         
│   │   ├── http/          # REST handlers
│   │   └── grpc/          # (Optional) gRPC handlers
│   └── scheduler/         # cron-like batch executor
├── pkg/                   # Shared utilities and log models
├── migrations/            # DB migration SQL
├── config/                # YAML/env configuration
├── deploy/                # Docker Compose, CI/CD files
└── README.md
```

📦 Services Overview
| Service	|Description|
|----|-----|
|api-server|	Exposes REST APIs for log stats/dashboard|
|log-ingestor|	Listens to Kafka and writes logs to DB|
|batch-job	|Aggregates daily user stats and keywords|

## 🧩 Phase Roadmap
### ✅ Phase 1: MVP Setup - Kafka Ingestion
Focus: replacing in-memory storage with Kafka → PostgreSQL

- [x] Set up Kafka + Zookeeper via Docker Compose

- [x] Create Kafka Producer (mock log generator)

- [ ] Create Kafka Consumer service (log-ingestor)

- [ ] Design PostgreSQL schema log_raw (normalized)

- [ ] Implement LogRepository interface

- [ ] Abstract Kafka via QueueSubscriber interface

- [ ] Add graceful shutdown & error handling for workers

### 🚀 Phase 2: API Server + Redis Cache
Focus: API endpoints + caching for dashboard queries

- [ ] Create REST API server using Gin (api-server)

- [ ] Define routes:

    - /stats/errors?window=10s

    - /stats/keywords?q=timeout

- [ ] Add Redis cache layer for recent queries

- [ ] Implement fallback: Redis miss → DB → re-cache

- [ ] Add JWT-based auth (optional)

### 📊 Phase 3: Batch Job for Aggregation
Focus: summary tables and offline aggregation logic

- [ ] Create batch-job binary

- [ ] Cron schedule aggregation every 1 day

- [ ] Create table user_log_summary (denormalized)

- [ ] Aggregate: total logs, error rate, keyword frequency

- [ ] Insert into summary table

- [ ] Track & log invalid/malformed rows

### ☁️ Phase 4: Observability & Metrics
Focus: system introspection & monitoring

- [ ] Add Prometheus metrics (/metrics)

    - log ingestion rate

    - error count per worker

    - batch job duration

- [ ] Enable pprof (/debug/pprof)

- [ ] Add structured logging (zap or zerolog)

- [ ] Log collector ready (stdout or file)

### 🔁 Phase 5: CI/CD + Deploy
Focus: DevOps and reproducible environments

- [ ] Create docker-compose.yml with:
    - Kafka + Zookeeper
    - PostgreSQL
    - Redis
    - All 3 services

- [ ] Add .env and YAML config loaders

- [ ] Create GitHub Actions workflow:
    - test
    - lint
    - docker build
- [ ] Add Makefile for build/dev tasks

🛠 Tools Used
|Area	    |Tech|
|-----------|-----|
|Language	|Go 1.21+|
|Queue	    |Kafka (segmentio/kafka-go)|
|DB	        |PostgreSQL / SQLite|
|API        |	Gin|
|Cache      |	Redis (go-redis)|
|Metrics	|Prometheus client_golang|
|Profiling	|net/http/pprof|
|Testing	|Go test, stretchr/testify|
|Deployment	|Docker, GitHub Actions|

📌 How to Run (MVP Phase)
```
bash
# Start Kafka + PostgreSQL + services
docker-compose up -d

# Run log producer to push logs into Kafka
go run cmd/log-generator/main.go

# Run Kafka consumer to persist to DB
go run cmd/log-ingestor/main.go

# Start API server
go run cmd/api-server/main.go

# Check stats
curl http://localhost:8080/stats/errors
```
## 🧠 Key Concepts Practiced
- Clean Architecture in Go
- Message queue (Kafka) consumer pipelines
- High-concurrency goroutine orchestration
- Time-windowed stat computation
- Caching & fallback strategies (Redis)
- System observability with Prometheus
- Modular service decomposition

## 🧪 Future Ideas
- gRPC endpoints with protobuf
- SQLite fallback mode (embedded DB)
- Live dashboard with WebSocket streaming
- ElasticSearch log search integration

### 🧑‍💻 Author
Ann Chen