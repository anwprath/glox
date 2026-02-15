package token

import "fmt"

type Token struct {
	tokenType TokenType
	Lexeme    string
	literal   any
	line      int
}

func New(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		tokenType: tokenType,
		Lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

func (t Token) String() string {
	return "Token(" + t.tokenType.String() + " " + t.Lexeme + " " + fmt.Sprintf("%v", t.literal) + ")"
}

var keywords map[string]TokenType = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func GetKeyword(key string) (TokenType, bool) {
	tok, ok := keywords[key]
	return tok, ok
}
