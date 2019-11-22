package database

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/basetype"
)

// HSET hash field value
// 将哈希表 hash 中域 field 的值设置为 value 。
//
// 返回值, 当 HSET 命令在哈希表中新创建 field 域并成功为它设置值时，命令返回 1
// 如果域 field 已经存在于哈希表， 并且 HSET 命令成功使用新值覆盖了它的旧值， 那么命令返回 0 。
func HSet(r *RedisDB, args ...string) (interface{}, error) {
	h := args[0]
	k := args[1]
	v := args[2]

	dict, value, err := r.getHashField(h, k)
	if err != nil {
		return nil, err
	}
	if dict == nil {
		dict = basetype.NewDict()
		r.setHash(h, dict)
	}
	dict.Set(k, v)

	if value == nil {
		return 1, nil
	} else {
		return 0, nil
	}
}

// HGET hash field
//
// 如果给定域不存在于哈希表中， 又或者给定的哈希表并不存在， 那么命令返回 nil 。
func HGet(r *RedisDB, args ...string) (interface{}, error) {
	h := args[0]
	k := args[1]

	_, value, err := r.getHashField(h, k)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (r *RedisDB) getHash(hash string) (*basetype.Dict, int64, error) {
	v := r.dict.Get(hash)
	if v == nil {
		return nil, TimeExpired, nil
	}

	vv, ok := v.(*basetype.Dict)
	if !ok {
		return nil, 0, fmt.Errorf("%v(%T): %w", v, v, ErrOperationWrongKindValue)
	}

	//expire := r.expires.Get(hash).(int64)
	//if expire == TimeNeverExpire {
	//	return vv, TimeNeverExpire, nil
	//} else if now := nowMillisecond(); now > expire {
	//	// 过期
	//	r.dict.Del(hash)
	//	r.expires.Del(hash)
	//	return nil, TimeExpired, nil
	//}

	return vv, TimeNeverExpire, nil
}

func (r *RedisDB) setHash(hash string, dict *basetype.Dict) {
	r.dict.Set(hash, dict)
}

func (r *RedisDB) getHashField(hash, field string) (*basetype.Dict, interface{}, error) {
	d, _, err := r.getHash(hash)
	if err != nil {
		return nil, nil, err
	}
	if d == nil {
		return nil, nil, nil
	}
	value := d.Get(field)
	if value == nil {
		return d, nil, nil
	}

	s, ok := value.(string)
	if !ok {
		return nil, nil, fmt.Errorf("%v(%T): %w", value, value, ErrOperationWrongKindValue)
	}

	return d, s, nil
}
