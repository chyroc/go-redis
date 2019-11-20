package objects

import (
	"github.com/chyroc/go-redis/internal/basetype"
)

type RedisZSetObject struct {
	*redisObjectImpl
}

func NewRedisZSetObject() *RedisZSetObject {
	// 只有满足下面两个条件的时候，才能用 intset
	// 元素个数小于 128
	// 每个元素字节长度小于 64
	return &RedisZSetObject{
		redisObjectImpl: &redisObjectImpl{
			type_:    RedisObjectTypeZSet,
			encoding: RedisObjectEncodingSkipList,
			ptr:      basetype.NewSkipList(),
		},
	}
}

// TODO
func (r *RedisZSetObject) ZAdd(data interface{}) {

}

func (r *RedisZSetObject) pinter() *basetype.Dict {
	return r.ptr.(*basetype.Dict)
}
