package helper

import (
	"encoding/binary"
)

// 备注
// 内存实现、rdb 持久化，用小端
// reply 协议等网络实现，用大端

//var byte2pool = sync.Pool{
//	New: func() interface{} {
//		return make([]byte, 2)
//	},
//}

func Int16ToBinary(i int16) []byte {
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, uint16(i))
	return result
}

func Int32ToBinary(i int16) []byte {
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, uint32(i))
	return result
}

func Int64ToBinary(i int16) []byte {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(i))
	return result
}
