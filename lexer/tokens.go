package lexer

import "fmt"

type TokenKind int

const (
	//  CMDS
	CMD_GET TokenKind = iota
	CMD_SET
	CMD_DEL
	CMD_EXISTS
	CMD_KEYS

	NEW_LINE
	VALUE
	WHITESPACE
)

type Token struct {
	Kind  TokenKind
	Value string
}

func (t Token) isOneOf(kinds ...TokenKind) bool {
	for _, kind := range kinds {
		if t.Kind == kind {
			return true
		}
	}
	return false
}

func NewToken(kind TokenKind, value string) Token {
	return Token{
		Kind:  kind,
		Value: value,
	}
}

var reserved_literal map[string]TokenKind = map[string]TokenKind{
	"GET":    CMD_GET,
	"SET":    CMD_SET,
	"DEL":    CMD_DEL,
	"EXISTS": CMD_EXISTS,
	"KEYS":   CMD_KEYS,
}

func TokenKindToString(kind TokenKind) string {
	switch kind {
	case CMD_GET:
		return "GET"
	case CMD_SET:
		return "SET"
	case CMD_DEL:
		return "DEL"
	case CMD_EXISTS:
		return "EXISTS"
	case CMD_KEYS:
		return "KEYS"
	case NEW_LINE:
		return "NEW_LINE"
	case VALUE:
		return "VALUE"
	case WHITESPACE:
		return "WHITESPACE"
	default:
		return fmt.Sprintf("Unknown TokenKind: %d", kind)
	}
}

func DebugToken(t Token) string {
	return fmt.Sprintf("{Kind: %s, Value: %s}", TokenKindToString(t.Kind), t.Value)
}
