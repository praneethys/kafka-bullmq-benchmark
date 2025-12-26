# Testing Guide

This document provides comprehensive information about the test suite for the Kafka vs BullMQ benchmark project.

## Test Structure

The project includes comprehensive test coverage across all packages:

```
cmd/
└── benchmark/
    └── main_test.go           # Main application tests (12 tests)
pkg/
├── common/
│   └── types_test.go          # Core type tests (3 tests)
├── kafka/
│   └── kafka_test.go          # Kafka implementation tests (10 tests)
├── redis/
│   └── redis_test.go          # Redis implementation tests (12 tests)
└── metrics/
    ├── benchmark_test.go      # Benchmark framework tests (3 tests)
    ├── collector_test.go      # Metrics collector tests (7 tests)
    └── exporter_test.go       # Export functionality tests (5 tests)
```

**Total Test Count**: 52 tests
- **Unit Tests**: 20 tests (always run)
- **Integration Tests**: 32 tests (require running services)

## Running Tests

### Quick Test (Unit Tests Only)

```bash
# Run all unit tests
make test

# Or directly with go test
go test ./...
```

This will run all unit tests and skip integration tests. Integration tests are automatically skipped when Kafka/Redis are not available.

### Full Test Suite (With Integration Tests)

Integration tests require running Kafka and Redis services:

```bash
# 1. Start services
docker compose up -d

# 2. Wait for services to be healthy (important!)
sleep 30

# 3. Run all tests including integration tests
KAFKA_TEST=true REDIS_TEST=true go test -v ./...
```

### Running Specific Test Suites

#### Common Package Tests
```bash
go test -v ./pkg/common/...
```

#### Metrics Package Tests
```bash
go test -v ./pkg/metrics/...
```

#### Kafka Integration Tests
```bash
# Requires Kafka running
docker compose up -d kafka
sleep 20
KAFKA_TEST=true go test -v ./pkg/kafka/...
```

#### Redis Integration Tests
```bash
# Requires Redis running
docker compose up -d redis
sleep 5
REDIS_TEST=true go test -v ./pkg/redis/...
```

### Test Options

#### Verbose Output
```bash
go test -v ./...
```

#### With Coverage
```bash
go test -cover ./...

# Detailed coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Race Detection
```bash
go test -race ./...
```

#### Short Mode (Skip Long-Running Tests)
```bash
go test -short ./...
```

## Test Categories

### 1. Unit Tests

**Always run, no external dependencies required.**

#### pkg/common/types_test.go
- `TestMessage` - Validates Message struct
- `TestBenchmarkConfig` - Validates BenchmarkConfig struct
- `TestBenchmarkResult` - Validates BenchmarkResult struct

#### pkg/metrics/collector_test.go
- `TestNewCollector` - Collector initialization
- `TestRecordLatency` - Latency recording
- `TestRecordError` - Error counting
- `TestAddBytesProcessed` - Byte counting
- `TestGetResults` - Results calculation
- `TestReset` - Collector reset
- `TestLatencyPercentiles` - Percentile calculations

#### pkg/metrics/exporter_test.go
- `TestExportToJSON` - JSON export functionality
- `TestExportToCSV` - CSV export functionality
- `TestGenerateReport` - Combined report generation
- `TestPrintResults` - Console output formatting
- `TestCompareResults` - Comparison table generation

#### pkg/metrics/benchmark_test.go
- `TestNewBenchmark` - Benchmark initialization
- `TestBenchmarkConfigValidation` - Config validation
- `TestCollectorIntegration` - Collector integration

#### cmd/benchmark/main_test.go
- `TestRunRedisBenchmarkInvalidAddr` - Error handling for invalid Redis address
- `TestBenchmarkConfigValidation` - Configuration validation
- `TestOutputDirectoryCreation` - Output directory creation

**Note**: `TestRunKafkaBenchmarkInvalidBroker` is skipped by default as it takes 10+ seconds to timeout with invalid brokers.

### 2. Integration Tests

**Require running services, skip automatically if not available.**

#### pkg/kafka/kafka_test.go

Set `KAFKA_TEST=true` to run these tests.

- `TestNewKafkaQueue` - Queue initialization
- `TestNewKafkaQueueInvalidBroker` - Error handling for invalid broker
- `TestKafkaGetName` - GetName() method
- `TestKafkaProduce` - Synchronous message production
- `TestKafkaProduceAsync` - Asynchronous message production
- `TestKafkaProduceAndConsume` - End-to-end produce/consume flow
- `TestKafkaProduceInvalidMessage` - Edge case handling
- `TestKafkaFlush` - Producer flush functionality
- `TestKafkaClose` - Cleanup and resource release
- `TestKafkaHighThroughput` - Performance under load (1000 messages)

#### pkg/redis/redis_test.go

Set `REDIS_TEST=true` to run these tests.

- `TestNewRedisQueue` - Queue initialization
- `TestNewRedisQueueInvalidAddr` - Error handling for invalid address
- `TestRedisGetName` - GetName() method
- `TestRedisProduce` - Message production
- `TestRedisProduceAsync` - Async production (same as Produce for Redis)
- `TestRedisProduceAndConsume` - End-to-end produce/consume flow
- `TestRedisProduceInvalidMessage` - Edge case handling
- `TestRedisGetStreamInfo` - Stream information retrieval
- `TestRedisTrimStream` - Stream trimming functionality
- `TestRedisClose` - Cleanup and resource release
- `TestRedisHighThroughput` - Performance under load (1000 messages)
- `TestRedisConsumerGroupHandling` - Consumer group management

#### cmd/benchmark/main_test.go

Set `KAFKA_TEST=true` and/or `REDIS_TEST=true` to run these tests.

- `TestRunKafkaBenchmark` - Full Kafka benchmark execution
- `TestRunRedisBenchmark` - Full Redis benchmark execution
- `TestRunKafkaBenchmarkSmallLoad` - Kafka with small message load (50 messages)
- `TestRunRedisBenchmarkSmallLoad` - Redis with small message load (50 messages)
- `TestRunKafkaBenchmarkLargeMessages` - Kafka with 10KB messages
- `TestRunRedisBenchmarkLargeMessages` - Redis with 10KB messages
- `TestRunKafkaBenchmarkMultipleProducersConsumers` - Kafka with 5 producers/consumers
- `TestRunRedisBenchmarkMultipleProducersConsumers` - Redis with 5 producers/consumers

## Test Design Principles

### 1. Smart Skipping
Integration tests automatically skip when services aren't available:

```go
func skipIfNoKafka(t *testing.T) {
    if os.Getenv("KAFKA_TEST") != "true" {
        t.Skip("Skipping Kafka integration test. Set KAFKA_TEST=true to run.")
    }
}
```

### 2. Isolated Testing
Each test creates its own resources with unique names:

```go
topicName := testTopic + "-produce-consume"
groupName := testGroup + "-pc"
```

### 3. Cleanup
Tests properly clean up resources:

```go
defer queue.Close()
defer queue.client.Del(queue.ctx, streamKey)
```

### 4. Timeout Protection
Tests include timeouts to prevent hanging:

```go
select {
case <-done:
    // Success
case <-time.After(10 * time.Second):
    t.Errorf("Timeout waiting for messages")
}
```

### 5. Mock Implementations
Unit tests use mocks to avoid external dependencies:

```go
type MockQueue struct {
    produceCount int
    consumeCount int
    produceError error
}
```

## CI/CD Integration

The GitHub Actions workflow automatically runs tests:

```yaml
- name: Run tests
  run: make test
