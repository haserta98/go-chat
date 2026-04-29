package internal

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	redisUrl := os.Getenv("REDIS_URL")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)
	if err != nil {
		redisDB = 0
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: redisPassword,
		DB:       int(redisDB),
	})
	return &RedisClient{
		client: client,
	}
}

func (c *RedisClient) Set(key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(context.Background(), key, value, ttl).Err()
}

func (c *RedisClient) Get(key string) (string, error) {
	return c.client.Get(context.Background(), key).Result()
}

func (c *RedisClient) SAdd(key string, members ...string) error {
	return c.client.SAdd(context.Background(), key, members).Err()
}

func (c *RedisClient) SRem(key string, members ...string) error {
	return c.client.SRem(context.Background(), key, members).Err()
}

func (c *RedisClient) SMembers(key string) ([]string, error) {
	return c.client.SMembers(context.Background(), key).Result()
}

func (c *RedisClient) Del(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

func (c *RedisClient) Exists(key string) (bool, error) {
	val, err := c.client.Exists(context.Background(), key).Result()
	return val > 0, err
}

