package resp

func NewWithErr(err error) *Reply {
	if err != nil {
		return &Reply{Err: err}
	} else {
		return &Reply{Null: true}
	}
}

func NewWithStr(str string) *Reply {
	return &Reply{Str: str}
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

func NewWithStringSlice(l []string) *Reply {
	var replies []*Reply
	for _, v := range l {
		replies = append(replies, &Reply{
			Err:     nil,
			Null:    false,
			Str:     v,
			Integer: 0,
			Replies: nil,
		})
	}
	return &Reply{Replies: replies}
}

func NewWithReplies(replies []*Reply) *Reply {
	return &Reply{Replies: replies}
}
