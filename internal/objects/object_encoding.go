package objects

type RedisObjectEncoding uint32

const (
	RedisObjectEncodingRaw RedisObjectEncoding = iota
	RedisObjectEncodingInt
	RedisObjectEncodingHashTable
	RedisObjectEncodingZipMap
	RedisObjectEncodingLinkedList
	RedisObjectEncodingZipList
	RedisObjectEncodingIntSet
	RedisObjectEncodingSkipList
	RedisObjectEncodingEmbStr
	RedisObjectEncodingQuickList
	RedisObjectEncodingStream
)

//#define OBJ_ENCODING_RAW 0     /* Raw representation */
//#define OBJ_ENCODING_INT 1     /* Encoded as integer */
//#define OBJ_ENCODING_HT 2      /* Encoded as hash table */
//#define OBJ_ENCODING_ZIPMAP 3  /* Encoded as zipmap */
//#define OBJ_ENCODING_LINKEDLIST 4 /* No longer used: old list encoding. */
//#define OBJ_ENCODING_ZIPLIST 5 /* Encoded as ziplist */
//#define OBJ_ENCODING_INTSET 6  /* Encoded as intset */
//#define OBJ_ENCODING_SKIPLIST 7  /* Encoded as skiplist */
//#define OBJ_ENCODING_EMBSTR 8  /* Embedded sds string encoding */
//#define OBJ_ENCODING_QUICKLIST 9 /* Encoded as linked list of ziplists */
//#define OBJ_ENCODING_STREAM 10 /* Encoded as a radix tree of listpacks */

func (r RedisObjectEncoding) String() string {
	if r <= 10 {
		return []string{
			"raw",
			"int",
			"hashtable",
			"zipmap",
			"linkedlist",
			"ziplist",
			"intset",
			"skiplist",
			"embstr",
			"quicklist",
			"stream",
		}[r]
	}
	return "unknown"
}
