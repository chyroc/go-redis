package server_test

import (
	"os"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/chyroc/go-redis/internal/server"
)

var client *redis.Client
var addr = os.Getenv("GO_REDIS_PORT") // GO_REDIS_PORT=:9091 go test ./...

func init() {
	go func() {
		if addr == ":9090" {
			return
		}
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
