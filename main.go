package main

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/basetype"
	"strconv"
)

func main() {
	{
		s := basetype.NewSkipList()
		for i := 1; i < 11; i++ {
			_ = s.Add(strconv.Itoa(i), float64(i))
		}

		fmt.Printf("[SkipList] %s\n", s)
		d, _ := s.Get(float64(1))
		node := d[0]
		fmt.Println(node)
	}
}
