# Quick Start Guide

Get up and running with the Kafka vs BullMQ benchmark in 5 minutes!

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- 8GB RAM

## 1. Start Infrastructure (2 minutes)

```bash
docker-compose up -d
```

Wait for services to be healthy (check with `docker-compose ps`).

## 2. Build & Run (3 minutes)

```bash
# Install dependencies
go mod download

# Build
go build -o benchmark cmd/benchmark/main.go

# Run quick test (100K messages)
./benchmark -messages 100000 -producers 10 -consumers 10

# Or run full test (1M messages)
./benchmark -messages 1000000 -producers 50 -consumers 50
```

## 3. View Results

Results are displayed in the console and saved to `./results/` directory.

## Using Make (Recommended)

```bash
# Start infrastructure
make docker-up

# Run quick benchmark
make run-quick

# Run full benchmark
make run-full

# View results
make results-view
```

## Common Commands

```bash
# Kafka only
./benchmark -queue kafka -messages 500000

# Redis only
./benchmark -queue redis -messages 500000

# Large messages (10KB)
./benchmark -messages 100000 -size 10240

# High concurrency
./benchmark -producers 100 -consumers 100
```

## Monitoring

- **Kafka UI**: http://localhost:8080
- **Redis Commander**: http://localhost:8081

## Troubleshooting

```bash
# Restart everything
make docker-restart

# View logs
docker-compose logs -f

# Clean and rebuild
make clean && make build
```

## Next Steps

See [SETUP.md](SETUP.md) for detailed configuration options.
See [README.md](README.md) for complete documentation.
