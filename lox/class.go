package lox

type LoxClass struct {
	Name       string
	Methods    map[string]*LoxFunction
	Superclass *LoxClass
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{name, methods, superclass}
}

func (c *LoxClass) String() string {
	return c.Name
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []any) any {
	instance := NewInstance(c)

	initializer := c.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (c *LoxClass) Arity() int {
	initializer := c.FindMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if val, ok := c.Methods[name]; ok {
		return val
	}

	if c.Superclass != nil {
		return c.Superclass.FindMethod(name)
	}

	return nil
}
