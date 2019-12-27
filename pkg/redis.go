package pkg

import (
	"github.com/go-redis/redis/v7"
)

type Redis struct {
	client redis.Client
}

func (d *Redis) Connect(hostPort, password string) {
	d.client = *redis.NewClient(&redis.Options{
		Addr:     hostPort,
		Password: password,
		DB:       0,
	})
}

func (d *Redis) RegisterNode(macAddress, tcpReadHostPort string) error {
	return d.client.Set(macAddress, tcpReadHostPort, 0).Err()
}
