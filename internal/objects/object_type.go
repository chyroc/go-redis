package objects

type RedisObjectType uint32

const (
	RedisObjectTypeString RedisObjectType = iota
	RedisObjectTypeList
	RedisObjectTypeHash
	RedisObjectTypeSet
	RedisObjectTypeZSet
)

//#define OBJ_STRING 0    /* String object. */
//#define OBJ_LIST 1      /* List object. */
//#define OBJ_SET 2       /* Set object. */
//#define OBJ_ZSET 3      /* Sorted set object. */
//#define OBJ_HASH 4      /* Hash object. */

func (r RedisObjectType) String() string {
	return []string{
		"string",
		"list",
		"hash",
		"set",
		"zset",
	}[r]
}
