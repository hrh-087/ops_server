package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"ops-server/global"
)

type RedisConfig struct {
	Addr      string
	Password  string
	Port      int64
	DB        int
	IsCluster bool
}

func NewRedisConn(config RedisConfig) (client redis.UniversalClient, err error) {

	if config.IsCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{config.Addr},
			Password: config.Password,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		})
	}
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		global.OPS_LOG.Error("redis connect ping failed, err:", zap.Error(err))
		return nil, err
	}

	return client, nil
}
