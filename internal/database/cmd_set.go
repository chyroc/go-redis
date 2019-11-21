package database

// 永久，-1
func Set(r *RedisDB, args ...string) (interface{}, error) {
	k := args[0]
	v := args[1]

	r.dict.Set(k, v)
	r.expires.Set(k, int64(-1))
	return "OK", nil
}
