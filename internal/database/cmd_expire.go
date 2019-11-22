package database

// TTL key
func Ttl(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]

	expire := r.expires.Get(k)

	if expire == nil {
		return -2, nil
	}
	if expire.(int64) == TimeNeverExpire {
		return TimeNeverExpire, nil
	}

	diff := expire.(int64) - nowMillisecond()
	if diff < 0 {
		r.dict.Del(k)
		r.expires.Del(k)
		return -2, nil
	}

	return diff / 1000, nil
}
