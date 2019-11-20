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

	list := &RedisListObject{
		redisObjectImpl: &redisObjectImpl{
			type_:    RedisObjectTypeList,
			encoding: RedisObjectEncodingLinkedList, // TODO: ziplist
			ptr:      basetype.NewLinkedList(),
		},
	}
	for _, v := range l {
		list.RPush(v)
	}

	return list
}

func (r *RedisListObject) Len() uint32 {
	return r.point().Len()
}

func (r *RedisListObject) LPush(value interface{}) {
	r.point().AddHead(value)
}

func (r *RedisListObject) RPush(value interface{}) {
	r.point().AddTail(value)
}

func (r *RedisListObject) LPop(value interface{}) interface{} {
	if r.point().Len() == 0 {
		return nil
	}
	data := r.point().First()
	r.point().DelHead()
	return data
}

func (r *RedisListObject) RPop(value interface{}) interface{} {
	if r.point().Len() == 0 {
		return nil
	}
	data := r.point().Last()
	r.point().DelTail()
	return data
}

func (r *RedisListObject) Index(idx uint32) interface{} {
	v := r.point().Index(idx)
	if v == nil {
		return nil
	}
	return v.Value()
}

func (r *RedisListObject) Insert(idx uint32, value interface{}) {
	r.point().Insert(idx, value)
}

func (r *RedisListObject) point() *basetype.LinkedList {
	return r.ptr.(*basetype.LinkedList)
}
