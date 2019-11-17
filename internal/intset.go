package internal

type intset struct {
	encoding uint32
	length   uint32
	contents []byte
}
