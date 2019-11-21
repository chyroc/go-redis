package main

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/server"
	"log"
)

func main() {
	r := server.New(":9090")
	go fmt.Println("listen :9090")
	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
