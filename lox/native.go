package lox

type NativeFunc struct {
	arity_ func() int
	call_  func(*Interpreter, []any) any
}

func NewNativeFunc(arity func() int, call func(*Interpreter, []any) any) *NativeFunc {
	return &NativeFunc{arity, call}
}

func (bf *NativeFunc) arity() int {
	return bf.arity_()
}

func (bf *NativeFunc) call(i *Interpreter, args []any) any {
	return bf.call_(i, args)
}

func (bf *NativeFunc) String() string {
	return "<native fn>"
}
