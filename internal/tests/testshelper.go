package tests

import (
	cryptorand "crypto/rand"
	"math/big"
	mathrand "math/rand"
	"time"
)

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

func RandomInt64(max int64) int64 {
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}
	return n.Int64()
}

func RandomIntIn100() int {
	return int(RandomInt64(100))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mathrand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandString32() string {
	b := make([]rune, 32)
	for i := range b {
		b[i] = letterRunes[mathrand.Intn(len(letterRunes))]
	}
	return string(b)
}
