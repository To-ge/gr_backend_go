package database

import (
	"crypto/tls"

	"github.com/To-ge/gr_backend_go/config"
	"github.com/go-redis/redis/v8"
)

type RedisConnector struct {
	Conn *redis.Client
}

func NewRedisConnector() *RedisConnector {
	address := config.LoadConfig().RedisInfo.Address
	client := redis.NewClient(&redis.Options{
		Addr: address,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	return &RedisConnector{
		Conn: client,
	}
}
