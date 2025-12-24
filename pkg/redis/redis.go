package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/praneethys/kafka-bullmq-benchmark/pkg/common"
)

// RedisQueue implements the MessageQueue interface using Redis Streams (BullMQ equivalent)
type RedisQueue struct {
	client      *redis.Client
	streamKey   string
	consumerGroup string
	consumerName  string
	ctx         context.Context
}

// NewRedisQueue creates a new Redis queue instance using Redis Streams
func NewRedisQueue(addr, streamKey, consumerGroup, consumerName string) (*RedisQueue, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		PoolSize:     100,
		MinIdleConns: 10,
		MaxRetries:   3,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	rq := &RedisQueue{
		client:       client,
		streamKey:    streamKey,
		consumerGroup: consumerGroup,
		consumerName:  consumerName,
		ctx:          ctx,
	}

	// Create consumer group (ignore error if already exists)
	client.XGroupCreateMkStream(ctx, streamKey, consumerGroup, "0")

	return rq, nil
}

// Produce sends a message to Redis Stream
func (r *RedisQueue) Produce(msg *common.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	args := &redis.XAddArgs{
		Stream: r.streamKey,
		Values: map[string]interface{}{
			"id":        msg.ID,
			"payload":   data,
			"timestamp": msg.Timestamp.Unix(),
		},
	}

	if _, err := r.client.XAdd(r.ctx, args).Result(); err != nil {
		return fmt.Errorf("failed to add message to stream: %w", err)
	}

	return nil
}

// ProduceAsync sends a message to Redis Stream (same as Produce for Redis)
func (r *RedisQueue) ProduceAsync(msg *common.Message) error {
	return r.Produce(msg)
}

// Consume reads messages from Redis Stream and processes them with the provided handler
func (r *RedisQueue) Consume(handler func(*common.Message) error) error {
	for {
		// Read from consumer group
		streams, err := r.client.XReadGroup(r.ctx, &redis.XReadGroupArgs{
			Group:    r.consumerGroup,
			Consumer: r.consumerName,
			Streams:  []string{r.streamKey, ">"},
			Count:    10,
			Block:    100 * time.Millisecond,
		}).Result()

		if err != nil {
			if err == redis.Nil {
				continue
			}
			return fmt.Errorf("consumer error: %w", err)
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				var msg common.Message

				if data, ok := message.Values["payload"].(string); ok {
					if err := json.Unmarshal([]byte(data), &msg); err != nil {
						continue
					}
				} else {
					continue
				}

				if err := handler(&msg); err != nil {
					continue
				}

				// Acknowledge the message
				r.client.XAck(r.ctx, r.streamKey, r.consumerGroup, message.ID)
			}
		}
	}
}

// Close closes the Redis client connection
func (r *RedisQueue) Close() error {
	return r.client.Close()
}

// GetName returns the name of this queue implementation
func (r *RedisQueue) GetName() string {
	return "Redis Streams (BullMQ)"
}

// GetStreamInfo returns information about the stream
func (r *RedisQueue) GetStreamInfo() (*redis.XInfoStream, error) {
	return r.client.XInfoStream(r.ctx, r.streamKey).Result()
}

// TrimStream trims the stream to a maximum length
func (r *RedisQueue) TrimStream(maxLen int64) error {
	return r.client.XTrimMaxLen(r.ctx, r.streamKey, maxLen).Err()
}
