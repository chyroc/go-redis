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
		fmt.Println(reply.String())
		p := reply.Replies

		if len(p) == 0 {
			return fmt.Errorf("空命令")
		}

		switch strings.ToLower(p[0].Str) {
		case "command":
			res := resp.NewWithStringSlice([]string{"GET"})
			_, _ = r.writer.Write(res.Bytes())
		case "get":
			str, err := r.db.Get(p[1].Str)
			if err != nil {
				_, _ = r.writer.Write(resp.NewWithErr(err).Bytes())
			} else {
				if str == nil {
					_, _ = r.writer.Write(resp.NewWithNull().Bytes())
				} else {
					_, _ = r.writer.Write(resp.NewWithStr(*str).Bytes())
				}
			}
		case "set":
			if err := r.db.Set(p[1].Str, p[2].Str); err != nil {
				_, _ = r.writer.Write(resp.NewWithErr(err).Bytes())
			} else {
				_, _ = r.writer.Write(resp.NewWithNull().Bytes())
			}
		}
	}
}
