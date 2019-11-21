package database

import (
	"github.com/chyroc/go-pointer"
	"time"
)

func Get(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]

	v := r.dict.Get(k)
	if v == nil {
		return nil, nil
	}

	expire := r.expires.Get(k).(int64)
	if expire == -1 {
		return pointer.String(v.(string)), nil
	}

	if now := time.Now().UnixNano(); now > expire {
		// 过期
		r.dict.Del(k)
		return nil, nil
	}

	return pointer.String(v.(string)), nil
}
