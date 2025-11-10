package lox

import (
	"fmt"
	"github.com/LucDeCaf/go-lox/internal/lox/ast"
)

type Interpreter struct {
	env    *Environment
	errors []error
}

type RuntimeError struct {
	message string
}

func (e *RuntimeError) Error() string {
	return "RuntimeError: " + e.message
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env:    NewEnvironment(),
		errors: []error{},
	}
}

func (i *Interpreter) VisitLiteralExpr(l *ast.LiteralExpr) any {
	return l.Value
}

func (i *Interpreter) VisitGroupingExpr(g *ast.GroupingExpr) any {
	return i.evaluate(g.Expression)
}

func (i *Interpreter) VisitBinaryExpr(b *ast.BinaryExpr) any {
	left := i.evaluate(b.Left)
	right := i.evaluate(b.Right)

	switch b.Operator.Type {
	case ast.MINUS:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left - right

	case ast.STAR:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left * right

	case ast.SLASH:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left / right

	case ast.PLUS:
		switch left := left.(type) {
		case float64:
			right, ok := right.(float64)
			if !ok {
				return nil
			}
			return left + right
		case string:
			right, ok := right.(string)
			if !ok {
				return nil
			}
			return left + right
		}

	case ast.GREATER:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left > right

	case ast.GREATER_EQUAL:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left >= right

	case ast.LESS:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left < right

	case ast.LESS_EQUAL:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left <= right

	case ast.BANG_EQUAL:
		return left != right

	case ast.EQUAL_EQUAL:
		return left == right
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitUnaryExpr(u *ast.UnaryExpr) any {
	right := i.evaluate(u.Right)

	switch u.Operator.Type {
	case ast.MINUS:
		right, ok := right.(float64)
		if !ok {
			return nil
		}

		return -right

	case ast.BANG:
		return !isTruthy(right)
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitVariableExpr(u *ast.VariableExpr) any {
	v, ok := i.env.Get(u.Name.Lexeme)
	if !ok {
		return &RuntimeError{
			message: "Undefined variable '" + u.Name.Lexeme + "'",
		}
	}
	return v
}

func (i *Interpreter) VisitExpressionStmt(s *ast.ExpressionStmt) error {
	i.evaluate(s.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(s *ast.PrintStmt) error {
	expr := i.evaluate(s.Expression)
	fmt.Printf("%v\n", expr)
	return nil
}

func (i *Interpreter) VisitVarStmt(s *ast.VarStmt) error {
	var value any = nil
	if s.Value != nil {
		value = i.evaluate(s.Value)
	}

	i.env.Define(s.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) Interpret(statements []ast.Stmt) {
	for _, stmt := range statements {
		if err := i.execute(stmt); err != nil {
			i.errors = append(i.errors, err)
		}
	}
}

func (i *Interpreter) execute(s ast.Stmt) error {
	return s.Accept(i)
}

func (i *Interpreter) evaluate(e ast.Expr) any {
	return e.Accept(i)
}

func isTruthy(obj any) bool {
	switch v := obj.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}

func extractFloats(a, b any) (float64, float64, bool) {
	aF, ok := a.(float64)
	if !ok {
		return 0, 0, false
	}
	bF, ok := b.(float64)
	if !ok {
		return 0, 0, false
	}
	return aF, bF, true
}
