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

	if l.hadError {
		os.Exit(65)
	}

	return nil
}

func (l *Lox) run(source string) {
	scanner := NewScanner(l)
	tokens := scanner.scanTokens(source)

	for _, token := range tokens {
		fmt.Printf("%s\n", token.String())
	}
}

func (l *Lox) reportError(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	l.hadError = true
}
