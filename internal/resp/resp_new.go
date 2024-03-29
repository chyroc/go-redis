package resp

func NewWithErr(err error) *Reply {
	if err != nil {
		return &Reply{Err: err}
	} else {
		return NewWithNull()
	}
}

func NewWithStatus(str string) *Reply {
	return &Reply{str: str, ReplyType: replyTypeStatus}
}

func NewWithStr(str string) *Reply {
	return &Reply{str: str, ReplyType: replyTypeString}
}

func NewWithInt64(i int64) *Reply {
	return &Reply{Integer: i, ReplyType: replyTypeInt}
}

func NewWithNull() *Reply {
	return &Reply{Null: true}
}

func NewWithStringSlice(l []string) *Reply {
	var replies []*Reply
	for _, v := range l {
		replies = append(replies, NewWithStr(v))
	}
	return NewWithReplies(replies)
}

func NewWithReplies(replies []*Reply) *Reply {
	return &Reply{Replies: replies, ReplyType: ReplyTypeReplies}
}
