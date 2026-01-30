package scanner

import (
	"github.com/anwprath/glox/token"
)

type Scanner struct {
	source               string
	tokens               []token.Token
	start, current, line int
}

func New(source string) Scanner {
	return Scanner{
		source:  source,
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s Scanner) scanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.New(token.EOF, "", nil, 999))
	return s.tokens
}

func (s Scanner) scanToken() {

}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
