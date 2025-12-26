package metrics

import (
	"testing"
	"time"

	"github.com/praneethys/kafka-bullmq-benchmark/pkg/common"
)

func TestNewBenchmark(t *testing.T) {
	config := &common.BenchmarkConfig{
		MessageCount:    1000,
		MessageSize:     1024,
		ProducerCount:   10,
		ConsumerCount:   10,
		DurationSeconds: 60,
	}

	benchmark := NewBenchmark(config)

	if benchmark == nil {
		t.Fatal("Expected non-nil benchmark")
	}

	if benchmark.config != config {
		t.Error("Config not set correctly")
	}

	if benchmark.collector == nil {
		t.Error("Collector not initialized")
	}
}

// MockQueue implements the MessageQueue interface for testing
type MockQueue struct {
	name         string
	produceCount int
	consumeCount int
	produceError error
	consumeError error
	shouldFail   bool
}

func (m *MockQueue) Produce(msg *common.Message) error {
	m.produceCount++
	if m.produceError != nil {
		return m.produceError
	}
	return nil
}

func (m *MockQueue) Consume(handler func(*common.Message) error) error {
	if m.consumeError != nil {
		return m.consumeError
	}

	// Simulate consuming messages
	for i := 0; i < 10; i++ {
		msg := &common.Message{
			ID:        "test",
			Payload:   []byte("test"),
			Timestamp: time.Now().Add(-time.Millisecond),
		}
		if err := handler(msg); err != nil {
			return err
		}
		m.consumeCount++
	}

	// Block forever (will be interrupted by test timeout or context)
	select {}
}

func (m *MockQueue) Close() error {
	return nil
}

func (m *MockQueue) GetName() string {
	return m.name
}

func TestBenchmarkConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *common.BenchmarkConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &common.BenchmarkConfig{
				MessageCount:    1000,
				MessageSize:     1024,
				ProducerCount:   10,
				ConsumerCount:   10,
				DurationSeconds: 60,
			},
			wantErr: false,
		},
		{
			name: "zero message count",
			config: &common.BenchmarkConfig{
				MessageCount:    0,
				MessageSize:     1024,
				ProducerCount:   10,
				ConsumerCount:   10,
				DurationSeconds: 60,
			},
			wantErr: false, // We allow this, it's up to the benchmark logic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			benchmark := NewBenchmark(tt.config)
			if benchmark == nil && !tt.wantErr {
				t.Error("Expected non-nil benchmark")
			}
		})
	}
}

func TestCollectorIntegration(t *testing.T) {
	collector := NewCollector()

	// Simulate benchmark data
	for i := 0; i < 100; i++ {
		latency := time.Duration(i+1) * time.Millisecond
		collector.RecordLatency(latency)
		collector.AddBytesProcessed(1024)
	}

	collector.RecordError()
	collector.RecordError()

	collector.Stop()

	result := collector.GetResults("Integration Test", 100)

	if result.SuccessCount != 100 {
		t.Errorf("Expected 100 successes, got %d", result.SuccessCount)
	}

	if result.ErrorCount != 2 {
		t.Errorf("Expected 2 errors, got %d", result.ErrorCount)
	}

	if result.BytesProcessed != 102400 {
		t.Errorf("Expected 102400 bytes, got %d", result.BytesProcessed)
	}

	if result.Throughput <= 0 {
		t.Error("Expected positive throughput")
	}
}
