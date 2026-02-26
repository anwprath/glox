package interpreter

import (
	"fmt"
	"reflect"

	"github.com/anwprath/glox/ast"
	"github.com/anwprath/glox/errors"
	"github.com/anwprath/glox/token"
)

var _ ast.Visitor = &Interpreter{}

type Interpreter struct{}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case token.EQUAL_EQUAL:
		return isEqual(left, right), nil
	case token.BANG_EQUAL:
		return !isEqual(left, right), nil
	}

	err = checkBinaryNumberOperand(expr.Operator, left, right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.TokenType {

	case token.LESS:
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		return left.(float64) <= right.(float64), nil
	case token.GREATER:
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		return left.(float64) >= right.(float64), nil
	case token.MINUS:
		return left.(float64) - right.(float64), nil
	case token.STAR:
		return left.(float64) * right.(float64), nil
	case token.SLASH:
		return left.(float64) / right.(float64), nil
	case token.PLUS:
		if l, okL := left.(float64); okL {
			if r, okR := right.(float64); okR {
				return l + r, nil
			}
		}
		if l, okL := left.(string); okL {
			if r, okR := right.(string); okR {
				return l + r, nil
			}
		}
		return nil, errors.RuntimeError{
			Token:   expr.Operator,
			Message: "Both operands must be either string or number.",
		}
	}

	panic("unreachable code in VisitBinaryExpr")
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) (any, error) {
	return i.evaluate(expr)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) (any, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.TokenType {
	case token.BANG:
		return !isTruthy(right), nil
	case token.MINUS:
		err := checkUnaryNumberOperand(expr.Operator, expr.Right)
		if err != nil {
			return nil, err
		}
		val, ok := right.(float64)
		if !ok {
			panic(fmt.Sprintf("Operand must be a number, got %T", right))
		}
		return -1 * val, nil
	}

	panic(fmt.Sprintf("Unknown unary operator: %s", expr.Operator))
}

func (i *Interpreter) evaluate(expr ast.Expr) (any, error) {
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

func checkUnaryNumberOperand(operator token.Token, operand any) error {
	if _, ok := operand.(float64); !ok {
		return errors.RuntimeError{
			Token:   operator,
			Message: "Operand must be a number.",
		}
	}
	return nil
}

func checkBinaryNumberOperand(operator token.Token, leftOperand, rightOperand any) error {
	if _, ok := leftOperand.(float64); !ok {
		return errors.RuntimeError{
			Token:   operator,
			Message: "Left operand must be a number.",
		}
	}
	if _, ok := rightOperand.(float64); !ok {
		return errors.RuntimeError{
			Token:   operator,
			Message: "Right operand must be a number.",
		}
	}
	return nil
}
