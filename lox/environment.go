package lox

import "fmt"

type Environment struct {
	Values    map[string]any
	Enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{make(map[string]any), enclosing}
}

func (e *Environment) Get(name *Token) any {
	val, ok := e.Values[name.Lexeme]
	if ok {
		return val
	} else if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	} else {
		panic(NewRuntimeError(
			name, fmt.Sprintf("Undefined variable: '%s'.", name.Lexeme),
		))
	}
}

func (e *Environment) Assign(name *Token, value any) {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return
	} else if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
	} else {
		panic(NewRuntimeError(
			name, fmt.Sprintf("Undefined variable: '%s'.", name.Lexeme),
		))
	}
}

func (e *Environment) Define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) Ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.Enclosing
	}
	return environment
}

func (e *Environment) GetAt(distance int, name string) any {
	return e.Ancestor(distance).Values[name]
}

func (e *Environment) AssignAt(distance int, name *Token, value any) {
	e.Ancestor(distance).Values[name.Lexeme] = value
}
