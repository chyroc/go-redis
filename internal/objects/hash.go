package objects

import (
	"github.com/chyroc/go-ptr"
	"github.com/chyroc/go-redis/internal/basetype"
)

type RedisHashObject struct {
	*redisObjectImpl
}

func NewRedisHashObject() *RedisHashObject {
	// 只有满足下面两个条件的时候，才能用 ziplist
	// 所有 key value 的长度都小于 64 字节
	// 元素个数小于 512

	list := &RedisHashObject{
		redisObjectImpl: &redisObjectImpl{
			type_:    RedisObjectTypeList,
			encoding: RedisObjectEncodingHashTable, // TODO ziplist
			ptr:      basetype.NewDict(),
		},
	}
	return list
}

func (r *RedisHashObject) point() *basetype.Dict {
	return r.ptr.(*basetype.Dict)
}

func (r *RedisHashObject) HSet(k, v string) {
	vv := basetype.NewSDSWithString(v)
	r.point().Set(k, vv)
}

func (r *RedisHashObject) HGet(k string) *string {
	v := r.point().Get(k)
	if v == nil {
		return nil
	}
	return ptr.String(v.(*basetype.SDS).String())
}

func (r *RedisHashObject) HExist(k string) bool {
	v := r.point().Get(k)
	return v != nil
}

func (r *RedisHashObject) HDel(k string) {
	r.point().Del(k)
}

func (r *RedisHashObject) HLen() uint32 {
	return r.point().Size()
}

// TODO: 不知道数据结构
func (r *RedisHashObject) HGetAll() {
}
