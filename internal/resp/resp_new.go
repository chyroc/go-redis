package resp

func NewWithErr(err error) *Reply {
	if err != nil {
		return &Reply{Err: err}
	}
	return nil
}

func NewWithBytes(bs []byte) *Reply {
	return &Reply{Str: string(bs)}
}

func NewWithInt64(i int64) *Reply {
	return &Reply{Integer: i}
}

func NewWithNull() *Reply {
	return &Reply{Null: true}
}

func NewWithReplies(replies []*Reply) *Reply {
	return &Reply{Replies: replies}
}
