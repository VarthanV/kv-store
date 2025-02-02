package resp

import (
	"bytes"
	"strconv"

	"github.com/VarthanV/kv-store/pkg/objects"
)

type Value struct {
	// type inferred from the first byte
	typ objects.ValueType
	// key received when asked to store key value pair
	key string
	// string received when asked to store simple string
	str string
	// string received when asked to store bulk string
	bulk string
	// num received when asked to store integer
	num int
	// array received when asked to store array
	arr []Value
}

// Marshal: marshals the current value to RESP format
func (v *Value) Marshal() []byte {
	switch v.typ {
	case objects.ARRAY:
		return v.marshalArray()
	case objects.SIMPLE_STRING:
		return v.marshalString()
	case objects.BULK_STRING:
		return v.marshalBulk()
	case objects.INTEGER:
		return v.marshalInteger()
	case objects.ERROR_MESSAGE:
		return v.marshallError()
	case objects.NULL:
		return v.marshallNull()
	default:
		return []byte{}
	}
}

// marshalString: marshals the string to RESP format
func (v *Value) marshalString() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(STRING))
	buf.WriteString(v.str)
	buf.WriteString("\r\n")
	return buf.Bytes()
}

// marshalBulk: marshals the bulk string to RESP format
func (v *Value) marshalBulk() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(BULK))
	// Append len
	buf.WriteString(strconv.Itoa(len(v.bulk)))
	buf.WriteString("\r\n")
	for i := 0; i < len(v.bulk); i++ {
		buf.WriteByte(v.bulk[i])
	}
	buf.WriteString("\r\n")
	return buf.Bytes()
}

// marshalInteger: marshals the integer to RESP format
func (v *Value) marshalInteger() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(INTEGER))
	buf.WriteString(strconv.Itoa(v.num))
	buf.WriteString("\r\n")
	return buf.Bytes()
}

// marshalArray: marshals the array to RESP format
func (v *Value) marshalArray() []byte {
	var buf bytes.Buffer
	len := len(v.arr)
	buf.WriteByte(byte(ARRAY))
	buf.WriteString(strconv.Itoa(len))
	// CRLF
	buf.WriteString(`\r\n`)
	for i := 0; i < len; i++ {
		buf.Write(v.arr[i].Marshal())
	}
	return buf.Bytes()
}

// marshallError: marshals the error to RESP format
func (v Value) marshallError() []byte {
	var buf bytes.Buffer
	buf.WriteByte(byte(ERROR))
	buf.WriteString(v.str)
	buf.WriteString(`\r\n`)
	return buf.Bytes()
}

// marshallNull: marshals the null value to RESP format
func (v *Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}
