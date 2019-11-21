package main

import (
	"fmt"
	"log"

	"github.com/chyroc/go-redis/internal/server"
)

func main() {
	r := server.New(":9090")
	go fmt.Println("listen :9090")
	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
