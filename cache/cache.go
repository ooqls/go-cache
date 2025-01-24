package cache

import (
	"sync"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	gocache_store "github.com/eko/gocache/store/go_cache/v4"
	"github.com/eko/gocache/store/redis/v4"
	gocache "github.com/patrickmn/go-cache"
	goredis "github.com/redis/go-redis/v9"
)

var m sync.Mutex
var opt CacheOptions = CacheOptions{
	defaultStore: gocache_store.NewGoCache(gocache.New(time.Minute*10, time.Minute*20)),
}

type CacheOptions struct {
	redisOptions *RedisOptions
	defaultStore store.StoreInterface
}

type Cache[T any] cache.CacheInterface[T]

func SetRedis(redisOpt RedisOptions) {
	m.Lock()
	defer m.Unlock()

	opt.redisOptions = &redisOpt
}

func NewRedisCache[V any]() *cache.Cache[V] {
	m.Lock()
	defer m.Unlock()

	if opt.redisOptions == nil {
		panic("redis options is not set")
	}

	c := cache.New[V](redis.NewRedis(goredis.NewClient(opt.redisOptions.GetGoRedisOptions())))
	return c
}

func NewCache[V any]() *cache.Cache[V] {
	c := cache.New[V](opt.defaultStore)
	return c
}
