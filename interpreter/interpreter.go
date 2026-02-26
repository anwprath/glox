package interpreter

import (
	"fmt"
	"reflect"

	"github.com/anwprath/glox/ast"
	"github.com/anwprath/glox/token"
)

var _ ast.Visitor = &Interpreter{}

type Interpreter struct{}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case token.EQUAL_EQUAL:
		return isEqual(left, right)
	case token.BANG_EQUAL:
		return !isEqual(left, right)
	case token.LESS:
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case token.GREATER:
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case token.MINUS:
		return left.(float64) - right.(float64)
	case token.STAR:
		return left.(float64) * right.(float64)
	case token.SLASH:
		return left.(float64) / right.(float64)
	case token.PLUS:
		if l, okL := left.(float64); okL {
			if r, okR := right.(float64); okR {
				return l + r
			}
		}
		if l, okL := left.(string); okL {
			if r, okR := right.(string); okR {
				return l + r
			}
		}
	}

	panic("unreachable code in VisitBinaryExpr")
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

func isEqual(a, b any) bool {
	// a == b panics if a and b are non-comparables
	return reflect.DeepEqual(a, b)
}
