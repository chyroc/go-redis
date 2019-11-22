package database

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chyroc/go-redis/internal/basetype"
)

func (r *RedisDB) getSDS(k string) (*basetype.SDS, int64, error) {
	v := r.dict.Get(k)
	if v == nil {
		return nil, TimeExpired, nil
	}

	vv, ok := v.(*basetype.SDS)
	if !ok {
		return nil, 0, fmt.Errorf("%v(%T): %w", v, v, ErrOperationWrongKindValue)
	}

	expire := r.expires.Get(k).(int64)
	if expire == TimeNeverExpire {
		return vv, TimeNeverExpire, nil
	} else if now := nowMillisecond(); now > expire {
		// 过期
		r.dict.Del(k)
		r.expires.Del(k)
		return nil, TimeExpired, nil
	}

	return vv, expire, nil
}

func (r *RedisDB) setSDS(k, v string, old interface{}, millisecond int64, nx bool, xx bool) interface{} {
	if nx && old != nil {
		return nil
	}
	if xx && old == nil {
		return nil
	}

	r.dict.Set(k, basetype.NewSDSWithString(v))
	r.expires.Set(k, millisecond)
	return status("OK")
}

func getMillisecond(args []string, offset int) (off int, ms int64, err error) {
	if len(args) < offset+1 {
		return offset, TimeNeverExpire, nil
	}

	r := strings.ToLower(args[offset])
	if r != "ex" && r != "px" {
		return offset, TimeNeverExpire, nil
	}
	if len(args) < offset {
		return 0, 0, fmt.Errorf("[Redis.Get] got params %q, but no seconds params", args[offset])
	}

	i, err := strconv.ParseInt(args[offset+1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("[Redis.Get] got params %q, but next param %q is not int", args[offset], args[offset+1])
	}

	if r == "ex" {
		// 秒
		return offset + 2, nowMillisecond() + 1000*i, nil
	}
	// 毫秒
	return offset + 2, nowMillisecond() + i, nil
}

func getNxXx(args []string, offset int) (nx bool, xx bool, err error) {
	if len(args) < offset+1 {
		return
	}
	r := strings.ToLower(args[offset])
	if r != "nx" && r != "xx" {
		err = fmt.Errorf("[Redis.Get] got params %q, but need %q or %q", args[offset], "NX", "XX")
		return
	}

	if len(args) > offset+1 {
		err = fmt.Errorf("[Redis.Get] got params %q, which endswith %q, but got extra params", args, args[offset])
		return
	}

	return r == "nx", r == "xx", nil
}

func nowMillisecond() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}
