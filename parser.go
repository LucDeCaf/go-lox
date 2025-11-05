package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	errors  []error
	current int
}

type ParseError struct {
	token   Token
	message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("ParseError <%s>: %s", e.token.String(), e.message)
}

func NewParser() Parser {
	return Parser{
		current: 0,
	}
}

func (p *Parser) parse(tokens []Token) (Expr, bool) {
	p.tokens = tokens
	p.errors = []error{}

	expr, err := p.expression()
	if err != nil {
		p.errors = append(p.errors, err)
	}

	if len(p.errors) > 0 {
		return nil, false
	}

	return expr, true
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &Binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = &Binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(PLUS, MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &Binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(STAR, SLASH) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = &Binary{
			left:     expr,
			operator: *operator,
			right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &Unary{
			operator: *operator,
			right:    right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(TRUE) {
		return &Literal{value: true}, nil
	}
	if p.match(FALSE) {
		return &Literal{value: false}, nil
	}
	if p.match(NIL) {
		return &Literal{value: nil}, nil
	}

	if p.match(NUMBER, STRING) {
		value := p.previous().literal
		return &Literal{value}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}

		return &Grouping{expression: expr}, nil
	}

	return nil, &ParseError{token: *p.peek(), message: "Invalid token in expr"}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS:
			return
		case FUN:
			return
		case VAR:
			return
		case FOR:
			return
		case IF:
			return
		case WHILE:
			return
		case PRINT:
			return
		case RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) addError(err *ParseError) {
	return
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) peek() *Token {
	return &p.tokens[p.current]
}

func (p *Parser) consume(tokenType TokenType, errorMessage string) (*Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	err := &ParseError{token: *p.peek(), message: errorMessage}
	return nil, err
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.peek().tokenType == EOF
}
