package resp

import "github.com/VarthanV/kv-store/pkg/objects"

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

func (v *Value) Marshal() []byte {
	switch v.typ {
	case objects.ARRAY:
		return []byte{}
	default:
		return []byte{}
	}
}
