package redis

import (
	"os"
	"testing"
	"time"

	"github.com/praneethys/kafka-bullmq-benchmark/pkg/common"
)

const (
	testAddr          = "localhost:6379"
	testStream        = "test-stream"
	testConsumerGroup = "test-consumer-group"
	testConsumerName  = "test-consumer"
)

func skipIfNoRedis(t *testing.T) {
	// Check if REDIS_TEST environment variable is set
	if os.Getenv("REDIS_TEST") != "true" {
		t.Skip("Skipping Redis integration test. Set REDIS_TEST=true to run.")
	}
}

func TestNewRedisQueue(t *testing.T) {
	skipIfNoRedis(t)

	queue, err := NewRedisQueue(testAddr, testStream, testConsumerGroup, testConsumerName)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()

	if queue == nil {
		t.Fatal("Expected non-nil queue")
	}

	if queue.streamKey != testStream {
		t.Errorf("Expected stream key '%s', got '%s'", testStream, queue.streamKey)
	}

	if queue.consumerGroup != testConsumerGroup {
		t.Errorf("Expected consumer group '%s', got '%s'", testConsumerGroup, queue.consumerGroup)
	}

	if queue.consumerName != testConsumerName {
		t.Errorf("Expected consumer name '%s', got '%s'", testConsumerName, queue.consumerName)
	}

	if queue.client == nil {
		t.Error("Redis client not initialized")
	}
}

func TestNewRedisQueueInvalidAddr(t *testing.T) {
	_, err := NewRedisQueue("invalid:9999", testStream, testConsumerGroup, testConsumerName)
	if err == nil {
		t.Error("Expected error for invalid Redis address, got nil")
	}
}

func TestRedisGetName(t *testing.T) {
	skipIfNoRedis(t)

	queue, err := NewRedisQueue(testAddr, testStream+"-name", testConsumerGroup, testConsumerName)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()

	name := queue.GetName()
	if name != "Redis Streams (BullMQ)" {
		t.Errorf("Expected name 'Redis Streams (BullMQ)', got '%s'", name)
	}
}

func TestRedisProduce(t *testing.T) {
	skipIfNoRedis(t)

	streamKey := testStream + "-produce"
	queue, err := NewRedisQueue(testAddr, streamKey, testConsumerGroup, testConsumerName)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()
	defer queue.client.Del(queue.ctx, streamKey) // Cleanup

	msg := &common.Message{
		ID:        "test-msg-1",
		Payload:   []byte("test payload"),
		Timestamp: time.Now(),
	}

	err = queue.Produce(msg)
	if err != nil {
		t.Errorf("Failed to produce message: %v", err)
	}

	// Verify message was added to stream
	info, err := queue.GetStreamInfo()
	if err != nil {
		t.Fatalf("Failed to get stream info: %v", err)
	}

	if info.Length < 1 {
		t.Errorf("Expected at least 1 message in stream, got %d", info.Length)
	}
}

func TestRedisProduceAsync(t *testing.T) {
	skipIfNoRedis(t)

	streamKey := testStream + "-async"
	queue, err := NewRedisQueue(testAddr, streamKey, testConsumerGroup, testConsumerName)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()
	defer queue.client.Del(queue.ctx, streamKey)

	msg := &common.Message{
		ID:        "test-msg-async-1",
		Payload:   []byte("test async payload"),
		Timestamp: time.Now(),
	}

	err = queue.ProduceAsync(msg)
	if err != nil {
		t.Errorf("Failed to produce async message: %v", err)
	}

	// Verify message was added
	info, err := queue.GetStreamInfo()
	if err != nil {
		t.Fatalf("Failed to get stream info: %v", err)
	}

	if info.Length < 1 {
		t.Errorf("Expected at least 1 message in stream, got %d", info.Length)
	}
}

func TestRedisProduceAndConsume(t *testing.T) {
	skipIfNoRedis(t)

	streamKey := testStream + "-produce-consume"
	groupName := testConsumerGroup + "-pc"

	// Create producer queue
	producerQueue, err := NewRedisQueue(testAddr, streamKey, groupName, testConsumerName+"-producer")
	if err != nil {
		t.Fatalf("Failed to create producer queue: %v", err)
	}
	defer producerQueue.Close()
	defer producerQueue.client.Del(producerQueue.ctx, streamKey)

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

	// Give Redis time to process
	time.Sleep(100 * time.Millisecond)

	// Create consumer queue
	consumerQueue, err := NewRedisQueue(testAddr, streamKey, groupName, testConsumerName+"-consumer")
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
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout waiting for messages. Received %d/%d", receivedCount, len(testMessages))
	}

	// Verify all messages were received
	for _, msg := range testMessages {
		if !receivedMessages[msg.ID] {
			t.Errorf("Message %s was not received", msg.ID)
		}
	}
}

