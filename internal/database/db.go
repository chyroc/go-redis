package database

import "github.com/chyroc/go-redis/internal/basetype"

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

func (d *Database) Start() error {
	return nil
}

func (d *Database) ExecCommand(cmd string, args ...string) {

}

type RedisDB struct {
	dict    *basetype.Dict // 键值对
	expires *basetype.Dict // 过期时间 uint64，毫秒时间戳
}

func newRedisSB() *RedisDB {
	return &RedisDB{
		dict: basetype.NewDict(),
	}
}
