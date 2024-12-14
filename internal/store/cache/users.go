package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lucianboboc/goBackendEngineering/internal/store"
	"github.com/redis/go-redis/v9"
	"time"
)

type UsersStore struct {
	rds *redis.Client
}

const UserExpTime = time.Minute

func (s *UsersStore) Get(ctx context.Context, id int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", id)
	data, err := s.rds.Get(ctx, cacheKey).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, nil
		default:
			return nil, err
		}
	}

	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *UsersStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rds.Set(ctx, cacheKey, data, UserExpTime).Err()
}
