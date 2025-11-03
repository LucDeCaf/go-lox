package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	hadError bool
}

func NewLox() Lox {
	return Lox{
		hadError: false,
	}
}

func (l *Lox) runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		if line == "\n" {
			break
		}
		l.run(line)
		l.hadError = false
	}
}

func (l *Lox) runFile(path string) (err error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	l.run(string(source))

	return nil
}

func (l *Lox) run(source string) {
	scanner := NewScanner()
	tokens := scanner.scanTokens(source)

	for _, token := range tokens {
		fmt.Printf("%s\n", token.toString())
	}
}

func (l *Lox) reportError(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report() {

}
