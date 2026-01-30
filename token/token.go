package token

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

func New(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

func (t Token) String() string {
	return t.tokenType.String() + " " + t.lexeme + " " + fmt.Sprintf("%v", t.literal)
}
