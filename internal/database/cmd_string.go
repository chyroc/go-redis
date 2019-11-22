package database

import (
	"fmt"
)

const TimeNeverExpire int64 = -1 // 永久，-1
const TimeExpired int64 = -2     // 已过期、不存在

func stringOrError(v interface{}) (interface{}, error) {
	if v == nil {
		return v, nil
	}
	switch v.(type) {
	case string:
		return v, nil
	default:
		return nil, ErrOperationWrongKindValue
	}
}

func Get(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]

	v := r.dict.Get(k)
	if v == nil {
		return nil, nil
	}
	if _, err := stringOrError(v); err != nil {
		return nil, err
	}

	expire := r.expires.Get(k).(int64)
	if expire == TimeNeverExpire {
		return v, nil
	}

	if now := nowMillisecond(); now > expire {
		// 过期
		r.dict.Del(k)
		r.expires.Del(k)
		return nil, nil
	}

	return v, nil
}

// SET key value [EX seconds] [PX milliseconds] [NX|XX]
func Set(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]
	v := args[1]

	offset, millisecond, err := getMillisecond(args, 2)
	if err != nil {
		return nil, err
	}
	fmt.Println("millisecond", millisecond)
	if millisecond == 0 {
		millisecond = TimeNeverExpire
	} else {
		millisecond += nowMillisecond()
	}
	nx, xx, err := getNxXx(args, offset)
	if err != nil {
		return nil, err
	}

	v2 := r.dict.Get(k)
	if nx && v2 != nil {
		return nil, nil
	}
	if xx && v2 == nil {
		return nil, nil
	}

	r.dict.Set(k, v)
	r.expires.Set(k, millisecond)
	return status("OK"), nil
}

// GETSET key value
func GetSet(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]
	v := args[1]

	if _, err := stringOrError(v); err != nil {
		return nil, err
	}

	old, err := Get(r, k)
	if err != nil {
		return nil, err
	}
	//stringOrError(old)

	if _, err := Set(r, k, v); err != nil {
		return nil, err
	}

	return old, nil
}

// STRLEN key
// 返回键 key 储存的字符串值的长度。
// 当键 key 不存在时， 命令返回 0 。
// 当 key 储存的不是字符串值时， 返回一个错误。
func StrLen(r *RedisDB, args ...string) (interface{}, error) {
	v, err := Get(r, args...)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return 0, nil
	}
	if _, err := stringOrError(v); err != nil {
		return nil, err
	}
	return len(v.(string)), nil
}
