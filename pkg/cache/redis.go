package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewRedisClient(host, port string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
