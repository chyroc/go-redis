package server_test

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/chyroc/go-redis/internal/tests"
)

func TestHGetHSet(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		hash := tests.RandString32()
		field := tests.RandString32()
		msg := fmt.Sprintf("hash=%q, field=%q", hash, field)

		as.Equal(redis.Nil, client.HGet(hash, field).Err(), msg)

		c1 := client.HSet(hash, field, tests.RandString32())
		as.Nil(c1.Err())
		as.Equal(true, c1.Val())

		c2 := client.HSet(hash, field, tests.RandString32())
		as.Nil(c2.Err())
		as.Equal(false, c2.Val())
	})
}
