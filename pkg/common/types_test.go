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

	if config.ProducerCount != 10 {
		t.Errorf("Expected ProducerCount 10, got %d", config.ProducerCount)
	}

	if config.ConsumerCount != 10 {
		t.Errorf("Expected ConsumerCount 10, got %d", config.ConsumerCount)
	}

	if config.BatchSize != 100 {
		t.Errorf("Expected BatchSize 100, got %d", config.BatchSize)
	}

	if config.DurationSeconds != 60 {
		t.Errorf("Expected DurationSeconds 60, got %d", config.DurationSeconds)
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

	// Verify all fields
	if result.QueueType != "Test Queue" {
		t.Errorf("Expected QueueType 'Test Queue', got '%s'", result.QueueType)
	}

	if result.MessageCount != 1000 {
		t.Errorf("Expected MessageCount 1000, got %d", result.MessageCount)
	}

	if result.Duration != 10*time.Second {
		t.Errorf("Expected Duration 10s, got %v", result.Duration)
	}

	if result.Throughput != 100.0 {
		t.Errorf("Expected Throughput 100.0, got %.2f", result.Throughput)
	}

	if result.AvgLatency != 5*time.Millisecond {
		t.Errorf("Expected AvgLatency 5ms, got %v", result.AvgLatency)
	}

	if result.P50Latency != 4*time.Millisecond {
		t.Errorf("Expected P50Latency 4ms, got %v", result.P50Latency)
	}

	if result.P95Latency != 8*time.Millisecond {
		t.Errorf("Expected P95Latency 8ms, got %v", result.P95Latency)
	}

	if result.P99Latency != 10*time.Millisecond {
		t.Errorf("Expected P99Latency 10ms, got %v", result.P99Latency)
	}

	if result.MinLatency != 1*time.Millisecond {
		t.Errorf("Expected MinLatency 1ms, got %v", result.MinLatency)
	}

	if result.MaxLatency != 15*time.Millisecond {
		t.Errorf("Expected MaxLatency 15ms, got %v", result.MaxLatency)
	}

	if result.ErrorCount != 0 {
		t.Errorf("Expected ErrorCount 0, got %d", result.ErrorCount)
	}

	if result.SuccessCount != 1000 {
		t.Errorf("Expected SuccessCount 1000, got %d", result.SuccessCount)
	}

	if result.BytesProcessed != 1024000 {
		t.Errorf("Expected BytesProcessed 1024000, got %d", result.BytesProcessed)
	}

	if result.MBPerSecond != 100.0 {
		t.Errorf("Expected MBPerSecond 100.0, got %.2f", result.MBPerSecond)
	}
}
