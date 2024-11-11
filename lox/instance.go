package lox

import "fmt"

type Instance struct {
	class  *LoxClass
	fields map[string]any
}

func NewInstance(class *LoxClass) *Instance {
	return &Instance{class, make(map[string]any)}
}

func (i *Instance) Get(name *Token) any {
	if val, ok := i.fields[name.Lexeme]; ok {
		return val
	}

	method := i.class.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(i)
	}

	panic(NewRuntimeError(name,
		fmt.Sprintf("Undefined property '%s'.", name.Lexeme)))
}

func (i *Instance) Set(name *Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i *Instance) String() string {
	return i.class.Name + " instance"
}
