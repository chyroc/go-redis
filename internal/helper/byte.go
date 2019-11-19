package helper

// idx: [0, 7]
func GetBit(b byte, idx int) byte {
	return (b >> uint(7-idx)) & 1
}
