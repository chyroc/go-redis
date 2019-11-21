package database

import (
	"github.com/chyroc/go-pointer"
	"time"
)

func (r *RedisDB) Get(k string) (*string, error) {
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

// 永久，-1
func (r *RedisDB) Set(k, v string) ( error) {
	r.dict.Set(k, v)
	r.expires.Set(k, int64(-1))
	return nil
}
