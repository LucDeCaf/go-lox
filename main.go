package main

import (
	"fmt"
	"github.com/LucDeCaf/go-lox/internal/lox"
	"github.com/LucDeCaf/go-lox/internal/lox/error_reporters"
	"os"
)

func main() {
	args := os.Args[1:]

	lox := lox.NewLox()
	lox.RegisterErrorReporter(error_reporters.NewStdoutReporter())

	switch len(args) {
	case 0:
		lox.RunPrompt()
	case 1:
		lox.RunFile(args[0])
	default:
		fmt.Println("Usage: lox [script]")
	}
}
