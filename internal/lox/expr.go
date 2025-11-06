package lox

type Expr interface {
	accept(exprVisitor) any
}

type exprVisitor interface {
	visitBinary(*Binary) any
	visitLiteral(*Literal) any
	visitGrouping(*Grouping) any
	visitUnary(*Unary) any
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

func (expr *Literal) accept(v exprVisitor) any {
	return v.visitLiteral(expr)
}

func (expr *Binary) accept(v exprVisitor) any {
	return v.visitBinary(expr)
}

func (expr *Grouping) accept(v exprVisitor) any {
	return v.visitGrouping(expr)
}

func (expr *Unary) accept(v exprVisitor) any {
	return v.visitUnary(expr)
}
