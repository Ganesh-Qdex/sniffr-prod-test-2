package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// ConnectRedis establishes connection to Redis
func ConnectRedis(addr, password string, db int) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     100, // Connection pool size
		MinIdleConns: 10,  // Minimum idle connections
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis!")
	return nil
}

// DisconnectRedis closes the Redis connection
func DisconnectRedis() error {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			return fmt.Errorf("failed to disconnect from Redis: %v", err)
		}
		log.Println("Disconnected from Redis!")
	}
	return nil
}

// Set stores a value in Redis with expiration
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}

	return RedisClient.Set(ctx, key, jsonValue, expiration).Err()
}

// Get retrieves a value from Redis
func Get(ctx context.Context, key string, dest interface{}) error {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get value: %v", err)
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete removes a key from Redis
func Delete(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// DeletePattern removes all keys matching a pattern
func DeletePattern(ctx context.Context, pattern string) error {
	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys: %v", err)
	}

	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}

// Exists checks if a key exists in Redis
func Exists(ctx context.Context, key string) (bool, error) {
	count, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
