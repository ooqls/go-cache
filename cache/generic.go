package cache

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/eko/gocache/lib/v4/cache"
)

func NewGenericCache(c *cache.Cache[[]byte]) *GenericCache {
	return &GenericCache{
		c: c,
	}
}

type GenericCache struct {
	c *cache.Cache[[]byte]
}

func (c *GenericCache) Set(ctx context.Context, key string, val any) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(val)
	if err != nil {
		return err
	}

	return c.c.Set(ctx, key, buff.Bytes())
}

func (c *GenericCache) Get(ctx context.Context, key string, target any) error {
	b, err := c.c.Get(ctx, key)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(bytes.NewReader(b))
	return dec.Decode(target)
}

func (c *GenericCache) Delete(ctx context.Context, key string) error {
	return c.c.Delete(ctx, key)
}

func (c *GenericCache) Clear(ctx context.Context) error {
	return c.c.Clear(ctx)
}
