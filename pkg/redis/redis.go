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
	SETNX(key string) (bool, error)
	HINCRBY(value, key string) error
	HGETALL(value string) (map[string]string, error)
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

// SETNX set a key from the Redis database and return if the key is new or not
func (c *Client) SETNX(key string) (bool, error) {
	conn := c.pool.Get()
	defer conn.Close()

	isNew, err := redis.Bool(conn.Do("SETNX", key, 1))
	if err != nil {
		return false, err
	}

	return isNew, nil
}

// HINCRBY increase a key from the Redis database based on key
func (c *Client) HINCRBY(value, key string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("HINCRBY", value, key, 1)
	if err != nil {
		return err
	}

	return err
}

// HGETALL is func to get all metrics from Redis database based on key
func (c *Client) HGETALL(value string) (map[string]string, error) {
	var values map[string]string
	conn := c.pool.Get()
	defer conn.Close()

	values, err := redis.StringMap(conn.Do("HGETALL", value))
	if err != nil {
		fmt.Println(err)
		return values, err
	}
	return values, err
}
