package lox

type ReturnValue struct {
	Value any
}

func NewReturnValue(value any) *ReturnValue {
	return &ReturnValue{value}
}
