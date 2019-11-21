package database

import "github.com/chyroc/go-redis/internal/basetype"

type Database struct {
	dbnum int // 1 - 16
	db    []*redisDB
}

func New() *Database {
	d := new(Database)
	d.dbnum = 16
	d.db = make([]*redisDB, d.dbnum)
	for i := 0; i < d.dbnum; i++ {
		d.db[i] = newRedisSB()
	}

	return d
}

func (d *Database) Run(cmd string, args ...string) {

}

type redisDB struct {
	dict    *basetype.Dict // 键值对
	expires *basetype.Dict // 过期时间 uint64，毫秒时间戳
}

func newRedisSB() *redisDB {
	return &redisDB{
		dict: basetype.NewDict(),
	}
}
