package database

import "errors"

var (
	ErrOperationWrongKindValue = errors.New(`WRONGTYPE Operation against a key holding the wrong kind of value`)
	ErrInvalidInt              = errors.New(`ERR value is not an integer or out of range`)
)
