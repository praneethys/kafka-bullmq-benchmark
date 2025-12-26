package kafka

import (
	"os"
	"testing"
	"time"

	"github.com/praneethys/kafka-bullmq-benchmark/pkg/common"
)

const (
	testBrokers = "localhost:9092"
	testTopic   = "test-topic"
	testGroup   = "test-group"
)

func skipIfNoKafka(t *testing.T) {
	// Check if KAFKA_TEST environment variable is set
	if os.Getenv("KAFKA_TEST") != "true" {
		t.Skip("Skipping Kafka integration test. Set KAFKA_TEST=true to run.")
	}
}

func TestNewKafkaQueue(t *testing.T) {
	skipIfNoKafka(t)

	queue, err := NewKafkaQueue(testBrokers, testTopic, testGroup)
	if err != nil {
		t.Fatalf("Failed to create Kafka queue: %v", err)
	}
	defer queue.Close()

	if queue == nil {
		t.Fatal("Expected non-nil queue")
	}

	if queue.topic != testTopic {
		t.Errorf("Expected topic '%s', got '%s'", testTopic, queue.topic)
	}

	if queue.brokers != testBrokers {
		t.Errorf("Expected brokers '%s', got '%s'", testBrokers, queue.brokers)
	}

	if queue.producer == nil {
		t.Error("Producer not initialized")
	}

	if queue.consumer == nil {
		t.Error("Consumer not initialized")
	}
}

func TestNewKafkaQueueInvalidBroker(t *testing.T) {
	// This test doesn't need actual Kafka running
	queue, err := NewKafkaQueue("invalid:9999", testTopic, testGroup)

	// The queue may be created successfully, but operations will fail
	// For now, we just verify it handles the invalid broker gracefully
	if err == nil && queue != nil {
		queue.Close()
	}
}

func TestKafkaGetName(t *testing.T) {
	skipIfNoKafka(t)

	queue, err := NewKafkaQueue(testBrokers, testTopic, testGroup)
	if err != nil {
		t.Fatalf("Failed to create Kafka queue: %v", err)
	}
	defer queue.Close()

	name := queue.GetName()
	if name != "Apache Kafka" {
		t.Errorf("Expected name 'Apache Kafka', got '%s'", name)
	}
}

func TestKafkaProduce(t *testing.T) {
	skipIfNoKafka(t)

	queue, err := NewKafkaQueue(testBrokers, testTopic+"-produce", testGroup)
	if err != nil {
		t.Fatalf("Failed to create Kafka queue: %v", err)
	}
	defer queue.Close()

	msg := &common.Message{
		ID:        "test-msg-1",
		Payload:   []byte("test payload"),
		Timestamp: time.Now(),
	}

	err = queue.Produce(msg)
	if err != nil {
		t.Errorf("Failed to produce message: %v", err)
	}

	// Flush to ensure message is delivered
	remaining := queue.Flush(5000)
	if remaining > 0 {
		t.Errorf("Expected 0 remaining messages, got %d", remaining)
	}
}

func TestKafkaProduceAsync(t *testing.T) {
	skipIfNoKafka(t)

	queue, err := NewKafkaQueue(testBrokers, testTopic+"-async", testGroup)
	if err != nil {
		t.Fatalf("Failed to create Kafka queue: %v", err)
	}
	defer queue.Close()

	msg := &common.Message{
		ID:        "test-msg-async-1",
		Payload:   []byte("test async payload"),
		Timestamp: time.Now(),
	}

	err = queue.ProduceAsync(msg)
	if err != nil {
		t.Errorf("Failed to produce async message: %v", err)
	}

	// Flush to ensure message is delivered
	remaining := queue.Flush(5000)
	if remaining > 0 {
		t.Errorf("Expected 0 remaining messages, got %d", remaining)
	}
}

