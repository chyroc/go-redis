package basetype
//
//import (
//	"github.com/chyroc/go-redis/internal/helper"
//	"math"
//
//	gr_binary "github.com/chyroc/go-redis/internal/helper/gr-binary"
//)
//
//func entryGetData(s interface{}) []byte {
//	switch i := s.(type) {
//	case int:
//		return entryGetData1(int64(i))
//	case int8:
//		return entryGetData1(int64(i))
//	case int16:
//		return entryGetData1(int64(i))
//	case int32:
//		return entryGetData1(int64(i))
//	case int64:
//		return entryGetData1(int64(i))
//	case []byte:
//		return i
//	}
//
//	panic("entryGetData")
//}
//
//// int
//func entryGetData1(i int64) (bs []byte) {
//	if 0 <= i && i <= 12 {
//		// 1111xxxx，int32，没有额外的字节，0-12
//		return []byte{}
//	}
//	if math.MinInt8 <= i && i <= math.MaxInt8 {
//		// 11111110，int8，额外 1 字节，8 位有符号整数
//		return []byte{byte(i)}
//	}
//	if math.MinInt16 <= i && i <= math.MaxInt16 {
//		// 11000000，int16，额外 2 字节
//		return helper.Int16ToBinary(int16(i))
//	}
//	if gr_binary.MinInt24 <= i && i <= gr_binary.MaxInt24 {
//		// 11110000，int24，额外 3 字节，24 位有符号整数
//		return helper.Int24ToBinary(int32(i))
//	}
//	if math.MinInt32 <= i && i <= math.MaxInt32 {
//		// 11010000，int32，额外 4 字节
//		return helper.Int32ToBinary(int32(i))
//	}
//
//	// 11100000，int64，额外 8 字节
//	return helper.Int64ToBinary(int64(i))
//}
