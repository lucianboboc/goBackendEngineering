package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisClient(addr, pass string, db int) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	ctx := context.Background()
	if _, err := cli.Ping(ctx).Result(); err != nil {
		log.Fatal(err)
		return nil
	}
	
	return cli
}
