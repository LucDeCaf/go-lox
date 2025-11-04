package main

import (
	"fmt"
)

type AstPrinter struct {
	astString string
}

func (a *AstPrinter) Print(e Expr) string {
	switch v := e.Accept(a).(type) {
	case string:
		return v
	default:
		return "<invalid type>"
	}
}

func (a *AstPrinter) VisitLiteral(l *Literal) any {
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

func (a *AstPrinter) VisitBinary(b *Binary) any {
	return fmt.Sprintf("(%s %s %s)", b.left.Accept(a), b.right.Accept(a), b.operator.lexeme)
}

func (a *AstPrinter) VisitGrouping(g *Grouping) any {
	return fmt.Sprintf("(group %s)", g.expression.Accept(a))
}

func (a *AstPrinter) VisitUnary(u *Unary) any {
	return fmt.Sprintf("(%s %s)", u.right.Accept(a), u.operator.lexeme)
}
