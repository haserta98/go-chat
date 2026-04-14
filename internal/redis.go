package internal

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
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

func (c *RedisClient) Publish(context context.Context, channel string, message []byte) error {
	return c.client.Publish(context, channel, message).Err()
}

func (c *RedisClient) Subscribe(context context.Context, channel string) (<-chan *redis.Message, error) {
	return c.client.Subscribe(context, channel).Channel(), nil
}
