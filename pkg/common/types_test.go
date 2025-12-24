package common

import (
	"testing"
	"time"
)

func TestMessage(t *testing.T) {
	msg := &Message{
		ID:        "test-123",
		Payload:   []byte("test payload"),
		Timestamp: time.Now(),
	}

	if msg.ID != "test-123" {
		t.Errorf("Expected ID 'test-123', got '%s'", msg.ID)
	}

	if string(msg.Payload) != "test payload" {
		t.Errorf("Expected payload 'test payload', got '%s'", string(msg.Payload))
	}

	if msg.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
}

func TestBenchmarkConfig(t *testing.T) {
	config := &BenchmarkConfig{
		MessageCount:    1000,
		MessageSize:     1024,
		ProducerCount:   10,
		ConsumerCount:   10,
		BatchSize:       100,
		DurationSeconds: 60,
	}

	if config.MessageCount != 1000 {
		t.Errorf("Expected MessageCount 1000, got %d", config.MessageCount)
	}

	if config.MessageSize != 1024 {
		t.Errorf("Expected MessageSize 1024, got %d", config.MessageSize)
	}
}

func TestBenchmarkResult(t *testing.T) {
	result := &BenchmarkResult{
		QueueType:      "Test Queue",
		MessageCount:   1000,
		Duration:       10 * time.Second,
		Throughput:     100.0,
		AvgLatency:     5 * time.Millisecond,
		P50Latency:     4 * time.Millisecond,
		P95Latency:     8 * time.Millisecond,
		P99Latency:     10 * time.Millisecond,
		MaxLatency:     15 * time.Millisecond,
		MinLatency:     1 * time.Millisecond,
		ErrorCount:     0,
		SuccessCount:   1000,
		BytesProcessed: 1024000,
		MBPerSecond:    100.0,
	}

	if result.QueueType != "Test Queue" {
		t.Errorf("Expected QueueType 'Test Queue', got '%s'", result.QueueType)
	}

	if result.SuccessCount != 1000 {
		t.Errorf("Expected SuccessCount 1000, got %d", result.SuccessCount)
	}

	if result.Throughput != 100.0 {
		t.Errorf("Expected Throughput 100.0, got %.2f", result.Throughput)
	}
}
