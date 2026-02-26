package errors

import (
	"fmt"

	"github.com/anwprath/glox/token"
)

type RuntimeError struct {
	Token token.Token
	Message    string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("runtime error at token %v: %s", e.Token, e.Message)
}
