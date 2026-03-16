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
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(69)
	}

	outputDir := args[1]

	defineAst(outputDir, "Expr", []string{
		"Binary   : Expr Left, token.Token Operator, Expr Right",
		"Grouping : Expr Expression",
		"Literal  : any Value",
		"Unary    : token.Token Operator, Expr Right",
	})

	defineAst(outputDir, "Stmt", []string{
		"Expression : Expr Expression",
		"Print      : Expr Expression",
	})
}

func defineAst(outputDir, baseName string, types []string) {
	defer func() { exec.Command("gofmt", "-w", "ast/"+baseName+".go").Run() }()
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
	fmt.Fprintf(w, "type %s interface{\nAccept(v %sVisitor) (any, error)\n}\n\n", baseName, baseName)

	defineVisitorInterface(w, baseName, types)

	for _, t := range types {
		structName := strings.TrimSpace(strings.Split(t, ":")[0])
		fields := strings.TrimSpace(strings.Split(t, ":")[1])
		defineType(w, baseName, structName, fields)
	}

	err = w.Flush()
	if err != nil {
		log.Fatal("Failed to flush:", err)
	}

}

func defineType(w *bufio.Writer, baseName, className, fieldList string) {
	fmt.Fprintf(w, "type %s struct {\n", className)

	fields := strings.Split(fieldList, ", ")
	for _, member := range fields {
		memberName := strings.Split(member, " ")[1]
		memberType := strings.Split(member, " ")[0]
		fmt.Fprintf(w, " %s %s\n", memberName, memberType)
	}
	fmt.Fprintf(w, "}\n\n")
	fmt.Fprintf(w, "func (node *%s)Accept(v %sVisitor) (any, error) {\n return v.Visit%sExpr(node)}\n", className, baseName, className)

}

func defineVisitorInterface(w *bufio.Writer, baseName string, types []string) {
	fmt.Fprintf(w, "type %sVisitor interface {\n", baseName)

	for _, t := range types {
		className := strings.TrimSpace(strings.Split(t, ":")[0])
		fmt.Fprintf(w, "Visit%s%s(expr *%s) (any, error)\n", className, baseName, className)
	}
	fmt.Fprintf(w, "}\n\n")

}
