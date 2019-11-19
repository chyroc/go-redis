package objects

import (
	"github.com/chyroc/go-redis/internal/basetype"
)

type RedisListObject struct {
	*redisObjectImpl
}

func NewRedisListObject(l []string) *RedisListObject {

	// 只有满足下面两个条件的时候，才能用 ziplist
	// 所有字符串长度都小于 64 字节
	// 元素个数小于 512

	var encoding = RedisObjectEncodingZipList
	var ll = basetype.NewLinkedList()
	for _, v := range l {
		s := basetype.NewSDSWithString(v)
		ll.AddTail()
		if encoding == RedisObjectEncodingZipList && s.Len() >= 64 {
			encoding = RedisObjectEncodingLinkedList
		}
	}

	return &RedisListObject{
		redisObjectImpl: &redisObjectImpl{
			type_:    RedisObjectTypeList,
			encoding: encoding,
			ptr:      ll,
		},
	}
}

func (r *RedisListObject) Len() uint32 {
	return r.linkedList().Len()
}

func (r *RedisListObject) LPush(value interface{}) {
	r.linkedList().AddHead(value)
}

func (r *RedisListObject) RPush(value interface{}) {
	r.linkedList().AddTail(value)
}

func (r *RedisListObject) LPop(value interface{}) interface{} {
	if r.linkedList().Len() == 0 {
		return nil
	}
	data := r.linkedList().First()
	r.linkedList().DelHead()
	return data
}

func (r *RedisListObject) RPop(value interface{}) interface{} {
	if r.linkedList().Len() == 0 {
		return nil
	}
	data := r.linkedList().Last()
	r.linkedList().DelTail()
	return data
}

func (r *RedisListObject) Index(idx uint32) interface{} {
	v := r.linkedList().Index(idx)
	if v == nil {
		return nil
	}
	return v.Value()
}

func (r *RedisListObject) Insert(idx uint32, value interface{}) {
	r.linkedList().Insert(idx, value)
}

//func (r *RedisListObject) Rem(value interface{}) {
//	r.linkedList().Insert(idx, value)
//}

func (r *RedisListObject) linkedList() *basetype.LinkedList {
	return r.ptr.(*basetype.LinkedList)
}
