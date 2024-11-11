package main

import (
	"fmt"
	"os"

	"lox/lox"
)

func main() {
	lox := lox.NewLox()
	if len(os.Args) > 2 {
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		lox.RunFile(os.Args[1])
	} else {
		lox.RunPrompt()
	}
}
