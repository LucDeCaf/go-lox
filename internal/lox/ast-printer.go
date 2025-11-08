package lox

import (
	"fmt"
)

type AstPrinter struct{}

func (a *AstPrinter) Print(e Expr) string {
	switch v := e.Accept(a).(type) {
	case string:
		return v
	default:
		return "<invalid type>"
	}
}

func (a *AstPrinter) VisitLiteralExpr(l *LiteralExpr) any {
	switch v := l.value.(type) {
	case float64:
		return fmt.Sprintf("%v", v)
	case string:
		return "'" + v + "'"
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	case nil:
		return "nil"
	default:
		return "<invalid literal>"
	}
}

func (a *AstPrinter) VisitBinaryExpr(b *BinaryExpr) any {
	return fmt.Sprintf("(%s %s %s)", b.left.Accept(a), b.right.Accept(a), b.operator.lexeme)
}

func (a *AstPrinter) VisitGroupingExpr(g *GroupingExpr) any {
	return fmt.Sprintf("(group %s)", g.expression.Accept(a))
}

func (a *AstPrinter) VisitUnaryExpr(u *UnaryExpr) any {
	return fmt.Sprintf("(%s %s)", u.right.Accept(a), u.operator.lexeme)
}

func (a *AstPrinter) VisitVariableExpr(u *VariableExpr) any {
	return u.name.lexeme
}
