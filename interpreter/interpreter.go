package interpreter

import (
	"fmt"

	"github.com/anwprath/glox/ast"
	"github.com/anwprath/glox/token"
)

var _ ast.Visitor = &Interpreter{}

type Interpreter struct{}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) any {
	panic("unimplemented")
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) any {
	return i.evaluate(expr)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) any {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {

	case token.BANG:
		return !isTruthy(right)

	case token.MINUS:
		val, ok := right.(float64)
		if !ok {
			panic(fmt.Sprintf("Operand must be a number, got %T", right))
		}
		return -1 * val
	}

	panic(fmt.Sprintf("Unknown unary operator: %s", expr.Operator))
}

func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}
	if val, ok := value.(bool); ok {
		return !val
	}
	return true
}
