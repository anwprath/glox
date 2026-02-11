package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	errors "github.com/anwprath/glox/errors"
	"github.com/anwprath/glox/scanner"
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
	sc := scanner.New(command)
	tokens := sc.ScanTokens()
	fmt.Println(tokens)
}

func runFile(args []string) {
	bytes, err := os.ReadFile(args[1])
	if err != nil {
		slog.Error("error reading file", "file", args[1], "error", err)
		log.Fatal()
	}
	fmt.Println("Reading file: ", args[1])
	run(string(bytes))
	if errors.HadError {
		os.Exit(69)
	}
}

func runPrompt() {
	inputScanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		// Ctrl+D works here
		if !inputScanner.Scan() {
			break
		}
		text := strings.TrimSpace(inputScanner.Text())
		if text == "" {
			break
		}
		run(inputScanner.Text())
		errors.HadError = false
	}
}
