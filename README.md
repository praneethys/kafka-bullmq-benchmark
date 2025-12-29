# Kafka vs BullMQ (Redis Streams) Performance Benchmark

A comprehensive distributed systems performance comparison between Apache Kafka and BullMQ (Redis Streams) implementation in Go, designed to handle millions of transactions per second.

## Overview

This project provides a robust benchmarking framework to compare the performance characteristics of two popular message queue systems:

- **Apache Kafka 4.1**: A distributed event streaming platform
- **BullMQ (Redis Streams)**: A Redis-based queue system

The benchmark measures:
- **Throughput**: Messages per second
- **Latency**: P50, P95, P99 percentiles
- **Bandwidth**: MB/s processed
- **Reliability**: Success/error rates
- **Scalability**: Performance under varying loads

## Features

- High-performance concurrent producers and consumers
- Configurable message sizes and volumes
- Real-time metrics collection
- Comprehensive latency analysis (min, avg, P50, P95, P99, max)
- Results export in JSON and CSV formats
- Docker Compose setup for easy infrastructure deployment
- Designed for million+ messages per second throughput
- Professional white paper template included

## Architecture

```
.
├── cmd/
│   └── benchmark/          # Main benchmark application
├── pkg/
│   ├── common/            # Shared types and interfaces
│   ├── kafka/             # Kafka implementation
│   ├── redis/             # Redis Streams (BullMQ) implementation
│   └── metrics/           # Benchmarking and metrics collection
├── results/               # Benchmark results output
├── whitepaper/           # Research white paper
├── docker-compose.yml    # Infrastructure setup
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- At least 8GB RAM for optimal performance
- 4+ CPU cores recommended

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/praneethys/kafka-bullmq-benchmark.git
cd kafka-bullmq-benchmark
```

### 2. Start Infrastructure

```bash
# Start Kafka (KRaft mode) and Redis
docker compose up -d

# Wait for services to be healthy (30-60 seconds)
docker compose ps
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Benchmark

```bash
# Run both Kafka and Redis benchmarks
go run cmd/benchmark/main.go

# Or build and run
go build -o benchmark cmd/benchmark/main.go
./benchmark
```

## Configuration Options

```bash
./benchmark [options]

Options:
  -messages int       Number of messages to send (default: 100000)
  -size int          Size of each message in bytes (default: 1024)
  -producers int     Number of producer goroutines (default: 10)
  -consumers int     Number of consumer goroutines (default: 10)
  -duration int      Maximum duration in seconds (default: 300)
  -queue string      Queue type: kafka, redis, or both (default: "both")
  -kafka-brokers     Kafka broker addresses (default: "localhost:9092")
  -kafka-topic       Kafka topic name (default: "benchmark-topic")
  -redis-addr        Redis server address (default: "localhost:6379")
  -redis-stream      Redis stream key (default: "benchmark-stream")
  -output string     Output directory for results (default: "./results")
```

## Example Usage

### High-Throughput Test (1 Million Messages)

```bash
./benchmark \
  -messages 1000000 \
  -size 1024 \
  -producers 50 \
  -consumers 50 \
  -queue both
```

### Large Message Test

```bash
./benchmark \
  -messages 100000 \
  -size 10240 \
  -producers 20 \
  -consumers 20 \
  -queue both
```

### Kafka Only Test

```bash
./benchmark \
  -messages 500000 \
  -size 2048 \
  -producers 30 \
  -consumers 30 \
  -queue kafka
```

### Redis Only Test

```bash
./benchmark \
  -messages 500000 \
  -size 2048 \
  -producers 30 \
  -consumers 30 \
  -queue redis
