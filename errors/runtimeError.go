package errors

import (
	"fmt"
	"strconv"

	"github.com/anwprath/glox/token"
)

type RuntimeError struct {
	Token   token.Token
	Message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("runtime error at token %v: %s", e.Token, e.Message)
}

func ReportRuntimeError(err RuntimeError) {
	fmt.Println(err.Error() + "\n[line " + strconv.Itoa(err.Token.Line) + "]")
	HadRuntimeError = true
}
