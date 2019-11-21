package resp

import (
	"bytes"
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

	} else if r.Integer != 0 {

		// 整数回复 ":"

		buf.WriteByte(':')
		buf.WriteString(strconv.FormatInt(r.Integer, 10))

		buf.WriteByte(cr)
		buf.WriteByte(lf)

	} else if r.Str != "" {

		// 状态回复 "+"

		buf.WriteByte('+')
		buf.WriteString(r.Str)

		buf.WriteByte(cr)
		buf.WriteByte(lf)

	} else if len(r.Replies) > 0 {

		// 多条批量回复 "*"

		buf.WriteByte('*')

		buf.WriteString(strconv.Itoa(len(r.Replies)))
		buf.WriteByte(cr)
		buf.WriteByte(lf)

		for _, v := range r.Replies {
			buf.Write(v.Bytes())
		}
	}

	return buf.Bytes()
}
