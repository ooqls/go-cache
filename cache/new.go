package cache

import (
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/ooqls/go-db/redis"
	gocache "github.com/patrickmn/go-cache"
)

func NewMemCache() *cache.Cache[[]byte] {
	memCache := gocache.New(time.Minute*5, time.Minute*10)
	memStore := gocache_store.NewGoCache(memCache)
	return cache.New[[]byte](memStore)
}

func NewRedisCache() *cache.Cache[[]byte] {
	r := redis.GetConnection()
	redisStore := redis_store.NewRedis(r)
	return cache.New[[]byte](redisStore)
}
