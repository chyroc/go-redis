package server

import (
	"github.com/chyroc/go-redis/internal/database"
	"github.com/chyroc/go-redis/internal/resp"
	"io"
	"net"
)

type redisCli struct {
	addr   string
	reader resp.Parser
	writer io.Writer

	database *database.Database
	db       *database.RedisDB
}

func newRedisCli(conn net.Conn, database_ *database.Database) *redisCli {
	return &redisCli{
		addr:     conn.RemoteAddr().String(),
		reader:   resp.NewParser(conn),
		writer:   conn,
		database: database_,
		db:       database_.DB(0),
	}
}
