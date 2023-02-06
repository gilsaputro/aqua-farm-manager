package redis

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisMethod list is all available method for redis
type RedisMethod interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Delete(key string) error
}

// RedisConfig is list config to create redis client
type RedisConfig struct {
	RedisHost        string
	Password         string
	MaxIdleInSec     int64
	IdleTimeoutInSec int64
}

// Client is a wrapper for Redigo Redis client
type Client struct {
	pool *redis.Pool
}

// NewRedisClient func to creates a new Redis client
func NewRedisClient(cfg RedisConfig) (RedisMethod, error) {
	var err error

	pool := &redis.Pool{
		MaxIdle:     int(cfg.MaxIdleInSec),
		IdleTimeout: time.Duration(cfg.IdleTimeoutInSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.RedisHost, redis.DialPassword(cfg.Password), redis.DialDatabase(0))
		},
	}

	return &Client{
		pool: pool,
	}, err
}

// Get retrieves a value from the Redis database
func (c *Client) Get(key string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", fmt.Errorf("error retrieving key %s: %v", key, err)
	}

	return value, nil
}

// Set stores a value in the Redis database
func (c *Client) Set(key, value string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("error setting key %s: %v", key, err)
	}

	return nil
}

// Delete removes a key-value pair from the Redis database
func (c *Client) Delete(key string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return fmt.Errorf("error deleting key %s: %v", key, err)
	}

	return nil
}
