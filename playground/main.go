package main

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/anwprath/glox/parser"
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
	tokenParser := parser.Parser{Tokens: tokens}
	exp, err := tokenParser.Parse()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf( "%s\n", exp)
	}

}

func runFile(args []string) {
	bytes, err := os.ReadFile(args[0])
	if err != nil {
		slog.Error("error reading file", "file", args[0], "error", err)
		log.Fatal()
	}

	run(string(bytes))
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
	}
}
