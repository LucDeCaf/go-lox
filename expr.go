package main

type Expr interface {
	Accept(ExprVisitor) any
}

type ExprVisitor interface {
	VisitBinary(*Binary) any
	VisitLiteral(*Literal) any
	VisitGrouping(*Grouping) any
	VisitUnary(*Unary) any
}

type Literal struct {
	value any
}

type Binary struct {
	left, right Expr
	operator    Token
}

type Grouping struct {
	expression Expr
}

type Unary struct {
	right    Expr
	operator Token
}

func (expr *Literal) Accept(v ExprVisitor) any {
	return v.VisitLiteral(expr)
}

func (expr *Binary) Accept(v ExprVisitor) any {
	return v.VisitBinary(expr)
}

func (expr *Grouping) Accept(v ExprVisitor) any {
	return v.VisitGrouping(expr)
}

func (expr *Unary) Accept(v ExprVisitor) any {
	return v.VisitUnary(expr)
}
