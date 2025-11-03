package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	lox := NewLox()

	switch len(args) {
	case 0:
		lox.runPrompt()
	case 1:
		lox.runFile(args[1])
	default:
		fmt.Println("Usage: lox [script]")
	}
}
