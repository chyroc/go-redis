package basetype_test

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/basetype"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDict(t *testing.T) {
	as := assert.New(t)
	dict := basetype.NewDict()

	m := []string{}
	for i := 0; i < 20; i++ {
		s := RandStringRunes(RandomInt())
		dict.Set(s, s)
		m = append(m, s)
	}

	for idx, v := range m {
		value := dict.Get(v)
		as.NotNil(value, fmt.Sprintf("%d: %q find nil data", idx, v))
		as.Equal(value.(string), v, fmt.Sprintf("%d: %q find nil data", idx, v))
	}
	t.Logf("dict is %+v", dict)
}
