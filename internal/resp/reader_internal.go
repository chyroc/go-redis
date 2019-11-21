package resp

import (
	"bufio"
	"strconv"
)

const (
	lf byte = 10 // \n
	cr byte = 13 // \r
)

type parser struct {
	reader *bufio.Reader
}

func (r *parser) readUntilCRLF() ([]byte, error) {
	bs, err := r.reader.ReadBytes(lf)
	if err != nil {
		return bs, err
	}

	l := len(bs)
	if l >= 2 && bs[l-2] == cr {
		return bs[:l-2], nil
	}

	return bs, nil
}

func (r *parser) readIntBeforeCRLF() (int64, error) {
	length, err := r.readUntilCRLF()
	if err != nil {
		return 0, err
	}
	c, err := strconv.ParseInt(string(length), 10, 64)
	if err != nil {
		return 0, err
	}
	return c, nil
}

func (r *parser) readBytes(length int) ([]byte, error) {
	bs := make([]byte, length)
	if _, err := r.reader.Read(bs); err != nil {
		return nil, err
	}
	return bs, nil
}
