package lox

import (
	"fmt"
	"github.com/LucDeCaf/go-lox/internal/lox/ast"
)

type AstPrinter struct{}

func (a *AstPrinter) Print(e ast.Expr) string {
	switch v := e.Accept(a).(type) {
	case string:
		return v
	default:
		return "<invalid type>"
	}
}

func (a *AstPrinter) VisitLiteralExpr(l *ast.LiteralExpr) any {
	switch v := l.Value.(type) {
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

func (a *AstPrinter) VisitBinaryExpr(b *ast.BinaryExpr) any {
	return fmt.Sprintf("(%s %s %s)", b.Left.Accept(a), b.Right.Accept(a), b.Operator.Lexeme)
}

func (a *AstPrinter) VisitGroupingExpr(g *ast.GroupingExpr) any {
	return fmt.Sprintf("(group %s)", g.Expression.Accept(a))
}

func (a *AstPrinter) VisitUnaryExpr(u *ast.UnaryExpr) any {
	return fmt.Sprintf("(%s %s)", u.Right.Accept(a), u.Operator.Lexeme)
}

func (a *AstPrinter) VisitVariableExpr(u *ast.VariableExpr) any {
	return u.Name.Lexeme
}