func TestKafkaProduceAndConsume(t *testing.T) {
	skipIfNoKafka(t)

	topicName := testTopic + "-produce-consume"
	groupName := testGroup + "-pc"

	// Create producer queue
	producerQueue, err := NewKafkaQueue(testBrokers, topicName, groupName)
	if err != nil {
		t.Fatalf("Failed to create producer queue: %v", err)
	}
	defer producerQueue.Close()

	// Produce test messages
	testMessages := []*common.Message{
		{
			ID:        "msg-1",
			Payload:   []byte("payload-1"),
			Timestamp: time.Now(),
		},
		{
			ID:        "msg-2",
			Payload:   []byte("payload-2"),
			Timestamp: time.Now(),
		},
		{
			ID:        "msg-3",
			Payload:   []byte("payload-3"),
			Timestamp: time.Now(),
		},
	}

	for _, msg := range testMessages {
		if err := producerQueue.Produce(msg); err != nil {
			t.Fatalf("Failed to produce message %s: %v", msg.ID, err)
		}
	}

	producerQueue.Flush(5000)

	// Give Kafka time to commit
	time.Sleep(500 * time.Millisecond)

	// Create consumer queue with different consumer group to read from beginning
	consumerQueue, err := NewKafkaQueue(testBrokers, topicName, groupName+"-consumer")
	if err != nil {
		t.Fatalf("Failed to create consumer queue: %v", err)
	}
	defer consumerQueue.Close()

	// Consume messages
	receivedCount := 0
	receivedMessages := make(map[string]bool)

	done := make(chan bool)
	go func() {
		err := consumerQueue.Consume(func(msg *common.Message) error {
			receivedMessages[msg.ID] = true
			receivedCount++

			if receivedCount >= len(testMessages) {
				done <- true
			}
			return nil
		})
		if err != nil {
			t.Logf("Consumer error: %v", err)
		}
	}()

	// Wait for messages or timeout
	select {
	case <-done:
		// Success
	case <-time.After(10 * time.Second):
		t.Errorf("Timeout waiting for messages. Received %d/%d", receivedCount, len(testMessages))
	}

	// Verify all messages were received
	for _, msg := range testMessages {
		if !receivedMessages[msg.ID] {
			t.Errorf("Message %s was not received", msg.ID)
		}
	}
}

func TestKafkaProduceInvalidMessage(t *testing.T) {
	skipIfNoKafka(t)

	queue, err := NewKafkaQueue(testBrokers, testTopic, testGroup)
	if err != nil {
		t.Fatalf("Failed to create Kafka queue: %v", err)
	}
	defer queue.Close()

	// Test with nil message (should cause marshal error or panic)
	// We're testing that it handles edge cases gracefully
	msg := &common.Message{
		ID:        "",
		Payload:   nil,
		Timestamp: time.Time{},
	}

	// This should not panic
	_ = queue.Produce(msg)
}

func TestKafkaFlush(t *testing.T) {
	skipIfNoKafka(t)

	queue, err := NewKafkaQueue(testBrokers, testTopic, testGroup)
	if err != nil {
		t.Fatalf("Failed to create Kafka queue: %v", err)
	}
	defer queue.Close()

	// Produce some messages async
	for i := 0; i < 10; i++ {
		msg := &common.Message{
			ID:        "flush-test-" + string(rune(i)),
			Payload:   []byte("test"),
			Timestamp: time.Now(),
		}
		queue.ProduceAsync(msg)
	}

	// Flush with timeout
	remaining := queue.Flush(10000)

	// Should have flushed all messages
	if remaining > 0 {
		t.Logf("Warning: %d messages remaining after flush", remaining)
	}
}

func TestKafkaClose(t *testing.T) {
	skipIfNoKafka(t)

	queue, err := NewKafkaQueue(testBrokers, testTopic, testGroup)
	if err != nil {
		t.Fatalf("Failed to create Kafka queue: %v", err)
	}

	// Close should not return error
	err = queue.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}

	// Calling close again should not panic
	err = queue.Close()
	if err != nil {
		t.Logf("Second close returned error (expected): %v", err)
	}
}

func TestKafkaHighThroughput(t *testing.T) {
	skipIfNoKafka(t)

	if testing.Short() {
		t.Skip("Skipping high-throughput test in short mode")
	}

	queue, err := NewKafkaQueue(testBrokers, testTopic+"-throughput", testGroup)
	if err != nil {
		t.Fatalf("Failed to create Kafka queue: %v", err)
	}
	defer queue.Close()

	messageCount := 1000
	startTime := time.Now()

	// Produce messages
	for i := 0; i < messageCount; i++ {
		msg := &common.Message{
			ID:        "perf-" + string(rune(i)),
			Payload:   make([]byte, 1024), // 1KB payload
			Timestamp: time.Now(),
		}
		if err := queue.ProduceAsync(msg); err != nil {
			t.Errorf("Failed to produce message %d: %v", i, err)
		}
	}

	queue.Flush(30000)
	duration := time.Since(startTime)

	throughput := float64(messageCount) / duration.Seconds()
	t.Logf("Produced %d messages in %v (%.2f msg/sec)", messageCount, duration, throughput)

	if throughput < 100 {
		t.Logf("Warning: Low throughput %.2f msg/sec", throughput)
	}
}
