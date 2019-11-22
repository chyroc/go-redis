package database

import (
	"fmt"
	"github.com/chyroc/go-pointer"
	"strconv"
	"strings"
	"time"
)

const TimeNeverExpire int64 = -1 // 永久，-1
const TimeExpired int64 = -2     // 已过期、不存在

func Get(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]

	v := r.dict.Get(k)
	if v == nil {
		return nil, nil
	}

	expire := r.expires.Get(k).(int64)
	if expire == TimeNeverExpire {
		return pointer.String(v.(string)), nil
	}

	if now := nowMillisecond(); now > expire {
		// 过期
		r.dict.Del(k)
		r.expires.Del(k)
		return nil, nil
	}

	return pointer.String(v.(string)), nil
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
	return "OK", nil
}

func getMillisecond(args []string, offset int) (off int, ms int64, err error) {
	if len(args) < offset+1 {
		return offset, 0, nil
	}

	r := strings.ToLower(args[offset])
	if r != "ex" && r != "px" {
		return offset, 0, nil
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
		return offset + 2, 1000 * i, nil
	}
	// 毫秒
	return offset + 2, i, nil
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
