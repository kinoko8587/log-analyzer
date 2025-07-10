# Log Analyzer

A Go-based high-performance log analysis system for processing and analyzing large volumes of log data with concurrency and race condition handling.

---

## 🗂 Project Structure

```
log-analyzer/
├── cmd/
│ └── log-analyzer/ # Main application entry point
├── internal/
│ ├── api/ # HTTP API server (Gin)
│ ├── config/ # Configuration logic
│ ├── generator/ # Concurrent log generator
│ ├── processor/ # Log processing logic
│ ├── storage/ # Thread-safe statistics storage
│ └── metrics/ # (Planned) Prometheus integration
├── pkg/
│ └── analyzer/ # Public log model (e.g. Log struct)
├── config/
│ └── config.yaml # (Planned) Config file for rate, keywords
```

---

## 🚀 Getting Started

### ✅ Prerequisites

- Go 1.21 or higher (recommended: 1.24.2)

### 📦 Installation

🛠 Build Binary
bash
```
go build -o bin/log-analyzer cmd/log-analyzer/main.go
./bin/log-analyzer
```
📡 HTTP API Endpoints
Endpoint	Description
/stats/errors	Get total count of error logs
/stats/all	Get all log stats (info/warn/error)
/health	Health check

Test with:

bash
```
curl http://localhost:8080/stats/errors
curl http://localhost:8080/stats/all
```

## ✅ Implemented Components
Log structure with timestamp, level, and message (pkg/analyzer/log.go)

Log generator with multiple goroutines and channels (internal/generator)

Log processor for real-time error counting (internal/processor)

Thread-safe statistics storage using sync.Mutex (internal/storage)

HTTP API with Gin framework (internal/api/server.go)

Graceful shutdown on OS signals

## 📍 Roadmap

### 🧩 Phase 1: MVP

- Define log structure

- Implement generator (3 workers, 100 logs/sec)

- Implement processor (2 workers)

- Thread-safe storage (mutex)

- HTTP API: /stats/errors, /health

- Graceful shutdown

### 🚦 Phase 2: Concurrency Experiments

- Add alternative storage: sync.Map

- Channel-based actor pattern

- Benchmark and compare performance

- Integrate net/http/pprof

### 📊 Phase 3: Advanced Stats

- Keyword counting (e.g. timeout)

- Time-windowed stats (last 10s)

- New API: /stats/keyword?q=timeout, /stats/errors?window=10s

☁️ Phase 4: Observability and Config
- Add Prometheus metrics

- Persist log stats (file or SQLite)

- Use context.Context for goroutine management

- External config file for log rate and keywords

### 🧪 Development Tools

Purpose	Tool
Race detection	go run -race
Benchmarking	testing.B
Profiling	net/http/pprof
Logging	log / zerolog / zap
Metrics	prometheus/client_golang

### 🧹 Graceful Shutdown

Uses os/signal to handle SIGINT/SIGTERM

Waits for all goroutines to complete using sync.WaitGroup

