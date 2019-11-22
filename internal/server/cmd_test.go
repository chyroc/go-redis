package server_test

import (
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/chyroc/go-redis/internal/server"
)

var client *redis.Client
var addr = ":9090"

func init() {
	go func() {
		if addr == ":9090" {
			return
		}
		panic("xx")
		if err := server.New(addr).Run(); err != nil {
			panic(err)
		}
	}()

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	time.Sleep(time.Second / 2)
}
