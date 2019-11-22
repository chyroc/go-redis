package database

import (
	"fmt"
	"github.com/chyroc/go-redis/internal/basetype"
	"github.com/chyroc/go-redis/internal/resp"
)

type status string

func interfaceToReply(i interface{}) *resp.Reply {
	if i == nil {
		return resp.NewWithNull()
	}

	switch r := i.(type) {
	case status:
		return resp.NewWithStatus(string(r))
	case string:
		return resp.NewWithStr(r)
	case *string:
		if r == nil {
			return resp.NewWithNull()
		}
		return resp.NewWithStr(*r)
	case basetype.SDS:
		return resp.NewWithStr(r.String())
	case *basetype.SDS:
		if r == nil {
			return resp.NewWithNull()
		}
		return resp.NewWithStr(r.String())
	case []string:
		return resp.NewWithStringSlice(r)
	case int:
		return resp.NewWithInt64(int64(r))
	case int8:
		return resp.NewWithInt64(int64(r))
	case int16:
		return resp.NewWithInt64(int64(r))
	case int32:
		return resp.NewWithInt64(int64(r))
	case int64:
		return resp.NewWithInt64(int64(r))
	}

	panic(fmt.Sprintf("%v(%T) 没有处理", i, i))
}
