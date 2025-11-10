package lox

import (
	"fmt"
	"github.com/LucDeCaf/go-lox/internal/lox/ast"
)

type Parser struct {
	tokens  []ast.Token
	errors  []error
	current int
}

type ParseError struct {
	token   ast.Token
	message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("ParseError <%s>: %s", e.token.String(), e.message)
}

func NewParser() *Parser {
	return &Parser{
		tokens:  []ast.Token{},
		errors:  []error{},
		current: 0,
	}
}
func (p *Parser) parse(tokens []ast.Token) ([]ast.Stmt, bool) {
	p.tokens = tokens
	p.errors = []error{}

	statements := []ast.Stmt{}

	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.errors = append(p.errors, err)
			p.synchronize()
		} else {
			statements = append(statements, stmt)
		}
	}

	if len(p.errors) > 0 {
		return nil, false
	}

	return statements, true
}

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(ast.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(ast.IDENTIFIER, "Expect identifier after VAR.")
	if err != nil {
		return nil, err
	}

	var value ast.Expr = nil
	if p.match(ast.EQUAL) {
		value, err = p.expression()
	}

	return &ast.VarStmt{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(ast.PRINT) {
		return p.printStmt()
	}

	return p.expressionStmt()
}

func (p *Parser) printStmt() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(ast.SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return &ast.PrintStmt{Expression: expr}, nil
}

func (p *Parser) expressionStmt() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(ast.SEMICOLON, "Expect ';'; after expression.")
	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStmt{Expression: expr}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(ast.BANG_EQUAL, ast.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(ast.GREATER, ast.GREATER_EQUAL, ast.LESS, ast.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(ast.PLUS, ast.MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(ast.STAR, ast.SLASH) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(ast.BANG, ast.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(ast.TRUE) {
		return &ast.LiteralExpr{Value: true}, nil
	}
	if p.match(ast.FALSE) {
		return &ast.LiteralExpr{Value: false}, nil
	}
	if p.match(ast.NIL) {
		return &ast.LiteralExpr{Value: nil}, nil
	}

	if p.match(ast.IDENTIFIER) {
		return &ast.VariableExpr{Name: p.previous()}, nil
	}

	if p.match(ast.NUMBER, ast.STRING) {
		value := p.previous().Literal
		return &ast.LiteralExpr{Value: value}, nil
	}

	if p.match(ast.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(ast.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return &ast.GroupingExpr{Expression: expr}, nil
	}

	return nil, &ParseError{token: *p.peek(), message: "Invalid token."}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == ast.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case ast.CLASS:
			return
		case ast.FUN:
			return
		case ast.VAR:
			return
		case ast.FOR:
			return
		case ast.IF:
			return
		case ast.WHILE:
			return
		case ast.PRINT:
			return
		case ast.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) advance() *ast.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() *ast.Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) match(tokenTypes ...ast.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType ast.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) peek() *ast.Token {
	return &p.tokens[p.current]
}

func (p *Parser) consume(tokenType ast.TokenType, errorMessage string) (*ast.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	err := &ParseError{token: *p.peek(), message: errorMessage}
	return nil, err
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.peek().Type == ast.EOF
}
