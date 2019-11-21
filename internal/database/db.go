package database

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/basetype"
	"github.com/chyroc/go-redis/internal/resp"
)

type Database struct {
	dbnum int // 1 - 16
	db    []*RedisDB
}

func New() *Database {
	d := new(Database)
	d.dbnum = 16
	d.db = make([]*RedisDB, d.dbnum)
	for i := 0; i < d.dbnum; i++ {
		d.db[i] = newRedisSB()
	}

	return d
}

func (d *Database) DB(i int) (*RedisDB) {
	return d.db[i]
}

func (r *RedisDB) ExecCommand(args ...string) *resp.Reply {
	if len(args) == 0 {
		return resp.NewWithErr(fmt.Errorf("至少需要一个命令"))
	}

	cmd := args[0]
	args = args[1:]

	t, ok := commandTemplates[cmd]
	if !ok {
		return resp.NewWithErr(fmt.Errorf("%q 命令不支持", cmd))
	}
	if t.argsCount >= 0 && len(args) != t.argsCount {
		return resp.NewWithErr(fmt.Errorf("%q 命令的参数需要 %d 个，但是传递了 %d 个", t.argsCount, len(args), ))
	}

	res, err := t.processor(r, args...)
	if err != nil {
		return resp.NewWithErr(err)
	} else {
		return interfaceToReply(res)
	}
}

type RedisDB struct {
	dict    *basetype.Dict // 键值对
	expires *basetype.Dict // 过期时间 int64，毫秒时间戳
}

func newRedisSB() *RedisDB {
	return &RedisDB{
		dict:    basetype.NewDict(),
		expires: basetype.NewDict(),
	}
}
