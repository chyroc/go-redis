package internal

type listNode struct {
	prev  *listNode
	next  *listNode
	value interface{}
}

type list struct {
	head *listNode
	tail *listNode
	len  uint32
}
