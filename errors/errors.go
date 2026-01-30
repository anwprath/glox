package errors

import (
	"log/slog"
	"strconv"
)

var HadError bool = false

func Error(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	slog.Error("[line " + strconv.Itoa(line) + "] Error" + where + ": " + message)
	HadError = true
}
