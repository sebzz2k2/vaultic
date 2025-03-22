package lexer

import (
	"regexp"
	"strings"
)

type lexer struct {
	input  string
	pos    int
	regex  []regexPattern
	tokens []Token
}

type regexPattern struct {
	pattern *regexp.Regexp
	handler regexHandler
}

type regexHandler func(lex *lexer, regex *regexp.Regexp)

func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.input[lex.pos:])

	if kind, ok := reserved_literal[match]; ok {
		lex.push(Token{Kind: kind, Value: match})
	} else {
		lex.push(Token{Kind: VALUE, Value: match})
	}
	lex.pos += len(match)
}

func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.input[lex.pos:])
	lex.push(Token{Kind: VALUE, Value: match})
	lex.pos += len(match)
}

func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.input[lex.pos:])
	lex.pos += len(match)
}

func (l *lexer) at_end() bool {
	return l.pos >= len(l.input)
}

func (l *lexer) push(t Token) {
	l.tokens = append(l.tokens, t)
}

func newLexer(input string) *lexer {
	return &lexer{
		input:  input,
		pos:    0,
		tokens: make([]Token, 0),
		regex: []regexPattern{
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
		},
	}
}
func Tokenize(input string) []Token {
	l := newLexer(strings.TrimSuffix(input, "\n")) // we know that a command ends with a newline
	for !l.at_end() {
		for _, p := range l.regex {
			loc := p.pattern.FindStringIndex(l.input[l.pos:])
			if loc != nil && loc[0] == 0 {
				p.handler(l, p.pattern)
				break
			}
		}
	}

	l.push(Token{Kind: NEW_LINE, Value: "nl"})
	return l.tokens
}
