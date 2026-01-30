package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	errors "github.com/anwprath/glox/errors"
)

func main() {
	args := os.Args

	if len(args) > 2 {
		slog.Info("Usage: glox [script]")
	} else if len(args) == 2 {
		runFile(args)
	} else {
		runPrompt()
	}
}

func run(command string) {
	fmt.Println(command)

}

func runFile(args []string) {
	bytes, err := os.ReadFile(args[0])
	if err != nil {
		slog.Error("error reading file", "file", args[0], "error", err)
		log.Fatal()
	}

	run(string(bytes))
	if errors.HadError {
		os.Exit(69)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		// Ctrl+D works here
		if !scanner.Scan() {
			break
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			break
		}
		run(scanner.Text())
		errors.HadError = false
	}
}
