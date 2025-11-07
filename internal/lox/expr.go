package lox

type Expr interface {
	Accept(ExprVisitor) any
}

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) any
	VisitLiteralExpr(*LiteralExpr) any
	VisitGroupingExpr(*GroupingExpr) any
	VisitUnaryExpr(*UnaryExpr) any
}

type LiteralExpr struct {
	value any
}

type BinaryExpr struct {
	left, right Expr
	operator    Token
}

type GroupingExpr struct {
	expression Expr
}

type UnaryExpr struct {
	right    Expr
	operator Token
}

func (expr *LiteralExpr) Accept(v ExprVisitor) any {
	return v.VisitLiteralExpr(expr)
}

func (expr *BinaryExpr) Accept(v ExprVisitor) any {
	return v.VisitBinaryExpr(expr)
}

func (expr *GroupingExpr) Accept(v ExprVisitor) any {
	return v.VisitGroupingExpr(expr)
}

func (expr *UnaryExpr) Accept(v ExprVisitor) any {
	return v.VisitUnaryExpr(expr)
}
