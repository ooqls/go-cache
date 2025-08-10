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

	store := NewGenericStore[Obj](redis.GetConnection(), 10*time.Second)

	obj := Obj{V: "value"}
	err := store.Set(context.Background(), "key", obj)
	assert.Nil(t, err)

	updatedObj, err := store.Get(context.Background(), "key")
	assert.Nil(t, err)
	assert.Equal(t, obj.V, updatedObj.V)

	err = store.Update(context.Background(), "key", func(obj *Obj) {
		obj.V = "updated value"
	})
	assert.Nil(t, err)

	updatedObj, err = store.Get(context.Background(), "key")
	assert.Nil(t, err)
	assert.Equal(t, "updated value", updatedObj.V)

}
