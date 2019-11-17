package basetype

// SDS 简单动态字符串
type SDS struct {
	// 已经用的
	len int

	// 尚未用的
	free int

	// 字节数组
	buf []byte
}

func NewSDS() *SDS {
	return &SDS{
		len:  0,
		free: 32,
		buf:  make([]byte, 32),
	}
}

func (s *SDS) Len() int {
	return s.len
}

func (s *SDS) Bytes() []byte {
	return s.buf[:s.len]
}

func (s *SDS) String() string {
	return string(s.buf[:s.len])
}

func (s *SDS) Append(str string) {
	strlen := len(str)
	for strlen > s.free {
		s.expansion()
	}

	copy(s.buf[s.len:s.len+strlen], str)
	s.len += strlen
	s.free -= strlen
}

// 扩容，容量调整为现在的两倍
func (s *SDS) expansion() {
	// 需要锁吗
	total := s.len + s.free
	newbuf := make([]byte, total*2)
	copy(newbuf, s.buf)
	s.buf = newbuf
	s.free = total*2 - s.len
}
