package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/praneethys/kafka-bullmq-benchmark/pkg/common"
)

// KafkaQueue implements the MessageQueue interface for Apache Kafka
type KafkaQueue struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
	topic    string
	brokers  string
}

// NewKafkaQueue creates a new Kafka queue instance
func NewKafkaQueue(brokers, topic string, consumerGroup string) (*KafkaQueue, error) {
	// High-throughput producer configuration
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"acks":              "1", // Wait for leader acknowledgment
		"compression.type":  "lz4",
		"linger.ms":         10,
		"batch.size":        1000000,
		"buffer.memory":     67108864, // 64MB
		"max.in.flight":     5,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// High-throughput consumer configuration
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"group.id":           consumerGroup,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
		"fetch.min.bytes":    1024,
		"fetch.max.wait.ms":  100,
		"max.partition.fetch.bytes": 10485760, // 10MB
	})
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	if err := consumer.Subscribe(topic, nil); err != nil {
		producer.Close()
		consumer.Close()
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return &KafkaQueue{
		producer: producer,
		consumer: consumer,
		topic:    topic,
		brokers:  brokers,
	}, nil
}

// Produce sends a message to Kafka
func (k *KafkaQueue) Produce(msg *common.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
		Key:   []byte(msg.ID),
	}

	deliveryChan := make(chan kafka.Event)
	if err := k.producer.Produce(kafkaMsg, deliveryChan); err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	close(deliveryChan)

	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}

	return nil
}

// ProduceAsync sends a message to Kafka without waiting for acknowledgment
func (k *KafkaQueue) ProduceAsync(msg *common.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &k.topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
		Key:   []byte(msg.ID),
	}

	if err := k.producer.Produce(kafkaMsg, nil); err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	return nil
}

// Consume reads messages from Kafka and processes them with the provided handler
func (k *KafkaQueue) Consume(handler func(*common.Message) error) error {
	for {
		msg, err := k.consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				continue
			}
			return fmt.Errorf("consumer error: %w", err)
		}

		var message common.Message
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			continue
		}

		if err := handler(&message); err != nil {
			continue
		}
	}
}

// Flush waits for all messages to be delivered
func (k *KafkaQueue) Flush(timeoutMs int) int {
	return k.producer.Flush(timeoutMs)
}

// Close closes the Kafka producer and consumer
func (k *KafkaQueue) Close() error {
	k.producer.Close()
	return k.consumer.Close()
}

// GetName returns the name of this queue implementation
func (k *KafkaQueue) GetName() string {
	return "Apache Kafka"
}
