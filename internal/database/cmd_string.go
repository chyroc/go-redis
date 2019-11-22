package database

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/basetype"
)

const TimeNeverExpire int64 = -1 // 永久，-1
const TimeExpired int64 = -2     // 已过期、不存在

func sdsStringOrError(v interface{}) (interface{}, error) {
	if v == nil {
		return v, nil
	}
	switch v.(type) {
	case *basetype.SDS:
		return v, nil
	default:
		return nil, fmt.Errorf("%v(%T): %w", v, v, ErrOperationWrongKindValue)
	}
}

func Get(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]

	v, _, err := r.getSDS(k)
	return v, err
}

// SET key value [EX seconds] [PX milliseconds] [NX|XX]
func Set(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]
	v := args[1]

	v2, _, err := r.getSDS(k)
	if err != nil {
		return nil, err
	}

	offset, millisecond, err := getMillisecond(args, 2)
	if err != nil {
		return nil, err
	}
	nx, xx, err := getNxXx(args, offset)
	if err != nil {
		return nil, err
	}

	return r.setSDS(k, v, v2, millisecond, nx, xx)
}

// GETSET key value
func GetSet(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]
	v := args[1]

	v2, _, err := r.getSDS(k)
	if err != nil {
		return nil, err
	}

	if _, err := r.setSDS(k, v, nil, TimeNeverExpire, false, false); err != nil {
		return nil, err
	}

	return v2, nil
}

// STRLEN key
// 返回键 key 储存的字符串值的长度。
// 当键 key 不存在时， 命令返回 0 。
// 当 key 储存的不是字符串值时， 返回一个错误。
func StrLen(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]

	v, _, err := r.getSDS(k)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return 0, nil
	}
	return v.Len(), nil
}

// APPEND key value
// 如果键 key 已经存在并且它的值是一个字符串， APPEND 命令将把 value 追加到键 key 现有值的末尾。
// 如果 key 不存在， APPEND 就简单地将键 key 的值设为 value ， 就像执行 SET key value 一样。
// 返回：追加 value 之后， 键 key 的值的长度。
func Append(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]
	v := args[1]

	v2, _, err := r.getSDS(k)
	if err != nil {
		return nil, err
	}
	if v2 == nil {
		if _, err := Set(r, k, v); err != nil {
			return nil, err
		}
		return len(v), nil
	}

	v2.Append(v)

	return v2.Len(), nil
}
