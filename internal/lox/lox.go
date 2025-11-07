package lox

import (
	"bufio"
	"fmt"
	"github.com/LucDeCaf/go-lox/internal/lox/error_reporters"
	"os"
)

type Lox struct {
	interpreter *Interpreter
	reporters   []error_reporters.ErrorReporter[error]
	hadError    bool
}

func NewLox() Lox {
	return Lox{
		interpreter: NewInterpreter(),
		reporters:   []error_reporters.ErrorReporter[error]{},
		hadError:    false,
	}
}

func (l *Lox) RunPrompt() {
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

func (l *Lox) RunFile(path string) (err error) {
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
	statements, parseOk := parser.parse(tokens)
	for _, err := range parser.errors {
		for _, r := range l.reporters {
			r.ReportError(err)
		}
	}

	if !scanOk || !parseOk {
		return
	}

	l.interpreter.Interpret(statements)
}

func (l *Lox) RegisterErrorReporter(r error_reporters.ErrorReporter[error]) {
	l.reporters = append(l.reporters, r)
}
