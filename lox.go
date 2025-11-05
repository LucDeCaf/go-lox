package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	interpreter *Interpreter
	reporters   []ErrorReporter[error]
	hadError    bool
}

func NewLox() Lox {
	return Lox{
		interpreter: NewInterpreter(),
		reporters:   []ErrorReporter[error]{},
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
	scanner := NewScanner()
	tokens, scanOk := scanner.scanTokens(source)
	for _, err := range scanner.errors {
		for _, r := range l.reporters {
			r.ReportError(err)
		}
	}

	parser := NewParser()
	ast, parseOk := parser.parse(tokens)
	for _, err := range parser.errors {
		for _, r := range l.reporters {
			r.ReportError(err)
		}
	}

	if !scanOk || !parseOk {
		return
	}

	v := l.interpreter.Interpret(ast)
	if v == nil {
		fmt.Println("WARNING: An error likely occured")
	}
	fmt.Printf("%v\n", v)
}

func (l *Lox) registerErrorReporter(r ErrorReporter[error]) {
	l.reporters = append(l.reporters, r)
}
