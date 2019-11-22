package basetype

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/tests"
	"testing"
)

import "github.com/stretchr/testify/assert"

func TestSDS(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		sds := NewSDS()
		as.Equal(sds.String(), "")
		as.Equal(sds.Len(), 0)

		total := 0
		for i := 0; i < 100; i++ {
			s := tests.RandStringRunes(tests.RandomIntIn100())
			sds.Append(s)
			total += len(s)
		}
		as.Equal(sds.Len(), total)
		//t.Logf("sds is %s", sds.string())
	})

	t.Run("", func(t *testing.T) {
		sds := NewSDSWithString("0")
		x, err := sds.Int64()
		as.Nil(err)
		as.Equal(int64(0), x)

		for i := 0; i < 1200; i++ {
			msg := fmt.Sprintf("i=%d, sds=%s", i, sds.String())

			ii, err := sds.Int64Incr()
			as.Nil(err, msg)
			as.Equal(int64(i+1), ii, msg)

			x, err := sds.Int64()
			as.Nil(err, msg)
			as.Equal(int64(i+1), x, msg)
		}
	})
}
