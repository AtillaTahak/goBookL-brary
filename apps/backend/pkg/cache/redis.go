package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

type CacheStats struct {
	Hits     int64
	Misses   int64
	Keys     int64
	Memory   string
	Uptime   string
	Connected bool
}

func NewRedisCache(addr, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 10,
		MinIdleConns: 5,
		MaxRetries: 3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx := context.Background()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Warning: Redis connection failed: %v\n", err)
	}

	return &RedisCache{
		client: rdb,
		ctx:    ctx,
	}
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = r.client.Set(r.ctx, key, jsonValue, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache key %s: %w", key, err)
	}

	return nil
}

func (r *RedisCache) Get(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get cache key %s: %w", key, err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return nil
}

func (r *RedisCache) Delete(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	err := r.client.Del(r.ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache keys: %w", err)
	}

	return nil
}

func (r *RedisCache) Exists(key string) (bool, error) {
	result := r.client.Exists(r.ctx, key)
	if result.Err() != nil {
		return false, fmt.Errorf("failed to check key existence: %w", result.Err())
	}

	return result.Val() > 0, nil
}

func (r *RedisCache) Expire(key string, expiration time.Duration) error {
	err := r.client.Expire(r.ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}

	return nil
}

func (r *RedisCache) Keys(pattern string) ([]string, error) {
	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys with pattern %s: %w", pattern, err)
	}

	return keys, nil
}

func (r *RedisCache) FlushAll() error {
	err := r.client.FlushAll(r.ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to flush all keys: %w", err)
	}

	return nil
}

func (r *RedisCache) Incr(key string) (int64, error) {
	result := r.client.Incr(r.ctx, key)
	if result.Err() != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, result.Err())
	}

	return result.Val(), nil
}

func (r *RedisCache) IncrBy(key string, value int64) (int64, error) {
	result := r.client.IncrBy(r.ctx, key, value)
	if result.Err() != nil {
		return 0, fmt.Errorf("failed to increment key %s by %d: %w", key, value, result.Err())
	}

	return result.Val(), nil
}

func (r *RedisCache) GetStats() (*CacheStats, error) {
	_, err := r.client.Info(r.ctx, "stats", "memory", "server").Result()
	if err != nil {
		return &CacheStats{Connected: false}, fmt.Errorf("failed to get cache stats: %w", err)
	}

	dbSize, _ := r.client.DBSize(r.ctx).Result()

	return &CacheStats{
		Keys:      dbSize,
		Connected: true,
		Memory:    "available via INFO command",
		Uptime:    "available via INFO command",
	}, nil
}

func (r *RedisCache) Ping() error {
	_, err := r.client.Ping(r.ctx).Result()
	if err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

func (r *RedisCache) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	result := r.client.SetNX(r.ctx, key, jsonValue, expiration)
	if result.Err() != nil {
		return false, fmt.Errorf("failed to set key %s: %w", key, result.Err())
	}

	return result.Val(), nil
}

func (r *RedisCache) TTL(key string) (time.Duration, error) {
	result := r.client.TTL(r.ctx, key)
	if result.Err() != nil {
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, result.Err())
	}

	return result.Val(), nil
}
