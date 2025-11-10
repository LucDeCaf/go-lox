package ast

type Expr interface {
	Accept(ExprVisitor) any
}

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) any
	VisitLiteralExpr(*LiteralExpr) any
	VisitGroupingExpr(*GroupingExpr) any
	VisitUnaryExpr(*UnaryExpr) any
	VisitVariableExpr(*VariableExpr) any
}

type LiteralExpr struct {
	Value any
}

type BinaryExpr struct {
	Left, Right Expr
	Operator    *Token
}

type GroupingExpr struct {
	Expression Expr
}

type UnaryExpr struct {
	Right    Expr
	Operator *Token
}

type VariableExpr struct {
	Name *Token
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

func (expr *VariableExpr) Accept(v ExprVisitor) any {
	return v.VisitVariableExpr(expr)
}
