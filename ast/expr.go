package ast

import "github.com/anwprath/glox/token"

type Expr interface {
	Accept(v Visitor) any
}

type Binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

func (node *Binary) Accept(v Visitor) any {
	return v.VisitBinaryExpr(node)
}

type Grouping struct {
	expression Expr
}

func (node *Grouping) Accept(v Visitor) any {
	return v.VisitGroupingExpr(node)
}

type Literal struct {
	value any
}

func (node *Literal) Accept(v Visitor) any {
	return v.VisitLiteralExpr(node)
}

type Unary struct {
	operator token.Token
	right    Expr
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
