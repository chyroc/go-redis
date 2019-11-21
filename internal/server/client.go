package server

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/database"
	"github.com/chyroc/go-redis/internal/resp"
	"io"
	"net"
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

		args, err := reply.StringSlice()
		if err != nil {
			_, _ = r.writer.Write(resp.NewWithErr(err).Bytes())
		} else {
			res := r.db.ExecCommand(args...)
			_, _ = r.writer.Write(res.Bytes())
		}
	}
}
