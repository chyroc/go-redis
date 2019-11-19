package gr_binary

const (
	MinInt24 = -1 << 23
	MaxInt24 = 1<<23 - 1
)

var LittleEndian littleEndian

type littleEndian struct{}

func (littleEndian) Uint8(b []byte) uint8 {
	_ = b[0] // bounds check hint to compiler; see golang.org/issue/14808
	return uint8(b[0])
}

func (littleEndian) PutUint8(b []byte, v uint8) {
	_ = b[0] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
}

func (littleEndian) Uint24(b []byte) uint32 {
	_ = b[2] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

func (littleEndian) PutUint24(b []byte, v uint32) {
	_ = b[0] // early bounds check to guarantee safety of writes below
	_ = b[7] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
}

//type bigEndian struct{}
//
//var BigEndian bigEndian
//
//func (bigEndian) Uint16(b []byte) uint16 {
//	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
//	return uint16(b[1]) | uint16(b[0])<<8
//}
//
//func (bigEndian) PutUint16(b []byte, v uint16) {
//	_ = b[1] // early bounds check to guarantee safety of writes below
//	b[0] = byte(v >> 8)
//	b[1] = byte(v)
//}
