package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/org/2112-space-lab/org/app-service/internal/config"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
)

func (r *RedisClient) Init() {}

// RedisClient wraps the Redis client to handle low-level operations.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient initializes a new Redis client.
func NewRedisClient(env *config.SEnv) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr: env.EnvVars.Redis.GetAddr(),
	})

	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

// HSet stores a hash in Redis.
func (r *RedisClient) HSet(ctx context.Context, key string, data map[string]interface{}) error {
	if err := r.client.HMSet(key, data).Err(); err != nil {
		return fmt.Errorf("failed to store hash: %w", err)
	}
	return nil
}

// HGetAll retrieves a hash from Redis.
func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve hash: %w", err)
	}
	return data, nil
}

// Del removes a key from Redis.
func (r *RedisClient) Del(ctx context.Context, key string) error {
	if err := r.client.Del(key).Err(); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}
	return nil
}

// Exists checks if a key exists in Redis.
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key: %w", err)
	}
	return count > 0, nil
}

// Expire sets an expiration time on a key in Redis.
func (r *RedisClient) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if err := r.client.Expire(key, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set expiration on key: %w", err)
	}
	return nil
}

// Get retrieves a value for a given key from Redis.
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Key does not exist
		}
		return "", fmt.Errorf("failed to get value: %w", err)
	}
	return value, nil
}

// Set sets a value for a key in Redis.
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	if err := r.client.Set(key, value, 0).Err(); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}
	return nil
}

// Publish sends a message to a Redis Pub/Sub channel.
func (r *RedisClient) Publish(ctx context.Context, channel string, message interface{}) error {
	// Log the publishing attempt
	log.Debugf("Publishing message to channel %s\n", channel)

	// Attempt to publish the message
	if err := r.client.Publish(channel, message).Err(); err != nil {
		log.Errorf("Failed to publish message to channel %s: %v\n", channel, err)
		return fmt.Errorf("failed to publish message to channel %s: %w", channel, err)
	}

	// Log successful publishing
	log.Debugf("Successfully published message to channel %s\n", channel)
	return nil
}

// Subscribe listens for messages on a Redis Pub/Sub channel and processes them using the provided handler.
func (r *RedisClient) Subscribe(ctx context.Context, channel string, handler func(string) error) error {
	pubsub := r.client.Subscribe(channel)

	// Ensure the subscription is successfully created
	_, err := pubsub.Receive()
	if err != nil {
		return fmt.Errorf("failed to subscribe to channel %s: %w", channel, err)
	}

	log.Debugf("Subscribed to Redis channel: %s\n", channel)

	// Start listening for messages
	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			log.Debugf("Received message from channel %s: %s\n", channel, msg.Payload)

			// Pass the message to the handler function
			if err := handler(msg.Payload); err != nil {
				log.Errorf("Error processing message from channel %s: %v\n", channel, err)
			}

		case <-ctx.Done():
			// Unsubscribe and stop listening if context is canceled
			if err := pubsub.Close(); err != nil {
				log.Errorf("Error unsubscribing from channel %s: %v\n", channel, err)
			}
			log.Debugf("Unsubscribed from Redis channel: %s\n", channel)
			return nil
		}
	}
}

// ZRangeByScore retrieves members of a sorted set by score range.
func (r *RedisClient) ZRangeByScore(ctx context.Context, key string, min, max string) ([]string, error) {
	// Query Redis for members within the specified score range
	results, err := r.client.ZRangeByScore(key, redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to query Redis for key %s with score range [%s, %s]: %w", key, min, max, err)
	}

	return results, nil
}

// Lock tries to acquire a lock on a given key with expiration.
func (r *RedisClient) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	locked, err := r.client.SetNX(key, "locked", ttl).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set lock for %s: %w", key, err)
	}
	if locked {
		log.Debugf("Lock acquired for %s", key)
	} else {
		log.Debugf("Lock already exists for %s", key)
	}
	return locked, nil
}

// Unlock releases the lock on a given key.
func (r *RedisClient) Unlock(ctx context.Context, key string) error {
	deleted, err := r.client.Del(key).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock for %s: %w", key, err)
	}
	if deleted > 0 {
		log.Debugf("Lock released for %s", key)
	} else {
		log.Debugf("No lock found for %s", key)
	}
	return nil
}

