# Setup Guide

This guide will help you set up and run the Kafka vs BullMQ benchmark.

## Prerequisites

### Required Software

1. **Go 1.21 or higher**
   ```bash
   # Check Go version
   go version

   # Install Go (macOS)
   brew install go

   # Install Go (Ubuntu/Debian)
   sudo apt-get update
   sudo apt-get install golang-1.21
   ```

2. **Docker and Docker Compose**
   ```bash
   # Check Docker version
   docker --version
   docker-compose --version

   # Install Docker Desktop (macOS/Windows)
   # Download from: https://www.docker.com/products/docker-desktop

   # Install Docker (Ubuntu/Debian)
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   ```

3. **Make (optional but recommended)**
   ```bash
   # Usually pre-installed on macOS/Linux
   make --version
   ```

### System Requirements

- **CPU**: 4+ cores recommended
- **RAM**: 8GB minimum, 16GB recommended
- **Disk**: 10GB free space
- **Network**: Internet connection for downloading dependencies

## Installation Steps

### 1. Clone or Download the Repository

If you haven't already:
```bash
git clone https://github.com/praneethys/kafka-bullmq-benchmark.git
cd kafka-bullmq-benchmark
```

### 2. Install Go Dependencies

```bash
# Using Make
make deps

# Or manually
go mod download
go mod verify
```

### 3. Start Infrastructure

```bash
# Using Make
make docker-up

# Or manually
docker-compose up -d

# Wait for services to be healthy (30-60 seconds)
docker-compose ps
```

You should see all services running:
```
NAME                COMMAND                  SERVICE             STATUS              PORTS
kafka               "kafka-server-start.…"   kafka               running (healthy)   0.0.0.0:9092->9092/tcp, 9093/tcp
kafka-ui            "java -jar kafka-ui.…"   kafka-ui            running             0.0.0.0:8080->8080/tcp
redis               "docker-entrypoint.s…"   redis               running (healthy)   0.0.0.0:6379->6379/tcp
redis-commander     "docker-entrypoint.s…"   redis-commander     running             0.0.0.0:8081->8081/tcp
```

**Note**: Kafka 4.1 uses KRaft mode (Kafka Raft) and does not require Zookeeper.

### 4. Verify Services

**Check Kafka:**
```bash
docker-compose exec kafka kafka-topics --bootstrap-server localhost:9092 --list
```

**Check Redis:**
```bash
docker-compose exec redis redis-cli ping
# Should return: PONG
```

**Access Web UIs:**
- Kafka UI: http://localhost:8080
- Redis Commander: http://localhost:8081

### 5. Build the Benchmark

```bash
# Using Make
make build

# Or manually
go build -o build/benchmark cmd/benchmark/main.go
```

## Running Benchmarks

### Quick Test (100K messages)

```bash
make run-quick
```

### Full Benchmark (1M messages)

```bash
make run-full
```

### Custom Benchmark

```bash
./build/benchmark \
  -messages 500000 \
  -size 2048 \
  -producers 30 \
  -consumers 30 \
  -queue both
```

### Benchmark Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `-messages` | 100000 | Number of messages to send |
| `-size` | 1024 | Message size in bytes |
| `-producers` | 10 | Number of producer goroutines |
| `-consumers` | 10 | Number of consumer goroutines |
| `-duration` | 300 | Maximum duration in seconds |
| `-queue` | both | Queue type: kafka, redis, or both |
| `-kafka-brokers` | localhost:9092 | Kafka broker addresses |
| `-kafka-topic` | benchmark-topic | Kafka topic name |
| `-redis-addr` | localhost:6379 | Redis server address |
| `-redis-stream` | benchmark-stream | Redis stream key |
| `-output` | ./results | Output directory for results |

## Viewing Results

### Console Output

Results are printed to the console after each benchmark run.

### Exported Files

Results are saved in the `results/` directory:
```bash
ls -l results/
# benchmark-results-20240115-143022.json
# benchmark-results-20240115-143022.csv
```

### View Latest Results

```bash
make results-view
```

## Troubleshooting

### Issue: "Cannot connect to Kafka"

**Solution:**
```bash
# Restart Kafka
make docker-restart

# Check Kafka logs
docker-compose logs kafka

# Verify Kafka is listening
docker-compose exec kafka nc -zv localhost 9092
```

### Issue: "Cannot connect to Redis"

**Solution:**
```bash
# Restart Redis
docker-compose restart redis

# Check Redis logs
docker-compose logs redis

# Test Redis connection
docker-compose exec redis redis-cli ping
```

### Issue: "Out of memory"

**Solution:**
1. Increase Docker memory limit (Docker Desktop → Settings → Resources)
2. Reduce message count or size:
   ```bash
   ./build/benchmark -messages 50000 -size 512
   ```

### Issue: "Too many open files"

**Solution (macOS/Linux):**
```bash
# Increase file descriptor limit
ulimit -n 10000

# Or add to ~/.bash_profile or ~/.zshrc:
echo "ulimit -n 10000" >> ~/.bash_profile
```

### Issue: Build fails with missing dependencies

**Solution:**
```bash
# Clean and reinstall
make clean
go clean -modcache
make deps
make build
```

## Advanced Configuration

### Kafka Tuning

Edit [docker-compose.yml](docker-compose.yml) to adjust Kafka settings:
```yaml
KAFKA_NUM_PARTITIONS: 10          # Increase for more parallelism
KAFKA_LOG_RETENTION_HOURS: 1      # Adjust retention
KAFKA_NUM_NETWORK_THREADS: 8      # Match CPU cores
```

### Redis Tuning

Edit [docker-compose.yml](docker-compose.yml) to adjust Redis settings:
```yaml
command: >
  redis-server
  --maxmemory 4gb                  # Increase memory
  --maxmemory-policy allkeys-lru
```

### Go Application Tuning

Set environment variables:
```bash
# Increase max goroutines
export GOMAXPROCS=8

# Run benchmark
./build/benchmark -producers 100 -consumers 100
```

## Running in Production

For production-like testing:

1. **Use separate machines** for producers, consumers, and brokers
2. **Enable replication** in Kafka (minimum 3 brokers)
3. **Enable persistence** in Redis (AOF or RDB)
4. **Use realistic network conditions** (add latency simulation)
5. **Monitor resources** (CPU, memory, disk, network)

Example production Kafka config:
```yaml
KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 3
KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 3
KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 2
```

## Cleanup

### Stop Services

```bash
make docker-down
```

### Remove All Data

```bash
make docker-clean
```

### Clean Build Artifacts

```bash
make clean
```

## Next Steps

1. Run the benchmarks with different configurations
2. Analyze the results in the `results/` directory
3. Update the white paper with your findings
4. Share your results with the community

## Getting Help

- Check the [README.md](README.md) for general information
- Review [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines
- Open an issue on GitHub for bugs or questions

## Useful Commands

```bash
# Development workflow
make dev              # Start dev environment
make build           # Build the benchmark
make test            # Run tests
make lint            # Run linter
make fmt             # Format code

# Benchmark shortcuts
make run-quick       # Quick test
make run-full        # Full benchmark
make run-kafka       # Kafka only
make run-redis       # Redis only

# Docker management
make docker-up       # Start services
make docker-down     # Stop services
make docker-restart  # Restart services
make docker-logs     # View logs

# Results management
make results-view    # View latest results
make results-clean   # Clean results
```
