package metrics

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/kafka-bullmq-benchmark/pkg/common"
)

// Benchmark runs performance tests on message queues
type Benchmark struct {
	config    *common.BenchmarkConfig
	collector *Collector
}

// NewBenchmark creates a new benchmark instance
func NewBenchmark(config *common.BenchmarkConfig) *Benchmark {
	return &Benchmark{
		config:    config,
		collector: NewCollector(),
	}
}

// RunProducerBenchmark runs a producer-only benchmark
func (b *Benchmark) RunProducerBenchmark(queue common.MessageQueue) (*common.BenchmarkResult, error) {
	b.collector.Reset()

	payload := make([]byte, b.config.MessageSize)
	for i := range payload {
		payload[i] = byte(i % 256)
	}

	var wg sync.WaitGroup
	messagesPerProducer := b.config.MessageCount / b.config.ProducerCount

	startTime := time.Now()

	for p := 0; p < b.config.ProducerCount; p++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()

			for i := 0; i < messagesPerProducer; i++ {
				msg := &common.Message{
					ID:        uuid.New().String(),
					Payload:   payload,
					Timestamp: time.Now(),
				}

				msgStart := time.Now()

				// Use async produce for Kafka if available
				if kq, ok := queue.(interface{ ProduceAsync(*common.Message) error }); ok {
					if err := kq.ProduceAsync(msg); err != nil {
						b.collector.RecordError()
						continue
					}
				} else {
					if err := queue.Produce(msg); err != nil {
						b.collector.RecordError()
						continue
					}
				}

				latency := time.Since(msgStart)
				b.collector.RecordLatency(latency)
				b.collector.AddBytesProcessed(int64(len(payload)))
			}
		}(p)
	}

	wg.Wait()

	// Flush Kafka producer if available
	if kq, ok := queue.(interface{ Flush(int) int }); ok {
		kq.Flush(30000) // 30 second timeout
	}

	b.collector.Stop()

	duration := time.Since(startTime)
	fmt.Printf("Producer benchmark completed in %v\n", duration)

	return b.collector.GetResults(queue.GetName(), b.config.MessageCount), nil
}

// RunConsumerBenchmark runs a consumer-only benchmark
func (b *Benchmark) RunConsumerBenchmark(queue common.MessageQueue, expectedMessages int) (*common.BenchmarkResult, error) {
	b.collector.Reset()

	var wg sync.WaitGroup
	messagesChan := make(chan *common.Message, 1000)
	stopChan := make(chan bool)
	receivedCount := 0
	var countMu sync.Mutex

	// Start consumers
	for c := 0; c < b.config.ConsumerCount; c++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()

			handler := func(msg *common.Message) error {
				latency := time.Since(msg.Timestamp)
				b.collector.RecordLatency(latency)
				b.collector.AddBytesProcessed(int64(len(msg.Payload)))

				countMu.Lock()
				receivedCount++
				if receivedCount >= expectedMessages {
					select {
					case stopChan <- true:
					default:
					}
				}
				countMu.Unlock()

				return nil
			}

			queue.Consume(handler)
		}(c)
	}

	// Wait for all messages or timeout
	select {
	case <-stopChan:
		fmt.Printf("All %d messages consumed\n", expectedMessages)
	case <-time.After(time.Duration(b.config.DurationSeconds) * time.Second):
		fmt.Printf("Timeout reached, consumed %d messages\n", receivedCount)
	}

	b.collector.Stop()

	return b.collector.GetResults(queue.GetName(), receivedCount), nil
}

// RunFullBenchmark runs both producer and consumer benchmarks
func (b *Benchmark) RunFullBenchmark(producerQueue, consumerQueue common.MessageQueue) (*common.BenchmarkResult, error) {
	fmt.Printf("Starting full benchmark for %s\n", producerQueue.GetName())
	fmt.Printf("Configuration: %d messages, %d bytes each, %d producers, %d consumers\n",
		b.config.MessageCount, b.config.MessageSize, b.config.ProducerCount, b.config.ConsumerCount)

	b.collector.Reset()

	payload := make([]byte, b.config.MessageSize)
	for i := range payload {
		payload[i] = byte(i % 256)
	}

	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	messagesPerProducer := b.config.MessageCount / b.config.ProducerCount
	receivedCount := 0
	var countMu sync.Mutex
	stopChan := make(chan bool, 1)

	// Start consumers first
	for c := 0; c < b.config.ConsumerCount; c++ {
		consumerWg.Add(1)
		go func(consumerID int) {
			defer consumerWg.Done()

			handler := func(msg *common.Message) error {
				latency := time.Since(msg.Timestamp)
				b.collector.RecordLatency(latency)
				b.collector.AddBytesProcessed(int64(len(msg.Payload)))

				countMu.Lock()
				receivedCount++
				if receivedCount >= b.config.MessageCount {
					select {
					case stopChan <- true:
					default:
					}
				}
				countMu.Unlock()

				return nil
			}

			consumerQueue.Consume(handler)
		}(c)
	}

	// Give consumers time to start
	time.Sleep(2 * time.Second)

	// Start producers
	for p := 0; p < b.config.ProducerCount; p++ {
		producerWg.Add(1)
		go func(producerID int) {
			defer producerWg.Done()

			for i := 0; i < messagesPerProducer; i++ {
				msg := &common.Message{
					ID:        uuid.New().String(),
					Payload:   payload,
					Timestamp: time.Now(),
				}

				if kq, ok := producerQueue.(interface{ ProduceAsync(*common.Message) error }); ok {
					if err := kq.ProduceAsync(msg); err != nil {
						b.collector.RecordError()
					}
				} else {
					if err := producerQueue.Produce(msg); err != nil {
						b.collector.RecordError()
					}
				}
			}
		}(p)
	}

	producerWg.Wait()
	fmt.Println("All producers finished")

	// Flush Kafka producer if available
	if kq, ok := producerQueue.(interface{ Flush(int) int }); ok {
		kq.Flush(30000)
	}

	// Wait for all messages to be consumed or timeout
	select {
	case <-stopChan:
		fmt.Printf("All %d messages consumed\n", b.config.MessageCount)
	case <-time.After(time.Duration(b.config.DurationSeconds) * time.Second):
		fmt.Printf("Timeout reached, consumed %d/%d messages\n", receivedCount, b.config.MessageCount)
	}

	b.collector.Stop()

	return b.collector.GetResults(producerQueue.GetName(), b.config.MessageCount), nil
}
