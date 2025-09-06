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
		b.crlf = "\\r\\n"
	} else {
		b.crlf = "\r\n"
	}
	return b
}

func (b *Builder) SimpleString(s string) *Builder {
	b.sb.WriteString("+" + s + b.crlf)
	return b
}

func (b *Builder) Error(err string) *Builder {
	b.sb.WriteString("-" + err + b.crlf)
	return b
}

func (b *Builder) Integer(i int64) *Builder {
	b.sb.WriteString(":" + strconv.FormatInt(i, 10) + b.crlf)
	return b
}

func (b *Builder) Bulk(s string) *Builder {
	b.sb.WriteString(fmt.Sprintf("$%d%s%s%s", len(s), b.crlf, s, b.crlf))
	return b
}

func (b *Builder) Array(elements []string) *Builder {
	b.sb.WriteString(fmt.Sprintf("*%d%s", len(elements), b.crlf))
	for _, e := range elements {
		b.Bulk(e)
	}
	return b
}

func (b *Builder) Null() *Builder {
	b.sb.WriteString("_" + b.crlf)
	return b
}

func (b *Builder) Boolean(v bool) *Builder {
	if v {
		b.sb.WriteString("#t" + b.crlf)
	} else {
		b.sb.WriteString("#f" + b.crlf)
	}
	return b
}

func (b *Builder) Double(f float64) *Builder {
	b.sb.WriteString("," + strconv.FormatFloat(f, 'f', -1, 64) + b.crlf)
	return b
}

func (b *Builder) BigNumber(n *big.Int) *Builder {
	b.sb.WriteString("(" + n.String() + b.crlf)
	return b
}

func (b *Builder) BulkError(err string) *Builder {
	b.sb.WriteString(fmt.Sprintf("!%d%s%s%s", len(err), b.crlf, err, b.crlf))
	return b
}

func (b *Builder) VerbatimString(format, data string) *Builder {
	content := format + ":" + data
	b.sb.WriteString(fmt.Sprintf("=%d%s%s%s", len(content), b.crlf, content, b.crlf))
	return b
}

func (b *Builder) Map(m map[string]string) *Builder {
	b.sb.WriteString(fmt.Sprintf("%%%d%s", len(m), b.crlf))
	for k, v := range m {
		b.Bulk(k)
		b.Bulk(v)
	}
	return b
}

func (b *Builder) Attribute(attrs map[string]string) *Builder {
	b.sb.WriteString(fmt.Sprintf("|%d%s", len(attrs), b.crlf))
	for k, v := range attrs {
		b.Bulk(k)
		b.Bulk(v)
	}
	return b
}

func (b *Builder) Set(elements []string) *Builder {
	b.sb.WriteString(fmt.Sprintf("~%d%s", len(elements), b.crlf))
	for _, e := range elements {
		b.Bulk(e)
	}
	return b
}

func (b *Builder) Push(elements []string) *Builder {
	b.sb.WriteString(fmt.Sprintf(">%d%s", len(elements), b.crlf))
	for _, e := range elements {
		b.Bulk(e)
	}
	return b
}

func (b *Builder) Build() string {
	return b.sb.String()
}
