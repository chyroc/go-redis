package helper

import (
	"github.com/spaolacci/murmur3"
)

func Hash(bs []byte) uint64 {
	hasher := murmur3.New64()
	_, _ = hasher.Write(bs)
	return hasher.Sum64()
}
