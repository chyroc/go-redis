package ints

import "fmt"

func IsInt(i interface{}) bool {
	switch i.(type) {
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
	default:
		return false
	}
	return true
}

func ToInt64(i interface{}) int64 {
	switch ii := i.(type) {
	case int:
		return int64(ii)
	case int8:
		return int64(ii)
	case int16:
		return int64(ii)
	case int32:
		return int64(ii)
	case int64:
		return ii
	}
	panic(fmt.Sprintf("%s(%T) 无法转为 int64", i, i))
}

func ByteLen(i interface{}) int {
	switch i.(type) {
	case int, int32:
		return 4
	case int8:
		return 1
	case int16:
		return 2
	case int64:
		return 8
	}
	panic(fmt.Sprintf("%s(%T) 无法转为 int64", i, i))
}
