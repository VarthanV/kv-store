package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/VarthanV/kv-store/pkg/objects"
)

// Specs for RESP protocol
// https://redis.io/docs/latest/develop/reference/protocol-spec/

// First byte - Type identifiers

type DataType rune

const (
	STRING  DataType = '+'
	ERROR   DataType = '-'
	INTEGER DataType = ':'
	BULK    DataType = '$'
	ARRAY   DataType = '*'
)

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

func (r *Resp) readIntegerFromInput() (val int, n int, err error) {
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

	switch DataType(typ) {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	case INTEGER:
		return r.readInteger()

	default:
		return nil, fmt.Errorf("unknown type %s", string(typ))
	}
}

func (r *Resp) readArray() (*Value, error) {
	v := Value{}
	v.Typ = "array"

	// length of array first to be read
	len, _, err := r.readIntegerFromInput()
	if err != nil {
		return nil, err
	}

	v.Arr = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return nil, err
		}

		// append parsed value to array
		v.Arr = append(v.Arr, *val)
	}

	return &v, nil
}

func (r *Resp) readBulk() (*Value, error) {
	v := Value{}

	v.Typ = "bulk"

	len, _, err := r.readIntegerFromInput()
	if err != nil {
		return nil, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.Bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return &v, nil
}

func (r *Resp) readInteger() (*Value, error) {
	v := Value{}
	v.Typ = objects.INTEGER

	val, _, err := r.readIntegerFromInput()
	if err != nil {
		return nil, err
	}
	v.Num = val
	return &v, nil
}
