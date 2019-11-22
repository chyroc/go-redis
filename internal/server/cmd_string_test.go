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

func TestStrLen(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		k := tests.RandString32()

		c1 := client.StrLen(k)
		as.Nil(c1.Err())
		as.Equal(int64(0), c1.Val())

		v := tests.RandString32()
		as.Nil(client.Set(k, v, 0).Err())

		c3 := client.StrLen(k)
		as.Nil(c3.Err())
		as.Equal(int64(len(v)), c3.Val())
	})
}

func TestAppend(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		k := tests.RandString32()
		v := tests.RandString32()

		c1 := client.Append(k, v)
		as.Nil(c1.Err())
		as.Equal(int64(len(v)), c1.Val())

		v2 := tests.RandString32()
		c2 := client.Append(k, v2)
		as.Nil(c2.Err())
		as.Equal(int64(len(v+v2)), c2.Val())
	})
}

func TestIncr(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		k := tests.RandString32()

		c1 := client.Incr(k)
		as.Nil(c1.Err())
		as.Equal(int64(1), c1.Val())
	})

	t.Run("", func(t *testing.T) {
		k := tests.RandString32()
		as.Nil(client.Set(k, 100, 0).Err())

		for i := 0; i < 1200; i++ {
			c1 := client.Incr(k)
			as.Nil(c1.Err())
			as.Equal(int64(100+i+1), c1.Val())
		}
	})
}

func TestMGetSet(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		c1 := client.MGet(tests.RandString32())
		as.Nil(c1.Err())
		as.Len(c1.Val(), 1)
		as.Nil(c1.Val()[0])
	})

	t.Run("", func(t *testing.T) {
		c1 := client.MGet(tests.RandString32(), tests.RandString32())
		as.Nil(c1.Err())
		as.Len(c1.Val(), 2)
		as.Nil(c1.Val()[0])
		as.Nil(c1.Val()[1])
	})

	t.Run("", func(t *testing.T) {
		k := tests.RandString32()

		as.Nil(client.Set(k, "x", 0).Err())

		c1 := client.MGet(k, tests.RandString32())
		as.Nil(c1.Err())
		as.Len(c1.Val(), 2)
		as.Equal("x", c1.Val()[0])
		as.Nil(c1.Val()[1])
	})
}
