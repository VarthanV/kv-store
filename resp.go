package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Specs for RESP protocol
// https://redis.io/docs/latest/develop/reference/protocol-spec/

// First byte - Type identifiers
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	// type inferred from the first byte
	typ string
	// string received when asked to store simple string
	str string
	// string received when asked to store bulk string
	bulk string
	// num received when asked to store integer
	num int
	// array received when asked to store array
	arr []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// readLine , read the line one byte at a time
// until we find the \r , return last -2 bytes which is \r\n
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}

	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(i), n, nil

}

func (r *Resp) Read() (*Value, error) {
	typ, err := r.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch typ {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		return nil, fmt.Errorf("unknown type %s", string(typ))
	}
}

func (r *Resp) readArray() (*Value, error) {
	v := Value{}
	v.typ = "array"

	// length of array first to be read
	len, _, err := r.readInteger()
	if err != nil {
		return nil, err
	}

	v.arr = make([]Value, len)
	for i := 0; i < len; i++ {

	}

	return &v, nil
}

func (r *Resp) readBulk() (*Value, error) {
	v := Value{}

	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return nil, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return &v, nil
}
