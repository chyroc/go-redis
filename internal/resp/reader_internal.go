package resp

import (
	"bufio"
	"fmt"
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
	bs := make([]byte, length+2)
	n, err := r.reader.Read(bs)
	if err != nil {
		return nil, err
	}
	if n != length+2 {
		return nil, fmt.Errorf("expect to read %d bytes, but got %d bytes", length+2, n)
	}
	return bs[:length], nil
}
