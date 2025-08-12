package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=store.go -destination=store_mock.go -package=store GenericInterface
type GenericInterface interface {
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string, target any) error
	Update(ctx context.Context, key string, fn func(func(target any) error) (any, error)) error
	Delete(ctx context.Context, key string) error
}

type GenericStore struct {
	db  *redis.Client
	ttl time.Duration
}

func NewGenericStore(db *redis.Client, ttl time.Duration) GenericInterface {
	return &GenericStore{db: db, ttl: ttl}
}

func (s *GenericStore) Set(ctx context.Context, key string, value any) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.db.Set(ctx, key, json, s.ttl).Err()
}

func (s *GenericStore) Get(ctx context.Context, key string, target any) error {
	res, err := s.db.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(res), target)
	return err
}

func (s *GenericStore) Update(ctx context.Context, key string, fn func(func(target any) error) (any, error)) error {
	return s.db.Watch(ctx, func(tx *redis.Tx) error {
		res, err := tx.Get(ctx, key).Result()
		if err != nil {
			return err
		}

		target, err := fn(func(target any) error {
			return json.Unmarshal([]byte(res), target)
		})
		if err != nil {
			return err
		}

		b, err := json.Marshal(target)
		if err != nil {
			return err
		}

		_, err = tx.Set(ctx, key, b, s.ttl).Result()
		return err
	}, key)
}

func (s *GenericStore) Delete(ctx context.Context, key string) error {
	return s.db.Del(ctx, key).Err()
}
