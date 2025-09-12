package resp

const (
	SIMPLE_STRING = "simple_string"
	ERROR         = "error"
	INTEGER       = "integer"
	BULK_STRING   = "bulk_string"
	ARRAY         = "array"
	NULL          = "null"
	BOOLEAN       = "boolean"
	DOUBLE        = "double"
	BIG_NUMBER    = "big_number"
	BULK_ERROR    = "bulk_error"
	VERBATIM      = "verbatim_string"
	MAP           = "map"
	ATTRIBUTES    = "attributes"
	SET           = "set"
	PUSH          = "push"
	CRLF          = "\r\n"
	ESC_CRLF      = "\\r\\n"
	TRUE          = "#t"
	FALSE         = "#f"
)

var respTypeToChar = map[string]byte{
	SIMPLE_STRING: '+',
	ERROR:         '-',
	INTEGER:       ':',
	BULK_STRING:   '$',
	ARRAY:         '*',
	NULL:          '_',
	BOOLEAN:       '#',
	DOUBLE:        ',',
	BIG_NUMBER:    '(',
	BULK_ERROR:    '!',
	VERBATIM:      '=',
	MAP:           '%',
	ATTRIBUTES:    '|',
	SET:           '~',
	PUSH:          '>',
}
