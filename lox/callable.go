package lox

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, args []any) any
}
