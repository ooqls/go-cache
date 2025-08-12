package store

import (
	"context"
	"testing"
	"time"

	"github.com/ooqls/go-db/redis"
	"github.com/ooqls/go-db/testutils"
	"github.com/stretchr/testify/assert"
)

type Obj struct {
	V string `json:"v"`
}

func TestRedisStore(t *testing.T) {
	testutils.InitRedis()

	store := NewGenericStore(redis.GetConnection(), 10*time.Second)

	obj := Obj{V: "value"}
	err := store.Set(context.Background(), "key", obj)
	assert.Nil(t, err)

	var updatedObj Obj
	err = store.Get(context.Background(), "key", &updatedObj)
	assert.Nil(t, err)
	assert.Equal(t, obj.V, updatedObj.V)

	err = store.Update(context.Background(), "key", func(fn func(target any) error) (any, error) {
		var obj Obj

		err := fn(&obj)
		assert.Nil(t, err)

		obj.V = "updated value"
		return obj, nil
	})
	assert.Nil(t, err)

	err = store.Get(context.Background(), "key", &updatedObj)
	assert.Nil(t, err)
	assert.Equal(t, "updated value", updatedObj.V)

}
