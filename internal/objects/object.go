package objects

// Redis 对象

// 只有字符串对象会被其他对象嵌套使用，其他对象不允许被嵌套

type RedisObject interface {
	Type() RedisObjectType
	Encoding() RedisObjectEncoding
}

type redisObjectImpl struct {
	type_    RedisObjectType     // 类型
	encoding RedisObjectEncoding // 编码
	ptr      interface{}         // 类型和编码，决定了 ptr 的解析方式
}

func (r *redisObjectImpl) Type() RedisObjectType {
	return r.type_
}

func (r *redisObjectImpl) Encoding() RedisObjectEncoding {
	return r.encoding
}
