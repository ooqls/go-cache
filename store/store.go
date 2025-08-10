package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type GenericInterface[T any] interface {
	Set(ctx context.Context, key string, value T) error
	Get(ctx context.Context, key string) (*T, error)
	Update(ctx context.Context, key string, fn func(*T)) error
	Delete(ctx context.Context, key string) error
}

type GenericStore[T any] struct {
	db  *redis.Client
	ttl time.Duration
}

func NewGenericStore[T any](db *redis.Client, ttl time.Duration) GenericInterface[T] {
	return &GenericStore[T]{db: db, ttl: ttl}
}

func (s *GenericStore[T]) Set(ctx context.Context, key string, value T) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.db.Set(ctx, key, json, s.ttl).Err()
}

func (s *GenericStore[T]) Get(ctx context.Context, key string) (*T, error) {
	res, err := s.db.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var t T
	err = json.Unmarshal([]byte(res), &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *GenericStore[T]) Update(ctx context.Context, key string, fn func(*T)) error {
	return s.db.Watch(ctx, func(tx *redis.Tx) error {
		res, err := tx.Get(ctx, key).Result()
		if err != nil {
			return err
		}

		var t T
		if err := json.Unmarshal([]byte(res), &t); err != nil {
			return err
		}

		fn(&t)

		b, err := json.Marshal(t)
		if err != nil {
			return err
		}

		_, err = tx.Set(ctx, key, b, s.ttl).Result()
		return err
	}, key)
}

func (s *GenericStore[T]) Delete(ctx context.Context, key string) error {
	return s.db.Del(ctx, key).Err()
}
