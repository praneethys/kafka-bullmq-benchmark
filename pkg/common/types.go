package common

import "time"

// Message represents a benchmark message
type Message struct {
	ID        string    `json:"id"`
	Payload   []byte    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

// MessageQueue interface for both Kafka and Redis implementations
type MessageQueue interface {
	Produce(msg *Message) error
	Consume(handler func(*Message) error) error
	Close() error
	GetName() string
}

// BenchmarkConfig holds configuration for benchmarks
type BenchmarkConfig struct {
	MessageCount    int
	MessageSize     int
	ProducerCount   int
	ConsumerCount   int
	BatchSize       int
	DurationSeconds int
}

// BenchmarkResult holds the results of a benchmark run
type BenchmarkResult struct {
	QueueType          string
	MessageCount       int
	Duration           time.Duration
	Throughput         float64 // messages per second
	AvgLatency         time.Duration
	P50Latency         time.Duration
	P95Latency         time.Duration
	P99Latency         time.Duration
	MaxLatency         time.Duration
	MinLatency         time.Duration
	ErrorCount         int
	SuccessCount       int
	BytesProcessed     int64
	MBPerSecond        float64
}