```

Integration tests are skipped in CI since external services aren't configured. For full CI testing, you would need to:

1. Add Kafka and Redis as services to the GitHub Actions workflow
2. Set environment variables `KAFKA_TEST=true` and `REDIS_TEST=true`
3. Add service readiness checks

## Coverage Goals

- **Unit Test Coverage**: 80%+ for core logic
- **Integration Test Coverage**: All public APIs tested
- **Edge Cases**: Invalid inputs, connection errors, timeouts
- **Performance Tests**: High-throughput scenarios

## Writing New Tests

### Adding a Unit Test

```go
func TestMyFeature(t *testing.T) {
    // Arrange
    input := setupTestData()

    // Act
    result := myFunction(input)

    // Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Adding an Integration Test

```go
func TestMyIntegrationFeature(t *testing.T) {
    skipIfNoKafka(t)  // or skipIfNoRedis(t)

    // Setup
    queue, err := NewKafkaQueue(...)
    if err != nil {
        t.Fatalf("Setup failed: %v", err)
    }
    defer queue.Close()

    // Test logic
    // ...
}
```

## Troubleshooting Tests

### Integration Tests Always Skip

Make sure you:
1. Started the services: `docker-compose up -d`
2. Waited for health: `sleep 30`
3. Set environment variables: `KAFKA_TEST=true REDIS_TEST=true`

### Tests Timeout

Increase timeout or check service health:

```bash
# Check Kafka
docker-compose exec kafka kafka-broker-api-versions --bootstrap-server localhost:9092

# Check Redis
docker-compose exec redis redis-cli ping
```

### Race Detector Warnings

Run with race detector to find concurrency issues:

```bash
go test -race ./...
```

### Flaky Tests

Integration tests may be flaky if:
- Services aren't fully initialized (increase wait time)
- System is under heavy load (reduce message count)
- Timeouts are too short (increase timeout values)

## Best Practices

1. **Always run unit tests before committing**: `make test`
2. **Run integration tests before pushing**: Full test suite
3. **Check coverage periodically**: `go test -cover ./...`
4. **Use `-v` flag for debugging**: See detailed test output
5. **Keep tests fast**: Use `-short` flag to skip long tests during development
6. **Clean up resources**: Always use `defer` for cleanup
7. **Test edge cases**: Invalid inputs, errors, timeouts
8. **Use table-driven tests**: For testing multiple scenarios

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Framework](https://github.com/stretchr/testify) (if you want to add it)
- [Go Test Coverage](https://go.dev/blog/cover)
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
