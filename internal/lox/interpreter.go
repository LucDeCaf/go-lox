package lox

type interpreter struct{}

func newInterpreter() *interpreter {
	return &interpreter{}
}

func (i *interpreter) visitLiteral(l *Literal) any {
	return l.value
}

func (i *interpreter) visitGrouping(g *Grouping) any {
	return i.evaluate(g.expression)
}

func (i *interpreter) visitBinary(b *Binary) any {
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

func (i *interpreter) visitUnary(u *Unary) any {
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

func (i *interpreter) Interpret(e Expr) any {
	// TODO error handling/reporting
	return i.evaluate(e)
}

func (i *interpreter) evaluate(e Expr) any {
	return e.accept(i)
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
