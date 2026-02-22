package errors

import (
	"fmt"
	"strconv"

	"github.com/anwprath/glox/token"
)

var HadError bool = false

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Println("[line " + strconv.Itoa(line) + "] Error" + where + ": " + message)
	HadError = true
}

func ReportParseError(t token.Token,  message string) error {
    if (t.TokenType == token.EOF) {
      report(t.Line, " at end", message);
    } else {
      report(t.Line, " at '" + t.Lexeme + "'", message);
    }

	return ParseErr{}
}

type ParseErr struct {} 

func (p ParseErr) Error() string {
	return ""
}	