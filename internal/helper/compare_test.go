package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompare(t *testing.T) {
	as := assert.New(t)

	t.Run("CompareInt", func(t *testing.T) {
		as.Equal(1, int(Compare(1, -1)))
		as.Equal(-1, int(Compare(-2, -1)))
		as.Equal(0, int(Compare(1, 1)))
	})

	t.Run("CompareFloat64", func(t *testing.T) {
		as.Equal(1, int(Compare(1, -1)))
		as.Equal(-1, int(Compare(-2, -1)))
		as.Equal(0, int(Compare(1, 1)))
	})

	t.Run("CompareString", func(t *testing.T) {
		as.Equal(0, int(Compare("", "")))
		as.Equal(0, int(Compare("1", "1")))

		as.Equal(1, int(Compare("12", "1")))
		as.Equal(1, int(Compare("2", "1")))
		as.Equal(1, int(Compare("2", "12")))
		as.Equal(-1, int(Compare("1", "12")))
	})
}
