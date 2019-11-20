package objects

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/basetype"
	"github.com/chyroc/go-redis/internal/helper/ints"
	"strconv"
)

type RedisSetObject struct {
	*redisObjectImpl
}

func NewRedisSetObject() *RedisSetObject {
	// 只有满足下面两个条件的时候，才能用 intset
	// 所有元素都是整数
	// 元素个数不超过 512
	list := &RedisSetObject{
		redisObjectImpl: &redisObjectImpl{
			type_:    RedisObjectTypeSet,
			encoding: RedisObjectEncodingIntSet,
			ptr:      basetype.NewIntset(),
		},
	}
	return list
}

func (r *RedisSetObject) SAdd(data interface{}) {
	r.upgrade(data)

	if r.encoding == RedisObjectEncodingIntSet {
		r.intsetPpinter().Add(data)
	} else {
		r.dictPointer().Add(data)
	}
}

func (r *RedisSetObject) SCard() uint32 {
	if r.encoding == RedisObjectEncodingIntSet {
		return r.intsetPpinter().Len()
	} else {
		return r.dictPointer().Len()
	}
}

// TODO
func (r *RedisSetObject) SIsMember(value interface{}) {
	if r.encoding == RedisObjectEncodingIntSet {

	} else {

	}
}

// TODO
func (r *RedisSetObject) SMembers() {
	if r.encoding == RedisObjectEncodingIntSet {

	} else {

	}
}

// TODO
func (r *RedisSetObject) SRandMember() {
	if r.encoding == RedisObjectEncodingIntSet {

	} else {

	}
}

// TODO
func (r *RedisSetObject) SPop() {
	if r.encoding == RedisObjectEncodingIntSet {

	} else {

	}
}

// TODO
func (r *RedisSetObject) SRem() {
	if r.encoding == RedisObjectEncodingIntSet {

	} else {

	}
}

func (r *RedisSetObject) intsetPpinter() *basetype.Intset {
	return r.ptr.(*basetype.Intset)
}

func (r *RedisSetObject) dictPointer() *basetype.Intset {
	return r.ptr.(*basetype.Intset)
}

// 升级
func (r *RedisSetObject) upgrade(data interface{}) {
	if !r.needUpgrade(data) {
		return
	}

	intset := r.intsetPpinter()
	dict := basetype.NewDict()
	intset.Range(func(idx uint32, data interface{}) (continue_ bool) {
		switch v := data.(type) {
		case int16:
			dict.Set(strconv.FormatInt(int64(v), 10), nil)
		case int32:
			dict.Set(strconv.FormatInt(int64(v), 10), nil)
		case int64:
			dict.Set(strconv.FormatInt(v, 10), nil)
		default:
			panic(fmt.Sprintf("intset 获取了一个不是 int16,32,64 的数据: %v(%T)", data, data))
		}
		return true
	})
	r.encoding = RedisObjectEncodingHashTable
	r.ptr = dict
}

func (r *RedisSetObject) needUpgrade(data interface{}) bool {
	// 只有满足下面两个条件的时候，才能用 intset
	// 所有元素都是整数
	// 元素个数不超过 512

	if r.encoding == RedisObjectEncodingHashTable {
		return false
	}

	if r.intsetPpinter().Len()+1 > 512 {
		return true
	}

	if !ints.IsInt(data) {
		return true
	}

	return false
}
