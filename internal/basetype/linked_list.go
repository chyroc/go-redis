package basetype

type ListNode struct {
	prev  *ListNode
	next  *ListNode
	value interface{}
}

type LinkedList struct {
	head *ListNode
	tail *ListNode
	len  uint32
}

// 创建一个新的双向链表
func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

// 长度
func (l *LinkedList) Len() uint32 {
	return l.len
}

// 最后一个节点
func (l *LinkedList) Last() *ListNode {
	return l.tail
}

// 第一个节点
func (l *LinkedList) First() *ListNode {
	return l.head
}
