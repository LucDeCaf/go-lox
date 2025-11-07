package lox

import "fmt"

type Interpreter struct {
	errors []error
}

type InterpretError struct {
	stmt    Stmt
	message string
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		errors: []error{},
	}
}

func (i *Interpreter) VisitLiteralExpr(l *LiteralExpr) any {
	return l.value
}

func (i *Interpreter) VisitGroupingExpr(g *GroupingExpr) any {
	return i.evaluate(g.expression)
}

func (i *Interpreter) VisitBinaryExpr(b *BinaryExpr) any {
	left := i.evaluate(b.left)
	right := i.evaluate(b.right)

	switch b.operator.tokenType {
	case MINUS:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left - right

	case STAR:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left * right

	case SLASH:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left / right

	case PLUS:
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

	case GREATER:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left > right

	case GREATER_EQUAL:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left >= right

	case LESS:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left < right

	case LESS_EQUAL:
		left, right, ok := extractFloats(left, right)
		if !ok {
			return nil
		}

		return left <= right

	case BANG_EQUAL:
		return left != right

	case EQUAL_EQUAL:
		return left == right
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitUnaryExpr(u *UnaryExpr) any {
	right := i.evaluate(u.right)

	switch u.operator.tokenType {
	case MINUS:
		right, ok := right.(float64)
		if !ok {
			return nil
		}

		return -right

	case BANG:
		return !isTruthy(right)
	}

	// Unreachable
	// TODO proper error handling regardless
	return nil
}

func (i *Interpreter) VisitExpressionStmt(s *ExpressionStmt) {
	i.evaluate(s.expression)
}

func (i *Interpreter) VisitPrintStmt(s *PrintStmt) {
	expr := i.evaluate(s.expression)
	fmt.Printf("%v\n", expr)
}

func (i *Interpreter) Interpret(statements []Stmt) {
	for _, stmt := range statements {
		if err := i.execute(stmt); err != nil {
			i.errors = append(i.errors, err)
		}
	}
}

func (i *Interpreter) execute(s Stmt) error {
	return s.Accept(i)
}

func (i *Interpreter) evaluate(e Expr) any {
	return e.Accept(i)
}

func (i *Interpreter) addError(err error) {
	i.errors = append(i.errors, err)
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
