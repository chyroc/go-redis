package database

import "github.com/chyroc/go-redis/internal/resp"

func interfaceToReply(i interface{}) *resp.Reply {
	if i == nil {
		return resp.NewWithNull()
	}

	switch r := i.(type) {
	case string:
		return resp.NewWithStr(r)
	case *string:
		if r == nil {
			return resp.NewWithNull()
		} else {
			return resp.NewWithStr(*r)
		}
	case []string:
		return resp.NewWithStringSlice(r)
	}

	panic(i)
	return nil
}