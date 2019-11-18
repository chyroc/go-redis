package helper

import (
	"fmt"
	"strconv"
	"strings"
)

func String(data interface{}) string {
	switch s := data.(type) {
	case string:
		return s
	case float64:
		if s == 0 {
			return "0"
		}
		return strings.TrimSuffix(strconv.FormatFloat(s, 'f', -1, 64), "0")
	case float32:
		if s == 0 {
			return "0"
		}
		return strings.TrimSuffix(strconv.FormatFloat(float64(s), 'f', -1, 32), "0")
	case int:
		return strconv.Itoa(s)
	}

	panic(fmt.Sprintf("%T 不支持转字符串", data))
}
