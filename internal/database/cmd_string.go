package database

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/basetype"
	"strconv"
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

// SETRANGE key offset value
// 从偏移量 offset 开始， 用 value 参数覆写(overwrite)键 key 储存的字符串值。
// 不存在的键 key 当作空白字符串处理。
// SETRANGE 命令会确保字符串足够长以便将 value 设置到指定的偏移量上
// 如果键 key 原来储存的字符串长度比偏移量小，那么原字符和偏移量之间的空白将用零字节(zerobytes, "\x00" )进行填充。
// 因为 Redis 字符串的大小被限制在 512 兆(megabytes)以内， 所以用户能够使用的最大偏移量为 2^29-1(536870911) ， 如果你需要使用比这更大的空间， 请使用多个 key 。

// INCR key
// 为键 key 储存的数字值加上一。
// 如果键 key 不存在， 那么它的值会先被初始化为 0 ， 然后再执行 INCR 命令。
// 如果键 key 储存的值不能被解释为数字， 那么 INCR 命令将返回一个错误。
// 返回值：INCR 命令会返回键 key 在执行加一操作之后的值。
func Incr(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]

	v, _, err := r.getSDS(k)
	if err != nil {
		return nil, err
	}
	if v == nil {
		if _, err := r.setSDS(k, strconv.Itoa(1), nil, 0, false, false); err != nil {
			return nil, err
		}
		return 1, nil
	}

	return v.Int64Incr()
}
