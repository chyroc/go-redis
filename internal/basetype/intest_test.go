package basetype

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/tests"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestIntset(t *testing.T) {
	as := assert.New(t)

	is := NewIntset()
	dict := make(map[int64]bool)
	arr := []int64{}
	var getInt = func() int64 {
		var i = tests.RandomInt64(math.MaxInt16 - 1)
		for dict[i] {
			i = tests.RandomInt64(math.MaxInt16 - 1)
		}
		dict[i] = true
		arr = append(arr, i)
		return i
	}

	lookcount := 50
	for i := 0; i < lookcount; i++ {
		is.Add(int16(getInt()))
	}
	for i := 0; i < lookcount; i++ {
		is.Add(int32(getInt()))
	}
	for i := 0; i < lookcount; i++ {
		is.Add(int64(getInt()))
	}
	for k := range dict {
		msg := fmt.Sprintf("[intset] encoding=%d, length=%d", is.encoding, is.length)
		//if is.length < 10 {
		msg = fmt.Sprintf("%s, contents=%v, arr=%v, intset=%v", msg, is.contents, arr, is.Int64Array())
		//}
		as.True(is.Exist(k), msg)
	}
}
