package redis

import (
	"biz-auto-api/pkg/logger"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	_client *redis.Client
)

func Setup(username, password, host string, db, port, timeout, poolSize int) {
	client, err := NewRedisClient(username, password, host, db, port, timeout, poolSize)
	if err != nil {
		logger.GetLogger().WithField("engine", "redis").Fatalf("redis connect error: %+v", err)
	}
	_client = client
	logger.GetLogger().WithField("engine", "redis").Info("init redis connection success")
}
func GetClient() *redis.Client {
	return _client
}

func NewRedisClient(username, password, host string, db, port, timeout, poolSize int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%v:%v", host, port),
		Password:     password,
		DB:           db,
		ReadTimeout:  time.Second * time.Duration(timeout),
		WriteTimeout: time.Second * time.Duration(timeout),
		DialTimeout:  time.Second * time.Duration(timeout),
		PoolSize:     poolSize,
		Username:     username,
	})
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return client, nil
}
