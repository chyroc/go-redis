package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type Parser interface {
	Read() (*Reply, error)
}

func NewParser(reader io.Reader) Parser {
	return &parser{
		reader: bufio.NewReader(reader),
	}
}

func (r *parser) Read() (*Reply, error) {
	respType, err := r.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch respType {
	case '+':
		res, err := r.readUntilCRLF()
		if err != nil {
			return nil, err
		}
		return NewWithStr(string(res)), nil
	case '-':
		message, err := r.readUntilCRLF()
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(message)) // TODO 错误类型
	case ':':
		length, err := r.readIntBeforeCRLF()
		if err != nil {
			return nil, err
		}
		return NewWithInt64(length), nil
	case '$':
		length, err := r.readIntBeforeCRLF()
		if err != nil {
			return nil, err
		}

		if length == -1 {
			return NewWithNull(), nil
		}

		bs, err := r.readBytes(int(length))
		if err != nil {
			return nil, err
		}

		return NewWithStr(string(bs)), nil
	case '*':
		// multi bulk reply
		count, err := r.readIntBeforeCRLF()
		if err != nil {
			return nil, err
		}

		var replys []*Reply
		for i := 0; i < int(count); i++ {
			reply, err := r.Read()
			if err != nil {
				return nil, err
			}
			replys = append(replys, reply)
		}

		return NewWithReplies(replys), nil
	}

	return nil, fmt.Errorf("%b(%s): %w", respType, []byte{respType}, ErrUnSupportRespType)
}
