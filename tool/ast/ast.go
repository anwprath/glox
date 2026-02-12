package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	defer func() { exec.Command("gofmt", "-w", "ast/expr.go").Run() }()
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(69)
	}

	outputDir := args[1]

	defineAst(outputDir, "expr", []string{
		"Binary   : Expr left, token.Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal  : any value",
		"Unary    : token.Token operator, Expr right",
	})
}

func defineAst(outputDir, baseName string, types []string) {
	outPath := fmt.Sprintf("%s/%s.go", outputDir, baseName)

	// Create directory
	os.MkdirAll(filepath.Dir(outPath), 0755)

	// Open file for writing (creates if doesn't exist)
	exprGo, err := os.Create(outPath)
	if err != nil {
		log.Fatal("Failed to create file:", err)
	}
	defer exprGo.Close()

	w := bufio.NewWriter(exprGo)
	fmt.Fprintf(w, "package ast\n\n")
	fmt.Fprintf(w, "import \"github.com/anwprath/glox/token\"\n\n")
	fmt.Fprintf(w, "type Expr interface{\nAccept(v Visitor) any\n}\n\n")

	for _, t := range types {
		structName := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(w, structName, fields)
	}

	defineVisitorInterface(w, types)

	err = w.Flush()
	if err != nil {
		log.Fatal("Failed to flush:", err)
	}

}

func defineType(w *bufio.Writer, className, fieldList string) {
	fmt.Fprintf(w, "type %s struct {\n", className)

	fields := strings.Split(fieldList, ", ")
	for _, member := range fields {
		memberName := strings.Split(member, " ")[1]
		memberType := strings.Split(member, " ")[0]
		fmt.Fprintf(w, " %s %s\n", memberName, memberType)
	}
	fmt.Fprintf(w, "}\n\n")
	fmt.Fprintf(w, "func (node *%s)Accept(v Visitor) any {\n return v.Visit%sExpr(node)}\n", className, className)

}

func defineVisitorInterface(w *bufio.Writer, types []string) {
	fmt.Fprintf(w, "type Visitor interface {\n")

	for _, t := range types {
		className := strings.TrimSpace(strings.Split(t, ":")[0])
		fmt.Fprintf(w, "Visit%sExpr(expr *%s) any\n", className, className)
	}
	fmt.Fprintf(w, "}\n\n")

}
