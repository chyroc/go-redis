package internal

import "github.com/chyroc/go-redis/internal/basetype"

type dictEntry struct {
	key  *basetype.SDS
	v    interface{}
	next *dictEntry
}

type dictht struct {
	table    []*dictEntry
	size     uint32
	sizemask uint32
	used     uint32
}

type dictType struct {
}

type dict struct {
	type_       *dictType
	privdata    interface{}
	ht          [2]dictht
	rehashindex int
}
