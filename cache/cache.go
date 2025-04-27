package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"sync"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	gocache "github.com/patrickmn/go-cache"
)

var m sync.Mutex
var opt CacheOptions = CacheOptions{
	defaultStore: gocache_store.NewGoCache(gocache.New(time.Minute*10, time.Minute*20)),
}

type CacheOptions struct {
	redisOptions *RedisOptions
	defaultStore store.StoreInterface
}

func New[V any](c *cache.Cache[[]byte]) *Cache[V] {
	return &Cache[V]{
		c: cache.Cache[[]byte](*c),
	}
}

type Cache[T any] struct {
	c cache.Cache[[]byte]
}

func (c *Cache[T]) Set(ctx context.Context, key string, val T) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(val)
	if err != nil {
		return err
	}

	return c.c.Set(ctx, key, buff.Bytes())
}

func (c *Cache[T]) Get(ctx context.Context, key string) (*T, error) {
	payload, err := c.c.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get item from cache: %v", err)
	}

	dec := gob.NewDecoder(bytes.NewBuffer(payload))
	var obj T
	err = dec.Decode(&obj)
	return &obj, err

}

func (c *Cache[T]) Clear(ctx context.Context) error {
	return c.c.Clear(ctx)
}

func (c *Cache[T]) Delete(ctx context.Context, key string) error {
	return c.c.Delete(ctx, key)
}

func SetRedis(redisOpt RedisOptions) {
	m.Lock()
	defer m.Unlock()

	opt.redisOptions = &redisOpt
}
