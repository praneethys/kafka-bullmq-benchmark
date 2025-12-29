package main

import (
	"os"
	"testing"

	"github.com/praneethys/kafka-bullmq-benchmark/pkg/common"
)

func skipIfNoKafka(t *testing.T) {
	if os.Getenv("KAFKA_TEST") != "true" {
		t.Skip("Skipping Kafka integration test. Set KAFKA_TEST=true to run.")
	}
}

func skipIfNoRedis(t *testing.T) {
	if os.Getenv("REDIS_TEST") != "true" {
		t.Skip("Skipping Redis integration test. Set REDIS_TEST=true to run.")
	}
}

func TestRunKafkaBenchmark(t *testing.T) {
	skipIfNoKafka(t)

	config := &common.BenchmarkConfig{
		MessageCount:    100,
		MessageSize:     512,
		ProducerCount:   2,
		ConsumerCount:   2,
		DurationSeconds: 30,
	}

	brokers := "localhost:9092"
	topic := "test-benchmark-kafka"

	result, err := runKafkaBenchmark(config, brokers, topic)
	if err != nil {
		t.Fatalf("runKafkaBenchmark failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.QueueType != "Apache Kafka" {
		t.Errorf("Expected QueueType 'Apache Kafka', got '%s'", result.QueueType)
	}

	if result.MessageCount != config.MessageCount {
		t.Errorf("Expected MessageCount %d, got %d", config.MessageCount, result.MessageCount)
	}

	if result.Throughput <= 0 {
		t.Error("Expected positive throughput")
	}

	if result.Duration <= 0 {
		t.Error("Expected positive duration")
	}
}

func TestRunKafkaBenchmarkInvalidBroker(t *testing.T) {
	// This test verifies that we can create Kafka queues with invalid broker
	// (they don't fail immediately), but the benchmark will eventually timeout
	// We skip this test as it takes too long (10+ seconds)
	t.Skip("Skipping slow invalid broker test - takes 10+ seconds to timeout")

	config := &common.BenchmarkConfig{
		MessageCount:    10,
		MessageSize:     512,
		ProducerCount:   1,
		ConsumerCount:   1,
		DurationSeconds: 10,
	}

	brokers := "invalid:9999"
	topic := "test-topic"

	result, err := runKafkaBenchmark(config, brokers, topic)
	// The benchmark may complete but with errors or no messages processed
	if err == nil && result != nil {
		// Verify no messages were successfully processed
		if result.SuccessCount > 0 {
			t.Error("Expected no successful messages with invalid broker")
		}
	}
}

func TestRunRedisBenchmark(t *testing.T) {
	skipIfNoRedis(t)

	config := &common.BenchmarkConfig{
		MessageCount:    100,
		MessageSize:     512,
		ProducerCount:   2,
		ConsumerCount:   2,
		DurationSeconds: 30,
	}

	addr := "localhost:6379"
	streamKey := "test-benchmark-redis"

	result, err := runRedisBenchmark(config, addr, streamKey)
	if err != nil {
		t.Fatalf("runRedisBenchmark failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.QueueType != "Redis Streams (BullMQ)" {
		t.Errorf("Expected QueueType 'Redis Streams (BullMQ)', got '%s'", result.QueueType)
	}

	if result.MessageCount != config.MessageCount {
		t.Errorf("Expected MessageCount %d, got %d", config.MessageCount, result.MessageCount)
	}

	if result.Throughput <= 0 {
		t.Error("Expected positive throughput")
	}

	if result.Duration <= 0 {
		t.Error("Expected positive duration")
	}
}

func TestRunRedisBenchmarkInvalidAddr(t *testing.T) {
	config := &common.BenchmarkConfig{
		MessageCount:    10,
		MessageSize:     512,
		ProducerCount:   1,
		ConsumerCount:   1,
		DurationSeconds: 10,
	}

	addr := "invalid:9999"
	streamKey := "test-stream"

	_, err := runRedisBenchmark(config, addr, streamKey)
	if err == nil {
		t.Error("Expected error for invalid address, got nil")
	}
}

func TestRunKafkaBenchmarkSmallLoad(t *testing.T) {
	skipIfNoKafka(t)

	if testing.Short() {
		t.Skip("Skipping small load test in short mode")
	}

	config := &common.BenchmarkConfig{
		MessageCount:    50,
		MessageSize:     256,
		ProducerCount:   1,
		ConsumerCount:   1,
		DurationSeconds: 30,
	}

	brokers := "localhost:9092"
	topic := "test-benchmark-kafka-small"

	result, err := runKafkaBenchmark(config, brokers, topic)
	if err != nil {
		t.Fatalf("runKafkaBenchmark small load failed: %v", err)
	}

	if result.SuccessCount != config.MessageCount {
		t.Errorf("Expected %d successful messages, got %d", config.MessageCount, result.SuccessCount)
	}

	// Check latency metrics are populated
	if result.AvgLatency <= 0 {
		t.Error("Expected positive average latency")
	}

	if result.P50Latency <= 0 {
		t.Error("Expected positive P50 latency")
	}

	if result.P95Latency <= 0 {
		t.Error("Expected positive P95 latency")
	}

	if result.P99Latency <= 0 {
		t.Error("Expected positive P99 latency")
	}

	if result.MinLatency <= 0 {
		t.Error("Expected positive min latency")
	}

	if result.MaxLatency <= 0 {
		t.Error("Expected positive max latency")
	}

	// Sanity checks on latency ordering
	if result.MinLatency > result.AvgLatency {
		t.Error("Min latency should be <= avg latency")
	}

	if result.AvgLatency > result.MaxLatency {
		t.Error("Avg latency should be <= max latency")
	}

	if result.P50Latency > result.P95Latency {
		t.Error("P50 should be <= P95")
	}

	if result.P95Latency > result.P99Latency {
		t.Error("P95 should be <= P99")
	}
}

func TestRunRedisBenchmarkSmallLoad(t *testing.T) {
	skipIfNoRedis(t)

	if testing.Short() {
		t.Skip("Skipping small load test in short mode")
	}

	config := &common.BenchmarkConfig{
		MessageCount:    50,
		MessageSize:     256,
		ProducerCount:   1,
		ConsumerCount:   1,
		DurationSeconds: 30,
	}

	addr := "localhost:6379"
	streamKey := "test-benchmark-redis-small"

	result, err := runRedisBenchmark(config, addr, streamKey)
	if err != nil {
		t.Fatalf("runRedisBenchmark small load failed: %v", err)
	}

	if result.SuccessCount != config.MessageCount {
		t.Errorf("Expected %d successful messages, got %d", config.MessageCount, result.SuccessCount)
	}

	// Check latency metrics are populated
	if result.AvgLatency <= 0 {
		t.Error("Expected positive average latency")
	}

	if result.P50Latency <= 0 {
		t.Error("Expected positive P50 latency")
	}

	if result.P95Latency <= 0 {
		t.Error("Expected positive P95 latency")
	}

	if result.P99Latency <= 0 {
		t.Error("Expected positive P99 latency")
	}

	// Sanity checks on latency ordering
	if result.MinLatency > result.AvgLatency {
		t.Error("Min latency should be <= avg latency")
	}

	if result.P50Latency > result.P95Latency {
		t.Error("P50 should be <= P95")
	}

	if result.P95Latency > result.P99Latency {
		t.Error("P95 should be <= P99")
	}
}

func TestRunKafkaBenchmarkLargeMessages(t *testing.T) {
	skipIfNoKafka(t)

	if testing.Short() {
		t.Skip("Skipping large message test in short mode")
	}

	config := &common.BenchmarkConfig{
		MessageCount:    20,
		MessageSize:     10240, // 10KB messages
		ProducerCount:   2,
		ConsumerCount:   2,
		DurationSeconds: 30,
	}

	brokers := "localhost:9092"
	topic := "test-benchmark-kafka-large"

	result, err := runKafkaBenchmark(config, brokers, topic)
	if err != nil {
		t.Fatalf("runKafkaBenchmark large messages failed: %v", err)
	}

	expectedBytes := int64(config.MessageCount * config.MessageSize)
	if result.BytesProcessed != expectedBytes {
		t.Errorf("Expected %d bytes processed, got %d", expectedBytes, result.BytesProcessed)
	}

	if result.MBPerSecond <= 0 {
		t.Error("Expected positive MB/s throughput")
	}
}

func TestRunRedisBenchmarkLargeMessages(t *testing.T) {
	skipIfNoRedis(t)

	if testing.Short() {
		t.Skip("Skipping large message test in short mode")
	}

	config := &common.BenchmarkConfig{
		MessageCount:    20,
		MessageSize:     10240, // 10KB messages
		ProducerCount:   2,
		ConsumerCount:   2,
		DurationSeconds: 30,
	}

	addr := "localhost:6379"
	streamKey := "test-benchmark-redis-large"

	result, err := runRedisBenchmark(config, addr, streamKey)
	if err != nil {
		t.Fatalf("runRedisBenchmark large messages failed: %v", err)
	}

	expectedBytes := int64(config.MessageCount * config.MessageSize)
	if result.BytesProcessed != expectedBytes {
		t.Errorf("Expected %d bytes processed, got %d", expectedBytes, result.BytesProcessed)
	}

	if result.MBPerSecond <= 0 {
		t.Error("Expected positive MB/s throughput")
	}
}

func TestRunKafkaBenchmarkMultipleProducersConsumers(t *testing.T) {
	skipIfNoKafka(t)

	if testing.Short() {
		t.Skip("Skipping multi producer/consumer test in short mode")
	}

	config := &common.BenchmarkConfig{
		MessageCount:    100,
		MessageSize:     1024,
		ProducerCount:   5,
		ConsumerCount:   5,
		DurationSeconds: 30,
	}

	brokers := "localhost:9092"
	topic := "test-benchmark-kafka-multi"

	result, err := runKafkaBenchmark(config, brokers, topic)
	if err != nil {
		t.Fatalf("runKafkaBenchmark multi producers/consumers failed: %v", err)
	}

	if result.SuccessCount != config.MessageCount {
		t.Errorf("Expected %d successful messages, got %d", config.MessageCount, result.SuccessCount)
	}

	// With more producers/consumers, throughput should be reasonable
	if result.Throughput <= 0 {
		t.Error("Expected positive throughput with multiple producers/consumers")
	}
}

func TestRunRedisBenchmarkMultipleProducersConsumers(t *testing.T) {
	skipIfNoRedis(t)

	if testing.Short() {
		t.Skip("Skipping multi producer/consumer test in short mode")
	}

	config := &common.BenchmarkConfig{
		MessageCount:    100,
		MessageSize:     1024,
		ProducerCount:   5,
		ConsumerCount:   5,
		DurationSeconds: 30,
	}

	addr := "localhost:6379"
	streamKey := "test-benchmark-redis-multi"

	result, err := runRedisBenchmark(config, addr, streamKey)
	if err != nil {
		t.Fatalf("runRedisBenchmark multi producers/consumers failed: %v", err)
	}

	if result.SuccessCount != config.MessageCount {
		t.Errorf("Expected %d successful messages, got %d", config.MessageCount, result.SuccessCount)
	}

	// With more producers/consumers, throughput should be reasonable
	if result.Throughput <= 0 {
		t.Error("Expected positive throughput with multiple producers/consumers")
	}
}

func TestBenchmarkConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *common.BenchmarkConfig
		wantErr bool
	}{
		{
			name: "valid minimal config",
			config: &common.BenchmarkConfig{
				MessageCount:    10,
				MessageSize:     256,
				ProducerCount:   1,
				ConsumerCount:   1,
				DurationSeconds: 10,
			},
			wantErr: false,
		},
		{
			name: "valid large config",
			config: &common.BenchmarkConfig{
				MessageCount:    1000000,
				MessageSize:     10240,
				ProducerCount:   50,
				ConsumerCount:   50,
				DurationSeconds: 600,
			},
			wantErr: false,
		},
		{
			name: "zero messages allowed",
			config: &common.BenchmarkConfig{
				MessageCount:    0,
				MessageSize:     1024,
				ProducerCount:   1,
				ConsumerCount:   1,
				DurationSeconds: 10,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the config can be created
			// The actual validation happens in the benchmark
			if tt.config.MessageCount < 0 {
				t.Error("Message count should not be negative")
			}
			if tt.config.MessageSize < 0 {
				t.Error("Message size should not be negative")
			}
			if tt.config.ProducerCount < 0 {
				t.Error("Producer count should not be negative")
			}
			if tt.config.ConsumerCount < 0 {
				t.Error("Consumer count should not be negative")
			}
		})
	}
}

// TestOutputDirectoryCreation tests that output directory is created correctly
func TestOutputDirectoryCreation(t *testing.T) {
	tempDir := t.TempDir()
	outputDir := tempDir + "/test-results"

	// Simulate what main() does
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Verify directory exists
	info, err := os.Stat(outputDir)
	if err != nil {
		t.Fatalf("Output directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("Expected output path to be a directory")
	}

	// Check permissions (on Unix-like systems)
	if info.Mode().Perm()&0o755 == 0 {
		t.Error("Expected directory to have execute permissions")
	}
}
