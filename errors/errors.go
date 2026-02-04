package errors

import (
	"fmt"
	"strconv"
)

var HadError bool = false

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Println("[line " + strconv.Itoa(line) + "] Error" + where + ": " + message)
	HadError = true
}
