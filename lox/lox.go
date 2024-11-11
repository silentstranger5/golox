package lox

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var LoxInstance *Lox

type Lox struct {
	hadError    bool
	interpreter *Interpreter
}

func NewLox() *Lox {
	return &Lox{false, NewInterpreter()}
}

func (l *Lox) RunFile(path string) {
	LoxInstance = l
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read file: %s\n", err)
		os.Exit(64)
	}
	l.Run(string(bytes))
	// Indicate an error in the exit code.
	if l.hadError {
		os.Exit(65)
	}
}

func (l *Lox) RunPrompt() {
	LoxInstance = l
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Failed to read line: %s\n", err)
			os.Exit(64)
		}
		l.Run(line)
		l.hadError = false
	}
}

func (l *Lox) Run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	if l.hadError {
		return
	}

	parser := NewParser(tokens)
	statements := parser.Parse()

	if l.hadError {
		return
	}

	resolver := NewResolver(l.interpreter)
	resolver.ResolveStatements(statements)

	if l.hadError {
		return
	}

	l.interpreter.Interpret(statements)
}

func (l *Lox) Report(error error) {
	fmt.Fprintln(os.Stderr, error.Error())
	l.hadError = true
}
