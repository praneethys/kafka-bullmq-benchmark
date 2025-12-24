package metrics

import (
	"sort"
	"sync"
	"time"

	"github.com/yourusername/kafka-bullmq-benchmark/pkg/common"
)

// Collector collects benchmark metrics
type Collector struct {
	mu              sync.Mutex
	latencies       []time.Duration
	errorCount      int
	successCount    int
	bytesProcessed  int64
	startTime       time.Time
	endTime         time.Time
}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	return &Collector{
		latencies:  make([]time.Duration, 0, 1000000),
		startTime:  time.Now(),
	}
}

// RecordLatency records a message latency
func (c *Collector) RecordLatency(latency time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.latencies = append(c.latencies, latency)
	c.successCount++
}

// RecordError records an error
func (c *Collector) RecordError() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errorCount++
}

// AddBytesProcessed adds to the bytes processed counter
func (c *Collector) AddBytesProcessed(bytes int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bytesProcessed += bytes
}

// Stop marks the end of the benchmark
func (c *Collector) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.endTime = time.Now()
}

// GetResults calculates and returns benchmark results
func (c *Collector) GetResults(queueType string, messageCount int) *common.BenchmarkResult {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.endTime.IsZero() {
		c.endTime = time.Now()
	}

	duration := c.endTime.Sub(c.startTime)

	result := &common.BenchmarkResult{
		QueueType:      queueType,
		MessageCount:   messageCount,
		Duration:       duration,
		ErrorCount:     c.errorCount,
		SuccessCount:   c.successCount,
		BytesProcessed: c.bytesProcessed,
	}

	if duration.Seconds() > 0 {
		result.Throughput = float64(c.successCount) / duration.Seconds()
		result.MBPerSecond = float64(c.bytesProcessed) / (1024 * 1024) / duration.Seconds()
	}

	if len(c.latencies) > 0 {
		sort.Slice(c.latencies, func(i, j int) bool {
			return c.latencies[i] < c.latencies[j]
		})

		result.MinLatency = c.latencies[0]
		result.MaxLatency = c.latencies[len(c.latencies)-1]
		result.P50Latency = c.latencies[len(c.latencies)*50/100]
		result.P95Latency = c.latencies[len(c.latencies)*95/100]
		result.P99Latency = c.latencies[len(c.latencies)*99/100]

		var sum time.Duration
		for _, l := range c.latencies {
			sum += l
		}
		result.AvgLatency = sum / time.Duration(len(c.latencies))
	}

	return result
}

// Reset resets all metrics
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.latencies = make([]time.Duration, 0, 1000000)
	c.errorCount = 0
	c.successCount = 0
	c.bytesProcessed = 0
	c.startTime = time.Now()
	c.endTime = time.Time{}
}