```

## Monitoring

### Kafka UI
Access Kafka UI at: http://localhost:8080
- View topics, partitions, and messages
- Monitor consumer groups
- Check broker health

### Redis Commander
Access Redis Commander at: http://localhost:8081
- View Redis streams
- Monitor memory usage
- Inspect data structures

## Results Analysis

Results are automatically saved to the `./results` directory in both JSON and CSV formats:

```
results/
├── benchmark-results-20240115-143022.json
└── benchmark-results-20240115-143022.csv
```

### Sample Output

```
================================================================================
Benchmark Results: Apache Kafka
================================================================================
Messages:           100000
Duration:           15.234s
Throughput:         6563.21 msg/s
Bandwidth:          6.41 MB/s
Success Count:      100000
Error Count:        0
Bytes Processed:    102400000 (97.66 MB)

Latency Statistics:
  Min:              0.12 ms
  Avg:              1.52 ms
  P50:              1.23 ms
  P95:              3.45 ms
  P99:              5.67 ms
  Max:              12.34 ms
================================================================================
```

## Performance Tuning

### Kafka Tuning

The Kafka implementation includes optimizations for high throughput:
- Batch size: 1MB
- Compression: LZ4
- Linger time: 10ms
- In-flight requests: 5
- 10 partitions for parallelism

### Redis Tuning

Redis is configured for maximum performance:
- No persistence (AOF/RDB disabled)
- 2GB max memory with LRU eviction
- Connection pooling (100 connections)
- Batch reads (10 messages per read)

## Testing

The project includes comprehensive test coverage for all components.

### Running Unit Tests

```bash
# Run all unit tests
make test

# Run tests for specific package
go test ./pkg/common/...
go test ./pkg/metrics/...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...
```

### Running Integration Tests

Integration tests require running Kafka and Redis services. Start the infrastructure first:

```bash
# Start services
docker compose up -d

# Wait for services to be healthy
sleep 30

# Run Kafka integration tests
KAFKA_TEST=true go test -v ./pkg/kafka/...

# Run Redis integration tests
REDIS_TEST=true go test -v ./pkg/redis/...

# Run all integration tests
KAFKA_TEST=true REDIS_TEST=true go test -v ./pkg/...
```

### Test Coverage

The test suite includes:
- **Unit Tests**: Core types, metrics collection, export functionality, main application
- **Integration Tests**: Kafka and Redis queue implementations, end-to-end benchmarks
- **Benchmark Tests**: Performance and throughput testing
- **Mock Implementations**: For isolated testing
- **Edge Case Tests**: Error handling, invalid inputs, resource cleanup

Total test coverage: 20 unit tests, 32 integration tests (52 tests total)

## Troubleshooting

### Kafka Connection Issues

```bash
# Check Kafka health
docker compose exec kafka kafka-broker-api-versions --bootstrap-server localhost:9092

# View Kafka logs
docker compose logs kafka
```

### Redis Connection Issues

```bash
# Check Redis health
docker compose exec redis redis-cli ping

# View Redis logs
docker compose logs redis
```

### Resource Constraints

If you encounter memory issues:
1. Reduce message count or size
2. Reduce producer/consumer counts
3. Increase Docker memory limits in Docker Desktop settings

## White Paper

A comprehensive white paper analyzing the results is available in the `whitepaper/` directory. It includes:
- Methodology
- Detailed results analysis
- Performance comparisons
- Use case recommendations
- Scalability analysis

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - see LICENSE file for details

## Citation

If you use this benchmark in your research, please cite:

```
Kafka vs BullMQ Performance Benchmark
https://github.com/praneethys/kafka-bullmq-benchmark
2024
```

## Acknowledgments

- Apache Kafka community
- Redis and BullMQ communities
- Go community

## Contact

For questions or issues, please open a GitHub issue or contact the maintainers.

## Roadmap

- [ ] Add Prometheus metrics export
- [ ] Implement custom partitioning strategies
- [ ] Add transaction support benchmarks
- [ ] Multi-node cluster benchmarks
- [ ] Network latency simulation
- [ ] Failure recovery testing
- [ ] Cost analysis comparisons
