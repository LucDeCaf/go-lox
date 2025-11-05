package main

import (
	"fmt"
	"github.com/LucDeCaf/go-lox/error_reporters"
	"os"
)

func main() {
	args := os.Args[1:]

	lox := NewLox()
	lox.registerErrorReporter(error_reporters.NewStdoutReporter())

	switch len(args) {
	case 0:
		lox.runPrompt()
	case 1:
		lox.runFile(args[0])
	default:
		fmt.Println("Usage: lox [script]")
	}
}
