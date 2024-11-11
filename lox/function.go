package lox

import "fmt"

type LoxFunction struct {
	Declaration   *Function
	Closure       *Environment
	IsInitializer bool
}

func NewLoxFunction(declaration *Function, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{declaration, closure, isInitializer}
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) Bind(instance *Instance) *LoxFunction {
	environment := NewEnvironment(f.Closure)
	environment.Define("this", instance)
	return NewLoxFunction(f.Declaration, environment, f.IsInitializer)
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []any) (ret any) {
	enclosing := interpreter.environment

	defer func() {
		if r := recover(); r != nil {
			if val, ok := r.(*ReturnValue); ok {
				if f.IsInitializer {
					ret = f.Closure.GetAt(0, "this")
				} else {
					ret = val.Value
				}
				interpreter.environment = enclosing
			} else {
				panic(r)
			}
		}
	}()

	environment := NewEnvironment(f.Closure)
	for i := 0; i < len(f.Declaration.Params); i++ {
		environment.Define(f.Declaration.Params[i].Lexeme,
			arguments[i])
	}

	interpreter.executeBlock(f.Declaration.Body, environment)

	if f.IsInitializer {
		return f.Closure.GetAt(0, "this")
	}

	return nil
}

func (f *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}
