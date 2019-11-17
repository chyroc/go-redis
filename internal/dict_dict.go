package internal

type dictEntry struct {
	key  sdshdr
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
