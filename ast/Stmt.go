package ast

import "github.com/anwprath/glox/token"

type Stmt interface {
	Accept(v StmtVisitor) (any, error)
}

type StmtVisitor interface {
	VisitExpressionStmt(expr *Expression) (any, error)
	VisitPrintStmt(expr *Print) (any, error)
}

type Expression struct {
	expression Expr
}

func (node *Expression) Accept(v StmtVisitor) (any, error) {
	return v.VisitExpressionExpr(node)
}

type Print struct {
	expression Expr
}

func (node *Print) Accept(v StmtVisitor) (any, error) {
	return v.VisitPrintExpr(node)
}
