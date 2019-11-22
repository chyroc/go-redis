package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/chyroc/go-redis/internal/database"
	"github.com/chyroc/go-redis/internal/tests"
)

func TestTTL(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		k := tests.RandString32()
		msg := fmt.Sprintf("k=%q", k)
		as.Nil(client.Set(k, tests.RandString32(), 0).Err(), msg)

		c2 := client.TTL(k)
		as.Nil(c2.Err(), msg)
		as.Equal(time.Duration(database.TimeNeverExpire), c2.Val(), msg)
	})

	t.Run("", func(t *testing.T) {
		k := tests.RandString32()
		as.Nil(client.Set(k, tests.RandString32(), time.Second/10).Err())

		time.Sleep(time.Second / 5)

		c2 := client.TTL(k)
		as.Nil(c2.Err())
		as.Equal(time.Duration(database.TimeExpired), c2.Val())
	})

	t.Run("", func(t *testing.T) {
		k := tests.RandString32()
		msg := fmt.Sprintf("k=%q", k)
		as.Nil(client.Set(k, tests.RandString32(), time.Second*2).Err(), msg)

		// time.Sleep(time.Second / 5)

		c2 := client.TTL(k)
		as.Nil(c2.Err(), msg)
		as.True(c2.Val() == time.Second || c2.Val() == time.Second*2, msg)
	})
}
