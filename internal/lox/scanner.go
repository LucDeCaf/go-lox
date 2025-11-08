package lox

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source string
	tokens []Token
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
		tokens:  []Token{},
		errors:  []error{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) scanTokens(source string) ([]Token, bool) {
	s.source = source
	s.errors = []error{}

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(
		EOF,
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
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
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

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, lexeme, literal, s.line))
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
	s.addTokenWithLiteral(STRING, literal)

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
	s.addTokenWithLiteral(NUMBER, number)
}

func (s *Scanner) parseIdent() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	ident := s.source[s.start:s.current]
	keyword, ok := keywords[ident]
	if !ok {
		s.addTokenWithLiteral(IDENTIFIER, ident)
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
