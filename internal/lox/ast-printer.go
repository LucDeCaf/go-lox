package lox

import (
	"fmt"
)

type AstPrinter struct{}

func (a *AstPrinter) Print(e Expr) string {
	switch v := e.accept(a).(type) {
	case string:
		return v
	default:
		return "<invalid type>"
	}
}

func (a *AstPrinter) visitLiteral(l *Literal) any {
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

func (a *AstPrinter) visitBinary(b *Binary) any {
	return fmt.Sprintf("(%s %s %s)", b.left.accept(a), b.right.accept(a), b.operator.lexeme)
}

func (a *AstPrinter) visitGrouping(g *Grouping) any {
	return fmt.Sprintf("(group %s)", g.expression.accept(a))
}

func (a *AstPrinter) visitUnary(u *Unary) any {
	return fmt.Sprintf("(%s %s)", u.right.accept(a), u.operator.lexeme)
}
