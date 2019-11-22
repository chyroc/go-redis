package server_test

import (
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"

	"github.com/chyroc/go-redis/internal/tests"
)

func TestGet(t *testing.T) {
	as := assert.New(t)

	t.Run("获取 空", func(t *testing.T) {
		c := client.Get(tests.RandString32())
		as.Equal(redis.Nil, c.Err())
	})

	t.Run("设置，然后获取", func(t *testing.T) {
		k, v := tests.RandString32(), tests.RandString32()
		c := client.Set(k, v, 0)
		as.Nil(c.Err())
		as.Equal("OK", c.Val())

		c2 := client.Get(k)
		as.Nil(c2.Err())
		as.Equal(v, c2.Val())
	})

	t.Run("重复设置", func(t *testing.T) {
		k, v := tests.RandString32(), tests.RandString32()
		c := client.Set(k, v, 0)
		as.Nil(c.Err())
		as.Equal("OK", c.Val())

		c2 := client.Set(k, tests.RandString32(), 0)
		as.Nil(c2.Err())
		as.Equal("OK", c2.Val())
	})
}

func TestGetSet(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		k := tests.RandString32()

		c1 := client.Get(k)
		as.Equal(redis.Nil, c1.Err())

		v := tests.RandString32()
		c2 := client.GetSet(k, v)
		as.Equal(redis.Nil, c2.Err())

		c3 := client.Get(k)
		as.Nil(c3.Err())
		as.Equal(v, c3.Val())
	})
}