func TestRedisProduceInvalidMessage(t *testing.T) {
	skipIfNoRedis(t)

	queue, err := NewRedisQueue(testAddr, testStream+"-invalid", testConsumerGroup, testConsumerName)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()

	// Test with message that has empty fields
	msg := &common.Message{
		ID:        "",
		Payload:   nil,
		Timestamp: time.Time{},
	}

	// This should not panic
	_ = queue.Produce(msg)
}

func TestRedisGetStreamInfo(t *testing.T) {
	skipIfNoRedis(t)

	streamKey := testStream + "-info"
	queue, err := NewRedisQueue(testAddr, streamKey, testConsumerGroup, testConsumerName)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()
	defer queue.client.Del(queue.ctx, streamKey)

	// Produce some messages
	for i := 0; i < 5; i++ {
		msg := &common.Message{
			ID:        "msg-" + string(rune(i)),
			Payload:   []byte("test"),
			Timestamp: time.Now(),
		}
		queue.Produce(msg)
	}

	info, err := queue.GetStreamInfo()
	if err != nil {
		t.Fatalf("Failed to get stream info: %v", err)
	}

	if info.Length != 5 {
		t.Errorf("Expected stream length 5, got %d", info.Length)
	}
}

func TestRedisTrimStream(t *testing.T) {
	skipIfNoRedis(t)

	streamKey := testStream + "-trim"
	queue, err := NewRedisQueue(testAddr, streamKey, testConsumerGroup, testConsumerName)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()
	defer queue.client.Del(queue.ctx, streamKey)

	// Produce 10 messages
	for i := 0; i < 10; i++ {
		msg := &common.Message{
			ID:        "msg-" + string(rune(i)),
			Payload:   []byte("test"),
			Timestamp: time.Now(),
		}
		queue.Produce(msg)
	}

	// Trim to 5 messages
	err = queue.TrimStream(5)
	if err != nil {
		t.Errorf("Failed to trim stream: %v", err)
	}

	// Verify stream length
	info, err := queue.GetStreamInfo()
	if err != nil {
		t.Fatalf("Failed to get stream info: %v", err)
	}

	if info.Length > 5 {
		t.Errorf("Expected stream length <= 5 after trim, got %d", info.Length)
	}
}

func TestRedisClose(t *testing.T) {
	skipIfNoRedis(t)

	queue, err := NewRedisQueue(testAddr, testStream+"-close", testConsumerGroup, testConsumerName)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
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

func TestRedisHighThroughput(t *testing.T) {
	skipIfNoRedis(t)

	if testing.Short() {
		t.Skip("Skipping high-throughput test in short mode")
	}

	streamKey := testStream + "-throughput"
	queue, err := NewRedisQueue(testAddr, streamKey, testConsumerGroup, testConsumerName)
	if err != nil {
		t.Fatalf("Failed to create Redis queue: %v", err)
	}
	defer queue.Close()
	defer queue.client.Del(queue.ctx, streamKey)

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

	duration := time.Since(startTime)
	throughput := float64(messageCount) / duration.Seconds()

	t.Logf("Produced %d messages in %v (%.2f msg/sec)", messageCount, duration, throughput)

	// Verify all messages were added
	info, err := queue.GetStreamInfo()
	if err != nil {
		t.Fatalf("Failed to get stream info: %v", err)
	}

	if info.Length != int64(messageCount) {
		t.Errorf("Expected %d messages in stream, got %d", messageCount, info.Length)
	}

	if throughput < 100 {
		t.Logf("Warning: Low throughput %.2f msg/sec", throughput)
	}
}

func TestRedisConsumerGroupHandling(t *testing.T) {
	skipIfNoRedis(t)

	streamKey := testStream + "-cg-handling"
	groupName := testConsumerGroup + "-cgh"

	// Create first queue (creates consumer group)
	queue1, err := NewRedisQueue(testAddr, streamKey, groupName, testConsumerName+"-1")
	if err != nil {
		t.Fatalf("Failed to create first Redis queue: %v", err)
	}
	defer queue1.Close()
	defer queue1.client.Del(queue1.ctx, streamKey)

	// Create second queue with same consumer group (should not error)
	queue2, err := NewRedisQueue(testAddr, streamKey, groupName, testConsumerName+"-2")
	if err != nil {
		t.Fatalf("Failed to create second Redis queue with same group: %v", err)
	}
	defer queue2.Close()

	// Both queues should work fine
	msg := &common.Message{
		ID:        "cg-test",
		Payload:   []byte("test"),
		Timestamp: time.Now(),
	}

	if err := queue1.Produce(msg); err != nil {
		t.Errorf("Failed to produce with first queue: %v", err)
	}

	if err := queue2.Produce(msg); err != nil {
		t.Errorf("Failed to produce with second queue: %v", err)
	}
}
