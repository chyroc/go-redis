package server_test

import (
	"github.com/chyroc/go-redis/internal/server"
	"github.com/chyroc/go-redis/internal/tests"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var client *redis.Client
var addr = ":9091"

func init() {
	go func() {
		if err := server.New(":9091").Run(); err != nil {
			panic(err)
		}
	}()

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	time.Sleep(time.Second / 2)
}

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
