package database

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func getMillisecond(args []string, offset int) (off int, ms int64, err error) {
	if len(args) < offset+1 {
		return offset, 0, nil
	}

	r := strings.ToLower(args[offset])
	if r != "ex" && r != "px" {
		return offset, 0, nil
	}
	if len(args) < offset {
		return 0, 0, fmt.Errorf("[Redis.Get] got params %q, but no seconds params", args[offset])
	}

	i, err := strconv.ParseInt(args[offset+1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("[Redis.Get] got params %q, but next param %q is not int", args[offset], args[offset+1])
	}

	if r == "ex" {
		// 秒
		return offset + 2, 1000 * i, nil
	}
	// 毫秒
	return offset + 2, i, nil
}

func getNxXx(args []string, offset int) (nx bool, xx bool, err error) {
	if len(args) < offset+1 {
		return
	}
	r := strings.ToLower(args[offset])
	if r != "nx" && r != "xx" {
		err = fmt.Errorf("[Redis.Get] got params %q, but need %q or %q", args[offset], "NX", "XX")
		return
	}

	if len(args) > offset+1 {
		err = fmt.Errorf("[Redis.Get] got params %q, which endswith %q, but got extra params", args, args[offset])
		return
	}

	return r == "nx", r == "xx", nil
}

func nowMillisecond() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}
