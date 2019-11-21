package server

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/database"
	"io"
	"net"
	"strings"

	"github.com/chyroc/go-redis/internal/resp"
)

type redisCli struct {
	addr   string
	parser resp.Parser
	writer io.Writer

	database *database.Database
	db       *database.RedisDB
}

func newRedisCli(conn net.Conn, database_ *database.Database) *redisCli {
	return &redisCli{
		addr:     conn.RemoteAddr().String(),
		parser:   resp.NewParser(conn),
		writer:   conn,
		database: database_,
		db:       database_.DB(0),
	}
}

func (r *redisCli) run() error {
	// todo: 多客户端
	for {
		reply, err := r.parser.Read()
		if err != nil {
			return err
		}

		for _, v := range reply.Replies {
			if v.Str != "" {
				switch strings.ToLower(v.Str) {
				case "command":
					res := resp.NewWithStringSlice([]string{"GET"})
					_, err := r.writer.Write(res.Bytes())
					return err
				case "get":
				case "set":
				}
			}
			fmt.Println(v.String())
		}
	}
}
