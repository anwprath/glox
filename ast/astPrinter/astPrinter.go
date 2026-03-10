package astPrinter

import (
	"fmt"
	"strings"

	"github.com/anwprath/glox/ast"
)

// Fails to compile if `AstPrinter` does not implement ast.Visitor
var _ ast.ExprVisitor = &AstPrinter{}

type AstPrinter struct{}

func (p *AstPrinter) Print(expr ast.Expr) (any, error) {
	return expr.Accept(p)
}

func (p *AstPrinter) VisitBinaryExpr(expr *ast.Binary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme,
		expr.Left, expr.Right), nil
}

func (p *AstPrinter) VisitGroupingExpr(expr *ast.Grouping) (any, error) {
	return p.parenthesize("group", expr.Expression), nil
}

func (p *AstPrinter) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return expr.Value, nil
}

func (p *AstPrinter) VisitUnaryExpr(expr *ast.Unary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (p *AstPrinter) parenthesize(name string, exprs ...ast.Expr) string {
	builder := strings.Builder{}
	builder.WriteString("(")
	builder.WriteString(name)

	for _, exp := range exprs {
		builder.WriteString(" ")
		if exp != nil {
			builder.WriteString(fmt.Sprint(exp.Accept(p)))
		} else {
			builder.WriteString("nil")
		}
	}
	builder.WriteString(")")

	return builder.String()
}
