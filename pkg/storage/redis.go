package storage

import (
	"fmt"
	"net"
	"strings"

	"github.com/go-redis/redis/v7"
)

// Redis is a Redis client, a temporary replacement for a proper DHT
type Redis struct {
	client *redis.Client
}

// NewRedis creates a new Redis client
func NewRedis(addr *net.TCPAddr, password string) *Redis {
	return &Redis{
		redis.NewClient(&redis.Options{
			Addr:     addr.String(),
			Password: password,
			DB:       0,
		}),
	}
}

// Apply applies relation deletions and puts
func (r *Redis) Apply(deletions [][2]string, additions [][2]string) error {
	for _, deletion := range deletions {
		if err := r.client.Del(fmt.Sprintf("node:%v:%v", deletion[0], deletion[1])).Err(); err != nil {
			return err
		}
	}

	for _, addition := range additions {
		if err := r.client.Set(fmt.Sprintf("node:%v:%v", addition[0], addition[1]), true, 0).Err(); err != nil {
			return err
		}
	}

	return nil
}

// GetAll returns all relations, a maximum of 1000 keys is supported
func (r *Redis) GetAll() ([][2]string, error) {
	keys, _, err := r.client.Scan(0, "node:*", 1000).Result()
	if err != nil {
		return nil, err
	}

	out := [][2]string{}
	for _, key := range keys {
		line := [2]string{}

		from := strings.Split(key, ":")
		line[0] = from[1]
		line[1] = from[2]

		out = append(out, line)
	}

	return out, nil
}
