package basetype

import (
	"strconv"
)

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

func NewSDSWithString(s string) *SDS {
	sds := NewSDS()
	sds.Append(s)
	return sds
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

func (s *SDS) Int64() (int64, error) {
	return strconv.ParseInt(s.String(), 10, 64)
}

func (s *SDS) Int64Incr() (int64, error) {
	bs := s.buf[:s.len]
	i, err := strconv.ParseInt(string(bs), 10, 64)
	if err != nil {
		return 0, err
	}
	i++
	j := strconv.FormatInt(i, 10)
	if len(j) > len(bs) {
		s.expansion(1)
		s.len++
		s.free--
	}

	for idx, c := range j {
		s.buf[idx] = byte(c)
	}
	return i,nil
}

func (s *SDS) Append(str string) {
	strlen := len(str)
	for strlen > s.free {
		s.expansion(strlen)
	}

	copy(s.buf[s.len:s.len+strlen], str)
	s.len += strlen
	s.free -= strlen
}

func (s *SDS) EqualToString(str string) bool {
	strlen := len(str)
	if strlen != s.len {
		return false
	}
	return s.String() == str
}

// 扩容
// 需要锁吗
func (s *SDS) expansion(neenlen int) {
	if s.free >= neenlen {
		return
	}

	total := s.len + s.free
	var newtotal int
	if total < 1024 {
		newtotal = total * 2 // 两步
	} else {
		newtotal = 1024 * (total/1024 + 1) // 1024的备注
	}

	newbuf := make([]byte, newtotal)
	copy(newbuf, s.buf)
	s.buf = newbuf
	s.free = newtotal - s.len
}
