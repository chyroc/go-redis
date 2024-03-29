package resp

import (
	"bytes"
	"fmt"
	"strconv"
)

// http://redisdoc.com/topic/protocol.html
func (r *Reply) Bytes() []byte {
	var buf = new(bytes.Buffer)

	if r.Err != nil {

		// 错误回复 "-"

		buf.WriteByte('-')
		buf.WriteString(r.Err.Error())

		buf.WriteByte(cr)
		buf.WriteByte(lf)

	} else if r.Null {

		buf.WriteByte('$')
		buf.WriteString(strconv.Itoa(-1))

		buf.WriteByte(cr)
		buf.WriteByte(lf)

	} else if r.ReplyType == replyTypeInt {

		// 整数回复 ":"

		buf.WriteByte(':')
		buf.WriteString(strconv.FormatInt(r.Integer, 10))

		buf.WriteByte(cr)
		buf.WriteByte(lf)

	} else if r.ReplyType == replyTypeStatus {

		// 状态回复 "+"

		buf.WriteByte('+')
		buf.WriteString(r.str)

		buf.WriteByte(cr)
		buf.WriteByte(lf)

	} else if r.ReplyType == ReplyTypeReplies {

		// 多条批量回复 "*"

		buf.WriteByte('*')

		buf.WriteString(strconv.Itoa(len(r.Replies)))
		buf.WriteByte(cr)
		buf.WriteByte(lf)

		for _, v := range r.Replies {
			buf.Write(v.Bytes())
		}
	} else if r.ReplyType == replyTypeString {

		buf.WriteByte('$')
		buf.WriteString(strconv.Itoa(len(r.str)))

		buf.WriteByte(cr)
		buf.WriteByte(lf)

		buf.WriteString(r.str)
		buf.WriteByte(cr)
		buf.WriteByte(lf)

	} else {
		panic(fmt.Sprintf("reply bytes 错误"))
	}

	return buf.Bytes()
}
