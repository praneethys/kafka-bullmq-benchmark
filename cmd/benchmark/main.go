package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/praneethys/kafka-bullmq-benchmark/pkg/common"
	"github.com/praneethys/kafka-bullmq-benchmark/pkg/kafka"
	"github.com/praneethys/kafka-bullmq-benchmark/pkg/metrics"
	"github.com/praneethys/kafka-bullmq-benchmark/pkg/redis"
)

func main() {
	// Command line flags
	messageCount := flag.Int("messages", 100000, "Number of messages to send")
	messageSize := flag.Int("size", 1024, "Size of each message in bytes")
	producers := flag.Int("producers", 10, "Number of producer goroutines")
	consumers := flag.Int("consumers", 10, "Number of consumer goroutines")
	duration := flag.Int("duration", 300, "Maximum duration in seconds")
	queueType := flag.String("queue", "both", "Queue type to test: kafka, redis, or both")
	kafkaBrokers := flag.String("kafka-brokers", "localhost:9092", "Kafka broker addresses")
	kafkaTopic := flag.String("kafka-topic", "benchmark-topic", "Kafka topic name")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis server address")
	redisStream := flag.String("redis-stream", "benchmark-stream", "Redis stream key")
	outputDir := flag.String("output", "./results", "Output directory for results")

	flag.Parse()

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	config := &common.BenchmarkConfig{
		MessageCount:    *messageCount,
		MessageSize:     *messageSize,
		ProducerCount:   *producers,
		ConsumerCount:   *consumers,
		DurationSeconds: *duration,
	}

	fmt.Println("Kafka vs BullMQ (Redis Streams) Benchmark")
	fmt.Println("==========================================")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Messages:       %d\n", config.MessageCount)
	fmt.Printf("  Message Size:   %d bytes\n", config.MessageSize)
	fmt.Printf("  Producers:      %d\n", config.ProducerCount)
	fmt.Printf("  Consumers:      %d\n", config.ConsumerCount)
	fmt.Printf("  Max Duration:   %d seconds\n", config.DurationSeconds)
	fmt.Println()

	var results []*common.BenchmarkResult

	// Run Kafka benchmark
	if *queueType == "kafka" || *queueType == "both" {
		fmt.Println("Starting Kafka benchmark...")
		kafkaResult, err := runKafkaBenchmark(config, *kafkaBrokers, *kafkaTopic)
		if err != nil {
			log.Printf("Kafka benchmark failed: %v", err)
		} else {
			results = append(results, kafkaResult)
			metrics.PrintResults(kafkaResult)
		}
	}

	// Run Redis benchmark
	if *queueType == "redis" || *queueType == "both" {
		fmt.Println("Starting Redis (BullMQ) benchmark...")
		redisResult, err := runRedisBenchmark(config, *redisAddr, *redisStream)
		if err != nil {
			log.Printf("Redis benchmark failed: %v", err)
		} else {
			results = append(results, redisResult)
			metrics.PrintResults(redisResult)
		}
	}

	// Compare results if both were run
	if len(results) > 1 {
		metrics.CompareResults(results)
	}

	// Generate report
	if err := metrics.GenerateReport(results, *outputDir); err != nil {
		log.Printf("Failed to generate report: %v", err)
	}

	fmt.Println("Benchmark completed successfully!")
}

func runKafkaBenchmark(config *common.BenchmarkConfig, brokers, topic string) (*common.BenchmarkResult, error) {
	// Create Kafka producer queue
	producerQueue, err := kafka.NewKafkaQueue(brokers, topic, "benchmark-producer-group")
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	defer producerQueue.Close()

	// Create Kafka consumer queue
	consumerQueue, err := kafka.NewKafkaQueue(brokers, topic, "benchmark-consumer-group")
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}
	defer consumerQueue.Close()

	// Run benchmark
	benchmark := metrics.NewBenchmark(config)
	result, err := benchmark.RunFullBenchmark(producerQueue, consumerQueue)
	if err != nil {
		return nil, fmt.Errorf("benchmark failed: %w", err)
	}

	return result, nil
}

func runRedisBenchmark(config *common.BenchmarkConfig, addr, streamKey string) (*common.BenchmarkResult, error) {
	// Create Redis producer queue
	producerQueue, err := redis.NewRedisQueue(addr, streamKey, "benchmark-group", "producer")
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis producer: %w", err)
	}
	defer producerQueue.Close()

	// Create Redis consumer queue
	consumerQueue, err := redis.NewRedisQueue(addr, streamKey, "benchmark-group", "consumer")
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis consumer: %w", err)
	}
	defer consumerQueue.Close()

	// Run benchmark
	benchmark := metrics.NewBenchmark(config)
	result, err := benchmark.RunFullBenchmark(producerQueue, consumerQueue)
	if err != nil {
		return nil, fmt.Errorf("benchmark failed: %w", err)
	}

	return result, nil
}
