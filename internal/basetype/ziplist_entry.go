package basetype
//
//import (
//	"encoding/binary"
//	"fmt"
//	"github.com/chyroc/go-redis/internal/helper"
//	gr_binary "github.com/chyroc/go-redis/internal/helper/gr-binary"
//)
//
////offset uint32 // TODO: 不知道要不要放在这里
//type entry struct {
//	prev     []byte // 表示【上个 entry 长度】的字节数组
//	encoding []byte // 编码
//	data     []byte // 数据
//}
//
//func newEntry(prev uint32, data interface{}) (e *entry) {
//	e = new(entry)
//
//	// prev
//	e.prev = e.prevToBytes(prev)
//
//	// encoding
//	e.encoding = entryGetEncoding(data)
//
//	// data
//	e.data = entryGetData(data)
//
//	return
//}
//
//func (e *entry) prevBytesLength() uint32 {
//	return uint32(len(e.prev))
//}
//
//func (e *entry) prevEntryLength() uint32 {
//	if e.prev[0] < 0xFE {
//		return uint32(e.prev[0])
//	} else {
//		return binary.LittleEndian.Uint32(e.prev[1:5])
//	}
//}
//
//func (e *entry) totalLength() uint32 {
//	return e.prevBytesLength() + e.encodingLen() + e.bytesLen()
//}
//
//func (e *entry) int() int64 {
//	switch e.encoding[0] {
//	case 0b11000000:
//		// 11000000，int16，额外 2 字节
//		return int64(binary.LittleEndian.Uint16((e.data)[:2]))
//	case 0b11010000:
//		// 11010000，int32，额外 4 字节
//		return int64(binary.LittleEndian.Uint32((e.data)[:4]))
//	case 0b11100000:
//		// 11100000，int64，额外 8 字节
//		return int64(binary.LittleEndian.Uint64((e.data)[:8]))
//	case 0b11110000:
//		// 11110000，int24，额外 3 字节，24 位有符号整数
//		return int64(gr_binary.LittleEndian.Uint24((e.data)[:3]))
//	case 0x11111110:
//		// 11111110，int8，额外 1 字节，08 位有符号整数
//		return int64(gr_binary.LittleEndian.Uint8((e.data)[:1]))
//	default:
//		// 1111xxxx，int4，没有额外的字节，0-12
//		if helper.GetBit(e.encoding[0], 2) == 1 && helper.GetBit(e.encoding[0], 3) == 1 {
//			return int64(e.encoding[0] - 0b11110000 - 1)
//		}
//	}
//	panic(fmt.Sprintf("求 entry.int 的时候出现了不允许的情况"))
//}
//
//func (e *entry) encodingLen() uint32 {
//	return uint32(len(e.encoding))
//}
//
//func (e *entry) dataLen() uint32 {
//	return uint32(len(e.data))
//}
//
//func (e *entry) bytesLen() uint32 {
//	encoding0 := helper.GetBit(e.encoding[0], 0)
//	encoding1 := helper.GetBit(e.encoding[0], 1)
//	switch encoding0 {
//	case 0:
//		switch encoding1 {
//		case 0:
//			// 00：encoding 为 1 个字节，剩下的 6 个 bit 存储了一个数字，数字的长度表示 e 存储的 bytes 的长度
//			x := uint32(e.encoding[0] - 0b00000000)
//			return x
//		case 1:
//			// 01：encoding 为 2 个字节，剩下的 6+8 个 bit 存储了一个数字，数字的长度表示 e 存储的 bytes 的长度
//			// ++++ 这里是大端序 ++++
//			x := uint32(binary.BigEndian.Uint16([]byte{
//				e.encoding[0] - 0b01000000, // 减去开头的 01
//				e.encoding[1],              // 下个字节
//			}))
//			return x
//		}
//	case 1:
//		switch encoding1 {
//		case 0:
//			// 10：encoding 为 5 个字节，剩下的 4*8 个 bit 存储了一个数字，数字的长度表示 e 存储的 bytes 的长度
//			// 这里是大端序
//			x := uint32(binary.BigEndian.Uint32(e.encoding[1:5]))
//			return x
//		}
//	}
//
//	panic(fmt.Sprintf("求 entry.bytes 的时候出现了不允许的情况"))
//}
//
//func (e *entry) bytes() []byte {
//	return (e.data)[:e.bytesLen()]
//}
//
//func (e *entry) prevToBytes(i uint32) []byte {
//	if i < 0xFE {
//		return []byte{
//			byte(i),
//		}
//	} else {
//		bs := make([]byte, 5)
//		bs[0] = 0xFE
//		binary.LittleEndian.PutUint32(bs[1:5], uint32(i))
//		return bs
//	}
//}
//
//func (e *entry) allbytes() []byte {
//	bs := make([]byte, 0, e.prevBytesLength()+e.encodingLen()+e.dataLen())
//	bs = append(bs, e.prev...)
//	bs = append(bs, e.encoding...)
//	bs = append(bs, e.data...)
//	return bs
//}
