package basetype

//
//import (
//	"encoding/binary"
//	"fmt"
//	"github.com/chyroc/go-redis/internal/helper"
//	"math"
//
//	gr_binary "github.com/chyroc/go-redis/internal/helper/gr-binary"
//)
//
//func entryEncodingToint(encoding []byte) int64 {
//	switch encoding[0] {
//	case 0b11000000:
//		// 11000000，int16，额外 2 字节
//		return int64(binary.LittleEndian.Uint16((e.ptr)[:2]))
//	case 0b11010000:
//		// 11010000，int32，额外 4 字节
//		return int64(binary.LittleEndian.Uint32((e.ptr)[:4]))
//	case 0b11100000:
//		// 11100000，int64，额外 8 字节
//		return int64(binary.LittleEndian.Uint64((e.ptr)[:8]))
//	case 0b11110000:
//		// 11110000，int24，额外 3 字节，24 位有符号整数
//		return int64(gr_binary.LittleEndian.Uint24((e.ptr)[:3]))
//	case 0x11111110:
//		// 11111110，int8，额外 1 字节，08 位有符号整数
//		return int64(gr_binary.LittleEndian.Uint8((e.ptr)[:1]))
//	default:
//		// 1111xxxx，int4，没有额外的字节，0-12
//		if helper.GetBit(e.encoding[0], 2) == 1 && helper.GetBit(e.encoding[0], 3) == 1 {
//			return int64(e.encoding[0] - 0b11110000 - 1)
//		}
//	}
//	panic(fmt.Sprintf("求 entry.int 的时候出现了不允许的情况"))
//}
//
//func entryGetEncoding(s interface{}) []byte {
//	switch i := s.(type) {
//	case int:
//		return entryGetEncoding1(int64(i))
//	case int8:
//		return entryGetEncoding1(int64(i))
//	case int16:
//		return entryGetEncoding1(int64(i))
//	case int32:
//		return entryGetEncoding1(int64(i))
//	case int64:
//		return entryGetEncoding1(int64(i))
//	case []byte:
//		return entryGetEncoding2(i)
//	}
//
//	panic("getEncoding")
//}
//
//// int
//func entryGetEncoding1(i int64) []byte {
//	if 0 <= i && i <= 12 {
//		// 1111xxxx，int32，没有额外的字节，0-12
//		return []byte{12 + 1 + 0b11110000}
//	}
//	if math.MinInt8 <= i && i <= math.MaxInt8 {
//		// 11111110，int8，额外 1 字节，8 位有符号整数
//		return []byte{0b11111110}
//	}
//	if math.MinInt16 <= i && i <= math.MaxInt16 {
//		// 11000000，int16，额外 2 字节
//		return []byte{11000000}
//	}
//	if gr_binary.MinInt24 <= i && i <= gr_binary.MaxInt24 {
//		// 11110000，int24，额外 3 字节，24 位有符号整数
//		return []byte{11110000}
//	}
//	if math.MinInt32 <= i && i <= math.MaxInt32 {
//		// 11010000，int32，额外 4 字节
//		return []byte{11110000}
//	}
//
//	// 11100000，int64，额外 8 字节
//	return []byte{11100000}
//}
//
//// []byte encoding
//func entryGetEncoding2(i []byte) []byte {
//	ilen := uint32(len(i))
//	if ilen <= 0b00111111 {
//		return []byte{byte(ilen)}
//	} else if ilen < 0b01111111 {
//		bs := make([]byte, 2)
//		binary.BigEndian.PutUint16(bs, uint16(ilen))
//		bs[0] += 0b01000000
//		return bs
//	} else if ilen < 0b10000000 {
//		bs := make([]byte, 5)
//		binary.BigEndian.PutUint32(bs[1:5], uint32(ilen))
//		bs[0] = 0b10000000
//		return bs
//	} else {
//		panic(fmt.Sprintf("getEncoding2"))
//	}
//}
