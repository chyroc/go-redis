package basetype

// 整数集合

const (
	INTSET_ENC_INT16 = 2
	INTSET_ENC_INT32 = 4
	INTSET_ENC_INT64 = 8
)

type Intset struct {
	encoding uint32
	length   uint32
	contents []byte
}

func NewIntset() *Intset {
	return &Intset{
		encoding: INTSET_ENC_INT16,
		length:   0,
		contents: make([]byte, 8),
	}
}

func (is *Intset) Add(s interface{}) {
	is.add(s)
}

func (is *Intset) Exist(s int64) bool {
	return is.exist(s)
}
