package basetype

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSkipList(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		s := NewSkipList()
		for i := 0; i < 1000; i++ {
			as.Nil(s.Add(strconv.Itoa(i), float64(i)))
		}

		for i := 0; i < 1000; i++ {
			nodes, err := s.Get(float64(i))
			as.Nil(err)
			as.Len(nodes, 1)
			as.Equal(float64(i), nodes[0].score)
		}
	})

	t.Run("", func(t *testing.T) {
		s := NewSkipList()
		as.Nil(s.Add("1", 1))
		fmt.Println(s.String())
		as.Equal(ErrDataRepeated, s.Add("1", 1))
	})
}
