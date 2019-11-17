package basetype_test

import (
	"github.com/chyroc/go-redis/internal/basetype"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinkedList(t *testing.T) {
	as := assert.New(t)

	list := basetype.NewLinkedList()
	as.Equal(list.Len(), uint32(0))
}
