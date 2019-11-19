package basetype

import (
	"github.com/chyroc/go-redis/internal/tests"
	"testing"
)

import "github.com/stretchr/testify/assert"

func TestSDS(t *testing.T) {
	as := assert.New(t)

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
}
