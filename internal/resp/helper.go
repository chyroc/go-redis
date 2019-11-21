package resp

import (
	"errors"
)

var ErrUnSupportRespType = errors.New("unsupported redis resp type")
var ErrKeyNotExist = errors.New("key not exist")
