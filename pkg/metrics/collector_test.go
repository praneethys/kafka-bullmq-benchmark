package metrics

import (
	"testing"
	"time"
)

func TestNewCollector(t *testing.T) {
	collector := NewCollector()

	if collector == nil {
		t.Fatal("Expected non-nil collector")
	}

	if len(collector.latencies) != 0 {
		t.Errorf("Expected empty latencies, got %d", len(collector.latencies))
	}

	if collector.errorCount != 0 {
		t.Errorf("Expected errorCount 0, got %d", collector.errorCount)
	}

	if collector.successCount != 0 {
		t.Errorf("Expected successCount 0, got %d", collector.successCount)
	}
}

func TestRecordLatency(t *testing.T) {
	collector := NewCollector()

	latency := 5 * time.Millisecond
	collector.RecordLatency(latency)

	if collector.successCount != 1 {
		t.Errorf("Expected successCount 1, got %d", collector.successCount)
	}

	if len(collector.latencies) != 1 {
		t.Errorf("Expected 1 latency recorded, got %d", len(collector.latencies))
	}

	if collector.latencies[0] != latency {
		t.Errorf("Expected latency %v, got %v", latency, collector.latencies[0])
	}
}

func TestRecordError(t *testing.T) {
	collector := NewCollector()

	collector.RecordError()
	collector.RecordError()

	if collector.errorCount != 2 {
		t.Errorf("Expected errorCount 2, got %d", collector.errorCount)
	}
}

func TestAddBytesProcessed(t *testing.T) {
	collector := NewCollector()

	collector.AddBytesProcessed(1024)
	collector.AddBytesProcessed(2048)

	if collector.bytesProcessed != 3072 {
		t.Errorf("Expected bytesProcessed 3072, got %d", collector.bytesProcessed)
	}
}

func TestGetResults(t *testing.T) {
	collector := NewCollector()

	// Record some metrics
	collector.RecordLatency(1 * time.Millisecond)
	collector.RecordLatency(2 * time.Millisecond)
	collector.RecordLatency(3 * time.Millisecond)
	collector.RecordLatency(10 * time.Millisecond)
	collector.RecordError()
	collector.AddBytesProcessed(4096)

	collector.Stop()

	result := collector.GetResults("Test Queue", 4)

	if result.QueueType != "Test Queue" {
		t.Errorf("Expected QueueType 'Test Queue', got '%s'", result.QueueType)
	}

	if result.MessageCount != 4 {
		t.Errorf("Expected MessageCount 4, got %d", result.MessageCount)
	}

	if result.SuccessCount != 4 {
		t.Errorf("Expected SuccessCount 4, got %d", result.SuccessCount)
	}

	if result.ErrorCount != 1 {
		t.Errorf("Expected ErrorCount 1, got %d", result.ErrorCount)
	}

	if result.BytesProcessed != 4096 {
		t.Errorf("Expected BytesProcessed 4096, got %d", result.BytesProcessed)
	}

	if result.MinLatency != 1*time.Millisecond {
		t.Errorf("Expected MinLatency 1ms, got %v", result.MinLatency)
	}

	if result.MaxLatency != 10*time.Millisecond {
		t.Errorf("Expected MaxLatency 10ms, got %v", result.MaxLatency)
	}
}

func TestReset(t *testing.T) {
	collector := NewCollector()

	// Add some data
	collector.RecordLatency(5 * time.Millisecond)
	collector.RecordError()
	collector.AddBytesProcessed(1024)

	// Reset
	collector.Reset()

	if collector.errorCount != 0 {
		t.Errorf("Expected errorCount 0 after reset, got %d", collector.errorCount)
	}

	if collector.successCount != 0 {
		t.Errorf("Expected successCount 0 after reset, got %d", collector.successCount)
	}

	if collector.bytesProcessed != 0 {
		t.Errorf("Expected bytesProcessed 0 after reset, got %d", collector.bytesProcessed)
	}

	if len(collector.latencies) != 0 {
		t.Errorf("Expected empty latencies after reset, got %d", len(collector.latencies))
	}
}

func TestLatencyPercentiles(t *testing.T) {
	collector := NewCollector()

	// Add 100 latencies from 1ms to 100ms
	for i := 1; i <= 100; i++ {
		collector.RecordLatency(time.Duration(i) * time.Millisecond)
	}

	collector.Stop()
	result := collector.GetResults("Test", 100)

	// P50 should be around 50ms
	if result.P50Latency < 49*time.Millisecond || result.P50Latency > 51*time.Millisecond {
		t.Errorf("Expected P50 around 50ms, got %v", result.P50Latency)
	}

	// P95 should be around 95ms
	if result.P95Latency < 94*time.Millisecond || result.P95Latency > 96*time.Millisecond {
		t.Errorf("Expected P95 around 95ms, got %v", result.P95Latency)
	}

	// P99 should be around 99ms
	if result.P99Latency < 98*time.Millisecond || result.P99Latency > 100*time.Millisecond {
		t.Errorf("Expected P99 around 99ms, got %v", result.P99Latency)
	}
}
