package basetype

import (
	cryptorand "crypto/rand"
	"math/big"
	mathrand "math/rand"
	"testing"
	"time"
)
import "github.com/stretchr/testify/assert"

func TestSDS(t *testing.T) {
	as := assert.New(t)

	sds := NewSDS()
	as.Equal(sds.String(), "")
	as.Equal(sds.Len(), 0)

	total := 0
	for i := 0; i < 100; i++ {
		s := RandStringRunes(RandomInt())
		sds.Append(s)
		total += len(s)
	}
	as.Equal(sds.Len(), total)
	//t.Logf("sds is %s", sds.string())
}

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

func RandomInt() int {
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(100))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mathrand.Intn(len(letterRunes))]
	}
	return string(b)
}
