package main

import "fmt"

type Expr interface {
	Eval() Expr
	String() string
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

func (b Binary) Eval() Expr {
	// TODO
	return b
}

func (g Grouping) Eval() Expr {
	// TODO
	return g
}

func (u Unary) Eval() Expr {
	// TODO
	return u
}

func (l Literal) Eval() Expr {
	return l
}

func (b Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.left.String(), b.right.String(), b.operator.lexeme)
}

func (g Grouping) String() string {
	return fmt.Sprintf("(group %s)", g.expression.String())
}

func (u Unary) String() string {
	return fmt.Sprintf("(%s %s)", u.right.String(), u.operator.lexeme)
}

func (l Literal) String() string {
	switch v := l.value.(type) {
	case float64:
		return fmt.Sprintf("%v", v)
	case string:
		return "'" + v + "'"
	case bool:
		if v {
			return "<true>"
		} else {
			return "<false>"
		}
	case nil:
		return "<nil>"
	default:
		return "<invalid>"
	}
}
