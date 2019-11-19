package basetype

//import (
//	"encoding/binary"
//)
//
//type ZipList []byte
//
//func NewZipList() ZipList {
//	return []byte{
//		0x0b, 0x00, 0x00, 0x00, // 11 个字节
//		0x0a, 0x00, 0x00, 0x00, // 最后一个的起始地址，10
//		0x00, 0x00, // 0 个元素
//		0xff, // ff 结尾
//	}
//}
//
//// 占用总字节
//func (r *ZipList) BytesLength() uint32 {
//	return binary.LittleEndian.Uint32(r[:4*8])
//}
//
//// 占用总字节
//func (r *ZipList) setBytesLength(i uint32) {
//	binary.LittleEndian.PutUint32(r[:4*8], i)
//}
//
//// 最后一个 entry 的起始字节地址
//func (r *ZipList) Tail() uint32 {
//	return binary.LittleEndian.Uint32(r[4*8 : 8*8])
//}
//
//// 最后一个 entry 的起始字节地址
//func (r *ZipList) setTail(i uint32) {
//	binary.LittleEndian.PutUint32(r[4*8:8*8], i)
//}
//
//// entry 个数
////
//// entry 的个数。当整个值等于 2^16-1 的是，表示个数大于 2 个字节能够存储的上限，需要遍历列表
//func (r *ZipList) Len() uint32 {
//	length := binary.LittleEndian.Uint16(r[8*8 : 10*8])
//	if length == 2^16-1 {
//		panic("不支持这么长的")
//	}
//	return uint32(length)
//}
//
//func (r *ZipList) setLen(length uint32) {
//	if length == 2^16-1 {
//		panic("不支持这么长的")
//	}
//	binary.LittleEndian.PutUint16(r[8*8:10*8], uint16(length))
//}
//
//// 在指定位置插入数据
////
//// idx 是位置，范围 [0, n]，0 表示从左push，n 表示从右push，x表示插入x节点后
//// data是数据，最大值是 int64
//func (r *ZipList) InsertInt(idx uint32, data int64) {
//	// 处理 data
//	//encoding, bs := intToEncodingBytes(data)
//
//	// 分三段
//
//	// 先计算数据从哪里开始复制
//	startOffset := r.getEntryOffsetByIdx(idx)
//	startEntry := r.getEntryByOffset(startOffset)
//
//	//entryBytes := newEntry(startEntry.prevEntryLength(), data)
//	preentry := newEntry(startEntry.prevEntryLength(), data)
//	newdata := make([]byte, 0, r.Tail()-startOffset+preentry.totalLength()+200)
//	var newdataoffset uint32 = 0
//
//	// append data
//	newdata = append(newdata, preentry.allbytes()...)
//	newdataoffset += preentry.totalLength()
//
//	var newPrevEntryLength uint32
//	var stopChainUpdate = false
//	for nextOffset := startOffset; nextOffset <= r.Tail(); {
//		nextEntry := r.getEntryByOffset(nextOffset)
//		nextOffset += nextEntry.totalLength()
//
//		nextEntry.prev = nextEntry.prevToBytes(preentry.totalLength())
//		copy(newdata[newdataoffset:newdataoffset+nextEntry.totalLength()], nextEntry.allbytes())
//		newdataoffset += nextEntry.totalLength()
//		preentry = nextEntry
//
//		if nextEntry.prevEntryLength() < 0xFE && preentry.totalLength() >= 0xFE {
//			// 只有这个情况下，会发生连锁更新
//		} else {
//			stopChainUpdate = true
//		}
//
//		if stopChainUpdate {
//			copy(newdata[newdataoffset:], (*r)[nextOffset:])
//			break
//		}
//	}
//
//	// 将 idx 之后的数据copy到数据结尾
//	// 如果发生了长度变化，重新计算 prev-length
//
//	length := int32(len(bs))
//
//	r.expansion(length)
//}
