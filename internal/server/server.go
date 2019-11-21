package server

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/database"
	"net"
)

type Server interface {
	Run() error
}

func New(listen string) Server {
	return &serverImpl{
		listen:   listen,
		clients:  map[string]*redisCli{},
		database: database.New(),
	}
}

type serverImpl struct {
	listen   string
	clients  map[string]*redisCli
	database *database.Database
}

func (r *serverImpl) Run() error {
	lister, err := net.Listen("tcp", r.listen)
	if err != nil {
		return err
	}
	defer lister.Close()

	// 单进程，所以不 `go`
	for {
		conn, err := lister.Accept()
		if err != nil {
			// 忽略，重新 accept
			fmt.Printf("[Err] tcp listen accept, %q\n", err.Error())
			continue
		}

		if err := r.clientRun(conn); err != nil {
			fmt.Printf("[Err] run client, %q\n", err.Error())
			continue
		}
	}
}

func (r *serverImpl) getClient(conn net.Conn) *redisCli {
	addr := conn.RemoteAddr().String()

	if _, ok := r.clients[addr]; !ok {
		r.clients[addr] = newRedisCli(conn, r.database)
		fmt.Println("【新客户端】addr", addr)
	} else {
		fmt.Println("【老客户端】addr", addr)
	}
	return r.clients[addr]
}

func (r *serverImpl) clientRun(conn net.Conn) error {
	//defer conn.Close()

	cli := r.getClient(conn)

	defer delete(r.clients, cli.addr)

	return cli.run()
}
