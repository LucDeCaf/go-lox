package lox

import (
	"fmt"
	"github.com/LucDeCaf/go-lox/internal/lox/ast"
	"strconv"
)

type Scanner struct {
	source string
	tokens []ast.Token
	errors []error

	start   int
	current int
	line    int
}

type ScanError struct {
	line    int
	message string
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("[line %d] Error: %s", e.line, e.message)
}

func NewScanner() *Scanner {
	return &Scanner{
		source:  "",
		tokens:  []ast.Token{},
		errors:  []error{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) scanTokens(source string) ([]ast.Token, bool) {
	s.source = source
	s.errors = []error{}

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, ast.NewToken(
		ast.EOF,
		"",
		nil,
		s.line,
	))

	return s.tokens, len(s.errors) == 0
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(ast.LEFT_PAREN)
	case ')':
		s.addToken(ast.RIGHT_PAREN)
	case '{':
		s.addToken(ast.LEFT_BRACE)
	case '}':
		s.addToken(ast.RIGHT_BRACE)
	case ',':
		s.addToken(ast.COMMA)
	case '.':
		s.addToken(ast.DOT)
	case '-':
		s.addToken(ast.MINUS)
	case '+':
		s.addToken(ast.PLUS)
	case ';':
		s.addToken(ast.SEMICOLON)
	case '*':
		s.addToken(ast.STAR)
	case '!':
		if s.match('=') {
			s.addToken(ast.BANG_EQUAL)
		} else {
			s.addToken(ast.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(ast.EQUAL_EQUAL)
		} else {
			s.addToken(ast.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(ast.LESS_EQUAL)
		} else {
			s.addToken(ast.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(ast.GREATER_EQUAL)
		} else {
			s.addToken(ast.GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(ast.SLASH)
		}

	case ' ':
		fallthrough
	case '\r':
		fallthrough
	case '\t':
		break

	case '\n':
		s.line++

	case '"':
		if err := s.parseString(); err != nil {
			s.errors = append(s.errors, err)
		}

	default:
		if isDigit(c) {
			s.parseNumber()
		} else if isAlpha(c) {
			s.parseIdent()
		} else {
			s.errors = append(s.errors, &ScanError{
				line:    s.line,
				message: "Unexpected character.",
			})
		}
	}
}

func (s *Scanner) addToken(tokenType ast.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType ast.TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, ast.NewToken(tokenType, lexeme, literal, s.line))
}

func (s *Scanner) parseString() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return &ScanError{
			line:    s.line,
			message: "Unterminated string.",
		}
	}

	s.advance()

	literal := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(ast.STRING, literal)

	return nil
}

func (s *Scanner) parseNumber() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	lexeme := s.source[s.start:s.current]
	number, _ := strconv.ParseFloat(lexeme, 64)
	s.addTokenWithLiteral(ast.NUMBER, number)
}

func (s *Scanner) parseIdent() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	ident := s.source[s.start:s.current]
	keyword, ok := ast.Keywords[ident]
	if !ok {
		s.addTokenWithLiteral(ast.IDENTIFIER, ident)
	} else {
		s.addToken(keyword)
	}
}

func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current++
	return c
}

func (s *Scanner) match(c byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != c {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
