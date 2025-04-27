package cache

import (
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/ooqls/go-db/redis"
	gocache "github.com/patrickmn/go-cache"
)

func NewMemCache[T any]() *cache.Cache[T] {
	memCache := gocache.New(time.Minute*5, time.Minute*10)
	memStore := gocache_store.NewGoCache(memCache)
	return cache.New[T](memStore)
}

func NewRedisCache[T any]() *cache.Cache[T] {
	r := redis.GetConnection()
	redisStore := redis_store.NewRedis(r)
	return cache.New[T](redisStore)
}
