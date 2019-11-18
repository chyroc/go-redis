package basetype

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinkedList(t *testing.T) {
	as := assert.New(t)

	list := NewLinkedList()
	as.Equal(list.Len(), uint32(0))
}
