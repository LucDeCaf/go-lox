package main

type Scanner struct {
	tokens []Token
}

func NewScanner() Scanner {
	return Scanner{
		tokens: []Token{},
	}
}

func (s *Scanner) scanTokens(source string) []Token {

	return s.tokens
}
