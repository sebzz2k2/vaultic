package resp

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type Builder struct {
	sb    strings.Builder
	debug bool
	crlf  string
}

func NewBuilder(debug bool) *Builder {
	b := &Builder{debug: debug}
	if debug {
		b.crlf = ESC_CRLF
	} else {
		b.crlf = CRLF
	}
	return b
}

func (b *Builder) SimpleString(s string) *Builder {

	b.sb.WriteString(string(respTypeToChar[SIMPLE_STRING]) + s + b.crlf)
	return b
}

func (b *Builder) Error(err string) *Builder {
	b.sb.WriteString(string(respTypeToChar[ERROR]) + err + b.crlf)
	return b
}

func (b *Builder) Integer(i int64) *Builder {
	b.sb.WriteString(string(respTypeToChar[INTEGER]) + strconv.FormatInt(i, 10) + b.crlf)
	return b
}

func (b *Builder) Bulk(s string) *Builder {
	b.sb.WriteString(fmt.Sprintf("%c%d%s%s%s", respTypeToChar[BULK_STRING], len(s), b.crlf, s, b.crlf))
	return b
}

func (b *Builder) Array(elements []string) *Builder {
	b.sb.WriteString(fmt.Sprintf("%c%d%s", respTypeToChar[ARRAY], len(elements), b.crlf))
	for _, e := range elements {
		b.Bulk(e)
	}
	return b
}

func (b *Builder) Null() *Builder {
	b.sb.WriteString(string(respTypeToChar[NULL]) + b.crlf)
	return b
}

func (b *Builder) Boolean(v bool) *Builder {
	if v {
		b.sb.WriteString(TRUE + b.crlf)
	} else {
		b.sb.WriteString(FALSE + b.crlf)
	}
	return b
}

func (b *Builder) Double(f float64) *Builder {
	b.sb.WriteString(string(respTypeToChar[DOUBLE]) + strconv.FormatFloat(f, 'f', -1, 64) + b.crlf)
	return b
}

func (b *Builder) BigNumber(n *big.Int) *Builder {
	b.sb.WriteString(string(respTypeToChar[BIG_NUMBER]) + n.String() + b.crlf)
	return b
}

func (b *Builder) BulkError(err string) *Builder {
	b.sb.WriteString(fmt.Sprintf("%c%d%s%s%s", respTypeToChar[BULK_ERROR], len(err), b.crlf, err, b.crlf))
	return b
}

func (b *Builder) VerbatimString(format, data string) *Builder {
	content := format + ":" + data
	b.sb.WriteString(fmt.Sprintf("%c%d%s%s%s", respTypeToChar[VERBATIM], len(content), b.crlf, content, b.crlf))
	return b
}

func (b *Builder) Map(m map[string]string) *Builder {
	b.sb.WriteString(fmt.Sprintf("%c%d%s", respTypeToChar[MAP], len(m), b.crlf))
	for k, v := range m {
		b.Bulk(k)
		b.Bulk(v)
	}
	return b
}

func (b *Builder) Attribute(attrs map[string]string) *Builder {
	b.sb.WriteString(fmt.Sprintf("%c%d%s", respTypeToChar[ATTRIBUTES], len(attrs), b.crlf))
	for k, v := range attrs {
		b.Bulk(k)
		b.Bulk(v)
	}
	return b
}

func (b *Builder) Set(elements []string) *Builder {
	b.sb.WriteString(fmt.Sprintf("%c%d%s", respTypeToChar[SET], len(elements), b.crlf))
	for _, e := range elements {
		b.Bulk(e)
	}
	return b
}

func (b *Builder) Push(elements []string) *Builder {
	b.sb.WriteString(fmt.Sprintf("%c%d%s", respTypeToChar[PUSH], len(elements), b.crlf))
	for _, e := range elements {
		b.Bulk(e)
	}
	return b
}

func (b *Builder) Build() string {
	return b.sb.String()
}
