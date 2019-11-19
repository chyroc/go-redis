package objects

import (
	"github.com/chyroc/go-redis/internal/basetype"
	"strconv"
)

// 字符串对象

type RedisStringObject struct {
	*redisObjectImpl
}

// 字符串对象
//
// RedisObjectEncodingInt + int64 数字，
// RedisObjectEncodingRaw + *SDS 字符串
// RedisObjectEncodingEmbStr + *SDS
//
// redis 中实现：RedisObjectEncodingEmbStr 对应的是只分配一次内存，这里不搞
func NewRedisStringObject(s string) *RedisStringObject {
	var encoding RedisObjectEncoding
	var data interface{}

	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		encoding = RedisObjectEncodingInt
		data = i
	} else if len(s) > 32 {
		encoding = RedisObjectEncodingRaw
		data = basetype.NewSDSWithString(s)
	} else {
		encoding = RedisObjectEncodingEmbStr
		data = basetype.NewSDSWithString(s)
	}

	return &RedisStringObject{
		redisObjectImpl: &redisObjectImpl{
			type_:    RedisObjectTypeString,
			encoding: encoding,
			ptr:      data,
		},
	}
}
