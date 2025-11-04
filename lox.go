package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	interpreter *Interpreter
	hadError    bool
}

func NewLox() Lox {
	return Lox{
		interpreter: NewInterpreter(),
		hadError:    false,
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

// TODO report errors here via struct that implements ErrorReporter
func (l *Lox) run(source string) {
	// Errors are already being reported, but they should probably be handled here
	scanner := NewScanner(l)
	tokens := scanner.scanTokens(source)
	parser := NewParser(l)
	ast := parser.Parse(tokens)

	if l.hadError || ast == nil {
		return
	}

	// Reuse same interpreter in every call
	v := l.interpreter.Interpret(ast)
	if v == nil {
		fmt.Println("WARNING: An error likely occured")
	}
	fmt.Printf("%v\n", v)
}

// TODO replace with actual error reporting
func (l *Lox) reportError(token Token, err error) {
	if token.tokenType == EOF {
		l.report(token.line, " at end", err)
	} else {
		l.report(token.line, " at '"+token.lexeme+"'", err)
	}
}

func (l *Lox) report(line int, where string, err error) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, err.Error())
	l.hadError = true
}
