package scanner

import (
	"strconv"

	"github.com/anwprath/glox/errors"
	"github.com/anwprath/glox/token"
)

type Scanner struct {
	source               []rune
	tokens               []token.Token
	start, current, line int
}

func New(source string) Scanner {
	return Scanner{
		source:  []rune(source),
		tokens:  make([]token.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.New(token.EOF, "", nil, 999))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH)
		}
	case ' ':
		break
	case '\r':
		break
	case '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.scanString()
	default:
		if isDigit(c) {
			s.scanNumber()
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			errors.Error(s.line, "Unexpected character.")
		}
	}

}

func (s *Scanner) addToken(token token.TokenType) {
	s.appendToken(token, nil)
}

func (s *Scanner) appendToken(tokenType token.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.New(tokenType, string(text), literal, s.line))
}

func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current++
	return c
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) scanString() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		errors.Error(s.line, "Unterminated string")
		return
	}

	// The closing ".
	s.advance()

	var value string = string(s.source[s.start+1 : s.current-1])
	s.appendToken(token.STRING, value)
}

func (s *Scanner) scanNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	var valueStr string = string(s.source[s.start:s.current])
	valueNum, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		panic(err)
	}
	s.appendToken(token.NUMBER, valueNum)
}

func (s *Scanner) scanIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	value := string(s.source[s.start:s.current])
	if t, ok := token.GetKeyword(value); !ok {
		s.addToken(token.IDENTIFIER)
	} else {
		s.addToken(t)
	}
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}
