package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

func NewRedisConnect(ctx context.Context, addrs []string, password string, poolSize int) (*redis.ClusterClient, error) {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: password,
		PoolSize: poolSize,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
