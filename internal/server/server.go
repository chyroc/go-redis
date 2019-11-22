package server

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/database"
	"github.com/chyroc/go-redis/internal/resp"
	"io"
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
		message:  make(chan message, 1000),
	}
}

type serverImpl struct {
	listen   string
	clients  map[string]*redisCli
	database *database.Database
	message  chan message
}

type message struct {
	reply *resp.Reply
	cli   *redisCli
}

func (r *serverImpl) handlequeue() {
	for {
		select {
		case message := <-r.message:
			//fmt.Printf("[got message] %v\n", message)

			cli := message.cli
			reply := message.reply

			fmt.Println(reply.String())

			args, err := reply.StringSlice()
			if err != nil {
				_, _ = cli.writer.Write(resp.NewWithErr(err).Bytes())
			} else {
				_, _ = cli.writer.Write(cli.db.ExecCommand(args...).Bytes())
			}
		}
	}
}

func (r *serverImpl) handleconn(conn net.Conn) {
	var f = func() error {
		cli := r.getClient(conn)
		for {
			reply, err := cli.reader.Read()
			if err != nil {
				return err
			}
			r.message <- message{
				reply: reply,
				cli:   cli,
			}
		}
	}

	defer conn.Close()
	defer delete(r.clients, conn.RemoteAddr().String())

	if err := f(); err != nil {
		if err == io.EOF {
			fmt.Printf("client closed.\n")
			return
		}
		fmt.Printf("[Err] run client, %q\n", err.Error())
	}
}

func (r *serverImpl) Run() error {
	lister, err := net.Listen("tcp", r.listen)
	if err != nil {
		return err
	}
	defer lister.Close()

	go r.handlequeue()

	// 单进程，所以不 `go`
	for {
		conn, err := lister.Accept()
		if err != nil {
			// 忽略，重新 accept
			fmt.Printf("[Err] tcp listen accept, %q\n", err.Error())
			continue
		}

		go r.handleconn(conn)
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
