package factory

import (
	"time"

	"github.com/ooqls/go-cache/cache"
	"github.com/ooqls/go-cache/store"
	"github.com/redis/go-redis/v9"
)

type CacheFactory interface {
	NewCache(key string, ttl time.Duration) cache.GenericCache
	NewStore(key string, ttl time.Duration) store.GenericInterface
}

func NewCacheFactory(rc redis.Client) CacheFactory {
	return &RedisCacheFactory{rc: rc}
}

type RedisCacheFactory struct {
	rc redis.Client
}

func (f *RedisCacheFactory) NewCache(key string, ttl time.Duration) cache.GenericCache {
	return *cache.NewGenericCache(key, cache.NewRedisCache(f.rc, ttl))
}

func (f *RedisCacheFactory) NewStore(key string, ttl time.Duration) store.GenericInterface {
	return store.NewRedisStore(key, f.rc, ttl)
}

func NewMemCacheFactory() CacheFactory {
	return &MemCacheFactory{}
}

type MemCacheFactory struct{}

func (f *MemCacheFactory) NewCache(key string, ttl time.Duration) cache.GenericCache {
	return *cache.NewGenericCache(key, cache.NewMemCache())
}

func (f *MemCacheFactory) NewStore(key string, ttl time.Duration) store.GenericInterface {
	return store.NewMemStore(key, ttl)
}
