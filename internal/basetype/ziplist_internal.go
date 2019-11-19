package basetype
//
//import (
//	"fmt"
//	"github.com/chyroc/go-redis/internal/helper"
//	gr_binary "github.com/chyroc/go-redis/internal/helper/gr-binary"
//	"math"
//)
//
//// 抄一下 redis 的文档
//// https://github.com/antirez/redis/blob/unstable/src/ziplist.c#L1-L180
//
//// 除了两个特殊情况，ziplist 都是小端存储的
////   - |01pppppp|qqqqqqqq| - 2 bytes：14 bit number is stored in big endian
////   - |10000000|qqqqqqqq|rrrrrrrr|ssssssss|tttttttt| - 5 bytes：32 bit number is stored in big endian.
//
///*
//
// * ZIPLIST OVERALL LAYOUT
// * ======================
// *
// * <zlbytes> <zltail> <zllen> <entry> <entry> ... <entry> <zlend>
// *
// * 小端序存储
// *
// * <zlbytes> 4 字节，整个 ziplist 的 bytes 个数
// *
// * <zltail> 4 字节，最后一个 entry 在 ziplist 中的偏移量
// *
// * <zllen> 2 字节， entry 的个数。当整个值等于 2^16-1 的是，表示个数大于 2 个字节能够存储的上限，需要遍历列表
// *
// * <zlend> 1 字节，为 0xFF 固定值，表示 ziplist 结束
// *
// * ZIPLIST ENTRIES
// * ===============
// *
// * <prevlen> <encoding> <entry-data>
// * 如果 data 是一个 0-12 之间的数字，就会直接编码到 encoding 中
// *
// * <prevlen> 上一个 entry 的长度
// * 如果小于 254，就 1 个字节，uint8
// * 否则，是 5 字节，并且第一个字节是 FE(254)，下面的4个字节表示长度 uint32；【所以先取第一个字节，判断接下来的动作】
// *
// * <encoding>
// * 如果 encoding 的前两个 bit 为 `00`，`01`，`10`，那么 content 是一个 bytes
//     * 00：encoding 为 1 个字节，剩下的 6 个 bit 存储了一个数字，数字的长度表示 entry 存储的 bytes 的长度
//     * 01：encoding 为 2 个字节，剩下的 6+8 个 bit 存储了一个数字，数字的长度表示 entry 存储的 bytes 的长度
//     * 10：encoding 为 5 个字节，剩下的 4*8 个 bit 存储了一个数字，数字的长度表示 entry 存储的 bytes 的长度
// * 11 开头，存了一个数字
//     * 11000000，int16，额外 2 字节
//     * 11010000，int32，额外 4 字节
//     * 11100000，int64，额外 8 字节
//     * 11110000，int24，额外 3 字节，24 位有符号整数
//     * 11111110，int8，额外 1 字节，08 位有符号整数
//     * 1111xxxx，int32，没有额外的字节，0-12
//       * xxxx between 0000 and 1101
//       * 0000 是 0，但是不能用
//       * 1110 是 13，但是不能用
//       * 所以范围是 1-13
//       * 然后减掉1，所以范围是 0-12
// *
// * EXAMPLES OF ACTUAL ZIPLISTS
// * ===========================
// *
// * 下面是一个存储了 2 和 5 的 ziplist
// *
// *  [0f 00 00 00] [0c 00 00 00] [02 00] [00 f3] [02 f6] [ff]
// *        |             |          |       |       |     |
// *     zlbytes        zltail    entries   "2"     "5"   end
// *
// * 0c = 12
// * f3 = 1111 0011 = 3，然后减-1，等于 2
// * f6 = 1111 0110 = 6，然后减-1，等于 5
// *
//*/
//
//// 4 + 4 + 2 + 1 = 11 ,就算没有数据，也是占用 11 个字节的数据
//// https://www.binaryhexconverter.com/hex-to-binary-converter
//
////type ZipList struct {
////	zlbytes uint32 // 记录 ziplist 占用字节大小
////	zltail  uint32 // 记录 ziplist 的最后一个节点的起始地址，和 ziplist 的起始地址的差值
////	zllen   uint16 // 节点数量
////	entry   []byte // 数据存储
////	zlend   uint8  // 0xFF 表示结束
////}
//
////
//
//// 扩容
////
//// 如果容量够，啥都不做
//// 如果扩容后的容量大于 1024，就扩容为 1024的整数倍
//// 否则，扩容为新容量的 2 倍
//func (r *ZipList) expansion(length uint32) {
//	// 判断是否真的需要扩容
//	total := r.BytesLength() + length
//	if total < uint32(len(*r)) {
//		return
//	}
//
//	// 计算新的容量
//	var l uint32
//	if total > 1024 {
//		l = 1024 * (total/1024 + 1)
//	} else {
//		l = 2 * total
//	}
//
//	// 迁移数据
//	buf := make([]byte, l)
//	copy(buf, (*r)[:r.BytesLength()])
//	*r = buf
//}
//
////// 添加数据到尾部
////func (r *ZipList) appendData(encodingAndData []byte, length int32) {
////	// append pre_length
////	bs := r.getOffsetEntryLengthBytes(uint32(r.Tail())) // 前个数据长度
////	bs = append(bs, encodingAndData[:length]...)        // encoding 和 数据
////	bs = append(bs, 0xff)                               // 0xff 结尾
////
////	total := int(r.BytesLength())
////
////	copy((*r)[total-1:total-1+len(bs)], bs) // 最后一个节点的数据 + 0xff
////
////	r.setTail(r.BytesLength() - 1)
////	r.setBytesLength(r.BytesLength() + uint32(len(bs)))
////	r.setLen(r.Len() + 1)
////}
//
////func (r *ZipList) getOffsetEncoding(offset int64) [2]byte {
////	b := (*r)[offset]
////	encoding_0 := (b >> uint(7)) & 1
////	encoding_1 := (b >> uint(6)) & 1
////	return [2]byte{encoding_0, encoding_1}
////	//	 * 如果 encoding 的前两个 bit 为 `00`，`01`，`10`，那么 content 是一个 bytes
////	//     * 00：encoding 为 1 个字节，剩下的 6 个 bit 存储了一个数字，数字的长度表示 entry 存储的 bytes 的长度
////	//     * 01：encoding 为 2 个字节，剩下的 6+8 个 bit 存储了一个数字，数字的长度表示 entry 存储的 bytes 的长度
////	//     * 10：encoding 为 5 个字节，剩下的 4*8 个 bit 存储了一个数字，数字的长度表示 entry 存储的 bytes 的长度
////	// * 11 开头，存了一个数字
////}
//
////// 获取 idx 节点的长度
////func (r *ZipList) getEntryLengthByIdx(idx uint32) uint32 {
////	entryOffset := r.getEntryOffsetByIdx(idx)
////	return r.getEntryLengthByOffset(entryOffset)
////}
//
//// 获取 idx 节点，所在的 offset
//func (r *ZipList) getEntryOffsetByIdx(idx uint32) uint32 {
//	if idx == 0 {
//		return 10
//	}
//
//	length := r.Len()
//	if idx < 0 || idx >= length {
//		panic(fmt.Sprintf("指定的节点ID(%d) 大于总长(%d)", idx, length))
//	}
//
//	// [0, n-1)
//
//	var entryOffset uint32
//	if idx < length/2 {
//		// 从前往后
//		var i uint32
//		entryOffset = uint32(10)
//		for i = 0; i < idx; i++ {
//			entry := r.getEntryByOffset(entryOffset)
//			entryOffset += entry.totalLength()
//		}
//	} else {
//		// 从后往前
//		var i uint32
//		entryOffset = r.Tail()
//		for i = 0; i < idx; i++ {
//			entry := r.getEntryByOffset(entryOffset)
//			entryOffset -= entry.prevEntryLength()
//		}
//	}
//	// entryOffset 就是 idx 对应的 offset
//	return entryOffset
//}
//
////
////// 获取 idx 节点的长度
////func (r *ZipList) getEntryLengthByIdx(idx uint32) uint32 {
////	entryOffset := r.getEntryOffsetByIdx(idx)
////	return r.getEntryLengthByOffset(entryOffset)
////}
////
////// 获取 idx 节点，所在的 offset
////func (r *ZipList) getEntryOffsetByIdx(idx uint32) uint32 {
////	if idx == 0 {
////		return 10
////	}
////
////	length := r.Len()
////	if idx < 0 || idx >= length {
////		panic(fmt.Sprintf("指定的节点ID(%d) 大于总长(%d)", idx, length))
////	}
////
////	// [0, n-1)
////
////	var entryOffset uint32
////	if idx < length/2 {
////		// 从前往后
////		var i uint32
////		entryOffset = uint32(10)
////		for i = 0; i < idx; i++ {
////			entryLength := uint32(r.getEntryLengthByOffset(entryOffset))
////			entryOffset = entryOffset + entryLength
////		}
////	} else {
////		// 从后往前
////		var i uint32
////		entryOffset = r.Tail()
////		for i = 0; i < idx; i++ {
////			prevEntryLength := uint32(r.getPrevEntryLengthByCurrentEntryOffset(entryOffset))
////			entryOffset = entryOffset - prevEntryLength
////		}
////	}
////	// entryOffset 就是 idx 对应的 offset
////	return entryOffset
////}
////
////// 根据 offset 获取 以 offset 作为起始地址的节点的长度
////func (r *ZipList) getEntryLengthByOffset(offset uint32) uint32 {
////	// 先计算出 prev，然后根据 prev的长度，计算 encoding 的 offset和数据
////	prev := r.getPrevEntryLengthBytesByCurrentEntryOffset(offset)
////	prevBytesLen := uint32(len(prev))
////	encodingBytes := (*r)[offset+prevBytesLen+1]
////
////	encoding0 := helper.GetBit(encodingBytes, 0)
////	encoding1 := helper.GetBit(encodingBytes, 1)
////	switch encoding0 {
////	case 0:
////		switch encoding1 {
////		case 0:
////			// 00：encoding 为 1 个字节，剩下的 6 个 bit 存储了一个数字，数字的长度表示 entry 存储的 bytes 的长度
////			return prevBytesLen + 1 + uint32(encodingBytes)
////		case 1:
////			// 01：encoding 为 2 个字节，剩下的 6+8 个 bit 存储了一个数字，数字的长度表示 entry 存储的 bytes 的长度
////			// 这里是大端序
////			return 1 + uint32(binary.BigEndian.Uint16([]byte{
////				encodingBytes - 0b01000000, // 减去开头的 01
////				(*r)[offset+1],             // 下个字节
////			}))
////		}
////	case 1:
////		switch encoding1 {
////		case 0:
////			// 10：encoding 为 5 个字节，剩下的 4*8 个 bit 存储了一个数字，数字的长度表示 entry 存储的 bytes 的长度
////			// 这里是大端序
////			return 1 + uint32(binary.BigEndian.Uint32((*r)[offset+1:offset+5]))
////		case 1:
////			switch encodingBytes {
////			case 0b11000000:
////				// 11000000，int16，额外 2 字节
////				return 3
////			case 0b11010000:
////				// 11010000，int32，额外 4 字节
////				return 5
////			case 0b11100000:
////				// 11100000，int64，额外 8 字节
////				return 9
////			case 0b11110000:
////				// 11110000，int24，额外 3 字节，24 位有符号整数
////				return 4
////			case 0x11111110:
////				// 11111110，int8，额外 1 字节，08 位有符号整数
////				return 2
////			default:
////				// 1111xxxx，int4，没有额外的字节，0-12
////				encoding2 := helper.GetBit(encodingBytes, 2)
////				encoding3 := helper.GetBit(encodingBytes, 3)
////				if encoding2 == 1 && encoding3 == 1 {
////					return 1
////				}
////			}
////		}
////	}
////
////	panic(fmt.Sprintf("不可能走到这里"))
////}
//
//func (r *ZipList) getEntryByIdx(idx uint32) (e *entry) {
//	offset := r.getEntryOffsetByIdx(idx)
//	return r.getEntryByOffset(offset)
//}
//
//func (r *ZipList) getEntryByOffset(offset uint32) (e *entry) {
//	e = new(entry)
//	//e.offset = offset
//
//	// prev
//	{
//		prevLengthByte := (*r)[offset]
//		if prevLengthByte < 0xFE {
//			e.prev = []byte{prevLengthByte}
//		} else {
//			e.prev = (*r)[offset : offset+5]
//		}
//	}
//	offset += uint32(len(e.prev))
//
//	// encoding
//	{
//		e.encoding = []byte{(*r)[offset]}
//	}
//
//	// data data
//	encoding0 := helper.GetBit(e.encoding[0], 0)
//	encoding1 := helper.GetBit(e.encoding[0], 1)
//	var extra int64 = -1
//	if encoding0 == 0 && encoding1 == 0 {
//		// 00：encoding 为 1 个字节，剩下的 6 个 bit 存储了一个数字，数字的长度表示 e 存储的 bytes 的长度
//		extra = int64(e.encoding[0])
//	} else if encoding0 == 0 && encoding1 == 1 {
//		// 01：encoding 为 2 个字节，剩下的 6+8 个 bit 存储了一个数字，数字的长度表示 e 存储的 bytes 的长度
//		// ++++ 这里是大端序 ++++
//		e.encoding = append(e.encoding, (*r)[offset+1])
//		extra = int64(e.bytesLen())
//	} else if encoding0 == 1 && encoding1 == 0 {
//		// 10：encoding 为 5 个字节，剩下的 4*8 个 bit 存储了一个数字，数字的长度表示 e 存储的 bytes 的长度
//		// 这里是大端序
//		e.encoding = append(e.encoding, (*r)[offset+1:offset+5]...)
//		extra = int64(e.bytesLen())
//	} else if encoding0 == 1 && encoding1 == 1 {
//		switch e.encoding[0] {
//		case 0b11000000:
//			// 11000000，int16，额外 2 字节
//			extra = 2
//		case 0b11010000:
//			// 11010000，int32，额外 4 字节
//			extra = 4
//		case 0b11100000:
//			// 11100000，int64，额外 8 字节
//			extra = 8
//		case 0b11110000:
//			// 11110000，int24，额外 3 字节，24 位有符号整数
//			extra = 3
//		case 0x11111110:
//			// 11111110，int8，额外 1 字节，08 位有符号整数
//			extra = 1
//		default:
//			// 1111xxxx，int4，没有额外的字节，0-12
//			if helper.GetBit(e.encoding[0], 2) == 1 && helper.GetBit(e.encoding[0], 3) == 1 {
//				extra = 0
//			}
//		}
//	}
//
//	if extra == -1 {
//		panic(fmt.Sprintf("不可能走到这里"))
//	} else if extra > 0 {
//		offset += e.encodingLen()
//		e.data = (*r)[offset : offset+uint32(extra)]
//	}
//
//	return
//}
//
////
////// 根据 entry offset 获取这个节点的上一个节点的长度
////func (r *ZipList) getPrevEntryLengthByCurrentEntryOffset(offset uint32) uint32 {
////	prevLengthBytes := (*r)[offset]
////
////	// 如果小于 254，就 1 个字节，uint8
////	if prevLengthBytes < 0xFE {
////		return uint32(prevLengthBytes)
////	}
////
////	// 否则，是 5 字节，并且第一个字节是 FE(254)，下面的4个字节表示长度 uint32
////	return binary.LittleEndian.Uint32((*r)[offset+1 : offset+5])
////}
////
////// 根据 entry offset 获取这个节点的上一个节点的长度
////func (r *ZipList) getPrevEntryLengthBytesByCurrentEntryOffset(offset uint32) []byte {
////	prevLengthBytes := (*r)[offset]
////
////	// 如果小于 254，就 1 个字节，uint8
////	if prevLengthBytes < 0xFE {
////		return []byte{prevLengthBytes}
////	}
////
////	// 否则，是 5 字节，并且第一个字节是 FE(254)，下面的4个字节表示长度 uint32
////	return (*r)[offset : offset+5]
////}
////
////// * <vlen> 上一个 entry 的长度
////// * 如果小于 254，就 1 个字节，uint8
////// * 否则，是 5 字节，并且第一个字节是 FE(254)，下面的4个字节表示长度 uint32；【所以先取第一个字节，判断接下来的动作】
////func (r *ZipList) getOffsetEntryLengthBytes(offset uint32) []byte {
////	tailLength := r.getEntryLengthByOffset(r.Tail())
////	if tailLength < 0xFE {
////		return []byte{byte(tailLength)}
////	} else {
////		bs := make([]byte, 5)
////		bs[0] = 0xfe
////		binary.LittleEndian.PutUint32(bs[1:5], tailLength)
////		return bs
////	}
////}
//
////func (r *ZipList) getEntryByOffset(offset uint32) []byte {
////	var bs = make([]byte, 0, 5+1+8) // 数字类型的话，最多14字节，bytes的话，靠扩容解决
////	var length uint32
////
////	// prev-length
////	prevLengthBytes := r.getPrevEntryLengthBytesByCurrentEntryOffset(offset)
////	bs = append(bs, prevLengthBytes...)
////	length += uint32(len(prevLengthBytes))
////
////	// encoding
////	r.getOffsetEncoding()
////}
//
//// 将 int64 转成 encoding 和 data 数据
//func intToEncodingBytes(i int64) (byte, []byte) {
//	if 0 <= i && i <= 12 {
//		// 1111xxxx，int32，没有额外的字节，0-12
//		return 12 + 1 + 0b11110000, nil
//	}
//	if math.MinInt8 <= i && i <= math.MaxInt8 {
//		// 11111110，int8，额外 1 字节，8 位有符号整数
//		return 0b11111110, []byte{byte(i)}
//	}
//	if math.MinInt16 <= i && i <= math.MaxInt16 {
//		// 11000000，int16，额外 2 字节
//		return 11000000, []byte{
//			byte(i),
//			byte(i >> 8),
//		}
//	}
//	if gr_binary.MinInt24 <= i && i <= gr_binary.MaxInt24 {
//		// 11110000，int24，额外 3 字节，24 位有符号整数
//		return 11110000, []byte{
//			byte(i),
//			byte(i >> 8),
//			byte(i >> 16),
//		}
//	}
//	if math.MinInt32 <= i && i <= math.MaxInt32 {
//		// 11010000，int32，额外 4 字节
//		return 11110000, []byte{
//			byte(i),
//			byte(i >> 8),
//			byte(i >> 16),
//			byte(i >> 24),
//		}
//	}
//
//	// 11100000，int64，额外 8 字节
//	return 11100000, []byte{
//		byte(i),
//		byte(i >> 8),
//		byte(i >> 16),
//		byte(i >> 24),
//		byte(i >> 32),
//		byte(i >> 40),
//		byte(i >> 48),
//		byte(i >> 56),
//	}
//}
//
//// 打包 int 数据
//func zipIntEntry(prevLen uint32, data interface{}) []byte {
//	e := newEntry(prevLen, data)
//	return e.allbytes()
//}
