package cache

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

func IsCacheMissErr(err error) bool {
	return errors.Is(err, redis.Nil)
}
