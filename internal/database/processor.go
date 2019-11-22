package database

type CommandProcessor func(r *RedisDB, args ...string) (interface{}, error)

type commandTemplate struct {
	argsCount int
	processor CommandProcessor
}

// 0 1 2 表示参数个数
// -1 表示最多 x 个
var commandTemplates = map[string]commandTemplate{
	// string
	"get": {
		argsCount: 1,
		processor: Get,
	},
	"set": {
		argsCount: -7,
		processor: Set,
	},
	"getset": {
		argsCount: 2,
		processor: GetSet,
	},
	"strlen": {
		argsCount: 1,
		processor: StrLen,
	},
	"append": {
		argsCount: 2,
		processor: Append,
	},
	"incr": {
		argsCount: 1,
		processor: Incr,
	},

	// expire
	"ttl": {
		argsCount: 1,
		processor: Ttl,
	},
}
