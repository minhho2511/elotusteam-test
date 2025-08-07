package cache

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type Config struct {
	Hosts     []string
	IsCluster bool
	Pass      string
	DB        string
	Debug     bool
}

func (c Config) GetDB() int {
	if c.DB == "" {
		return 0
	}
	db, err := strconv.Atoi(c.DB)
	if err != nil {
		return 0
	}
	return db
}

func NewRedis(c Config, logger log.Logger) (redis.UniversalClient, error) {
	addr := c.Hosts
	opts := redis.UniversalOptions{
		Addrs:        addr,
		Password:     c.Pass,
		DB:           c.GetDB(),
		PoolSize:     100,
		MinIdleConns: 10,
		MaxRetries:   3,
	}
	cluster := c.IsCluster || len(addr) > 1
	var client redis.UniversalClient
	if cluster {
		client = redis.NewClusterClient(opts.Cluster())
	} else {
		client = redis.NewUniversalClient(&opts)
	}

	if c.Debug {
		hook := NewDebugHook(logger)
		client.AddHook(&hook)
	}

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return client, nil
}
