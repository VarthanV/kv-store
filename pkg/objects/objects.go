package objects

type ValueType string

const (
	SIMPLE_STRING ValueType = "simple_string"
	ERROR_MESSAGE ValueType = "error_message"
	INTEGER       ValueType = "integer"
	BULK_STRING   ValueType = "bulk_string"
	ARRAY         ValueType = "array"
	NULL                    = "null"
)
