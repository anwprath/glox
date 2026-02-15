package ast

import "github.com/anwprath/glox/token"

type Expr interface {
	Accept(v Visitor) any
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (node *Binary) Accept(v Visitor) any {
	return v.VisitBinaryExpr(node)
}

type Grouping struct {
	Expression Expr
}

func (node *Grouping) Accept(v Visitor) any {
	return v.VisitGroupingExpr(node)
}

type Literal struct {
	Value any
}

func (node *Literal) Accept(v Visitor) any {
	return v.VisitLiteralExpr(node)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (node *Unary) Accept(v Visitor) any {
	return v.VisitUnaryExpr(node)
}

type Visitor interface {
	VisitBinaryExpr(expr *Binary) any
	VisitGroupingExpr(expr *Grouping) any
	VisitLiteralExpr(expr *Literal) any
	VisitUnaryExpr(expr *Unary) any
}
