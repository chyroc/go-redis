package objects_test

import (
	"github.com/chyroc/go-redis/internal/objects"
	"github.com/chyroc/go-redis/internal/tests"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestRedisObjectString(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		for i := 0; i <= 32; i++ {
			r := objects.NewRedisStringObject(tests.RandStringRunes(i))
			as.Equal(objects.RedisObjectTypeString, r.Type())
			as.Equal(objects.RedisObjectEncodingEmbStr, r.Encoding())
		}

		for i := 33; i <= 100; i++ {
			r := objects.NewRedisStringObject(tests.RandStringRunes(i))
			as.Equal(objects.RedisObjectTypeString, r.Type())
			as.Equal(objects.RedisObjectEncodingRaw, r.Encoding())
		}
	})

	t.Run("", func(t *testing.T) {
		for i := 0; i <= 100; i++ {
			r := objects.NewRedisStringObject(strconv.Itoa(i))
			as.Equal(objects.RedisObjectTypeString, r.Type())
			as.Equal(objects.RedisObjectEncodingInt, r.Encoding())
		}
	})
}
