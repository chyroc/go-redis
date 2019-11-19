package basetype

import "fmt"

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
	head := &ListNode{}
	tail := &ListNode{}
	head.next = tail
	tail.prev = head

	return &LinkedList{
		head: head,
		tail: tail,
		len:  0,
	}
}

// 长度
func (r *LinkedList) Len() uint32 {
	return r.len
}

// 长度
func (r *ListNode) Value() interface{} {
	return r.value
}

// 最后一个节点
func (r *LinkedList) Last() *ListNode {
	return r.tail
}

// 第一个节点
func (r *LinkedList) First() *ListNode {
	return r.head
}

// 添加到表头
func (r *LinkedList) AddHead(value interface{}) {
	r.Insert(0, value)
}

// 添加到表尾
func (r *LinkedList) AddTail(value interface{}) {
	r.Insert(r.len, value)
}

// 删除第一个节点
func (r *LinkedList) DelHead() {
	if r.len == 0 {
		return
	}

	r.head.next = r.head.next.next
	r.head.next.prev = r.head
	r.len--
}

// 删除第一个节点
func (r *LinkedList) DelTail() {
	if r.len == 0 {
		return
	}

	r.tail.prev.prev.next = r.tail
	r.tail.prev = r.tail.prev.prev
	r.len--
}

// 获取第 idx 个节点，idx: [0, n)
func (r *LinkedList) Index(idx uint32) *ListNode {
	if idx < 0 || idx >= r.len {
		panic(fmt.Sprintf("linkedlist 尝试获取 %d，但是长度是 %d", idx, r.len))
	}

	mid := r.len / 2
	var i uint32
	var node *ListNode
	if idx < mid {
		var node *ListNode = r.head
		for i = 0; i <= idx; i++ {
			node = node.next
		}
	} else {
		var node *ListNode = r.tail
		for i = r.len - 1; i >= idx; i-- {
			node = node.prev
		}
	}

	return node
}

// 插入第 idx 个节点，idx: [0, n]
// 当 idx 等于 n 的时候，表示插入最后一个节点
func (r *LinkedList) Insert(idx uint32, value interface{}) {
	var node = &ListNode{value: value}
	var prev *ListNode
	var next *ListNode
	if idx == r.len {
		prev = r.tail.prev
		next = r.tail
	} else if idx == 0 {
		prev = r.head
		next = r.head.next
	} else {
		prev = r.Index(idx)
		next = prev.next
	}

	node.next = node
	node.prev = prev
	prev.next = node
	next.prev = node
	r.len++
}
