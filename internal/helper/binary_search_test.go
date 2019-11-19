package helper

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chyroc/go-redis/internal/tests"
)

func TestSearchFirstGreaterEqual(t *testing.T) {
	as := assert.New(t)

	for length := 0; length < 99; length++ {
		var arr IntSequentialSequence
		var dict = make(map[int]bool)
		for i := 0; i < length; i++ {
			data := tests.RandomIntIn100()
			for dict[data] {
				data = tests.RandomIntIn100()
			}
			dict[data] = true
			arr = append(arr, data)
			sort.Ints(arr)
		}

		fmt.Println("arr", arr)

		for idx, v := range arr {
			idx2, found := SearchFirstGreaterEqual(arr, int64(v))
			msg := fmt.Sprintf("length=%d, data=%v, idx=%d, idx2=%d, found=%v\n", len(arr), v, idx, idx2, found)
			t.Logf(msg)
			as.True(found, msg)
			as.Equal(uint32(idx), idx2, msg)
		}
	}
}
