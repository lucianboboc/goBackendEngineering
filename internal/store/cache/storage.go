package cache

import (
	"context"
	"github.com/lucianboboc/goBackendEngineering/internal/store"
	"github.com/redis/go-redis/v9"
)

type UsersStorage interface {
	Get(ctx context.Context, id int64) (*store.User, error)
	Set(ctx context.Context, user *store.User) error
}

type Storage struct {
	Users UsersStorage
}

func NewRedisStorage(rds *redis.Client) Storage {
	return Storage{
		Users: &UsersStore{rds: rds},
	}
}
