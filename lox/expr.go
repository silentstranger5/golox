package lox

type ExprVisitor interface {
	VisitAssignExpr(expr *Assign) any
	VisitBinaryExpr(expr *Binary) any
	VisitCallExpr(expr *Call) any
	VisitGetExpr(expr *Get) any
	VisitGroupingExpr(expr *Grouping) any
	VisitLiteralExpr(expr *Literal) any
	VisitLogicalExpr(expr *Logical) any
	VisitSetExpr(expr *Set) any
	VisitSuperExpr(expr *Super) any
	VisitThisExpr(expr *This) any
	VisitUnaryExpr(expr *Unary) any
	VisitVariableExpr(expr *Variable) any
}

type Expr interface {
	Accept(v ExprVisitor) any
}

type Assign struct {
	Name *Token
	Value Expr
}

func NewAssign(name *Token, value Expr, ) Expr {
	return &Assign{ name, value,  }
}

func (a *Assign) Accept(ev ExprVisitor) any {
	return ev.VisitAssignExpr(a)
}

type Binary struct {
	Left Expr
	Operator *Token
	Right Expr
}

func NewBinary(left Expr, operator *Token, right Expr, ) Expr {
	return &Binary{ left, operator, right,  }
}

func (b *Binary) Accept(ev ExprVisitor) any {
	return ev.VisitBinaryExpr(b)
}

type Call struct {
	Callee Expr
	Paren *Token
	Arguments []Expr
}

func NewCall(callee Expr, paren *Token, arguments []Expr, ) Expr {
	return &Call{ callee, paren, arguments,  }
}

func (c *Call) Accept(ev ExprVisitor) any {
	return ev.VisitCallExpr(c)
}

type Get struct {
	Object Expr
	Name *Token
}

func NewGet(object Expr, name *Token, ) Expr {
	return &Get{ object, name,  }
}

func (g *Get) Accept(ev ExprVisitor) any {
	return ev.VisitGetExpr(g)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr, ) Expr {
	return &Grouping{ expression,  }
}

func (g *Grouping) Accept(ev ExprVisitor) any {
	return ev.VisitGroupingExpr(g)
}

type Literal struct {
	Value any
}

func NewLiteral(value any, ) Expr {
	return &Literal{ value,  }
}

func (l *Literal) Accept(ev ExprVisitor) any {
	return ev.VisitLiteralExpr(l)
}

type Logical struct {
	Left Expr
	Operator *Token
	Right Expr
}

func NewLogical(left Expr, operator *Token, right Expr, ) Expr {
	return &Logical{ left, operator, right,  }
}

func (l *Logical) Accept(ev ExprVisitor) any {
	return ev.VisitLogicalExpr(l)
}

type Set struct {
	Object Expr
	Name *Token
	Value Expr
}

func NewSet(object Expr, name *Token, value Expr, ) Expr {
	return &Set{ object, name, value,  }
}

func (s *Set) Accept(ev ExprVisitor) any {
	return ev.VisitSetExpr(s)
}

type Super struct {
	Keyword *Token
	Method *Token
}

func NewSuper(keyword *Token, method *Token, ) Expr {
	return &Super{ keyword, method,  }
}

func (s *Super) Accept(ev ExprVisitor) any {
	return ev.VisitSuperExpr(s)
}

type This struct {
	Keyword *Token
}

func NewThis(keyword *Token, ) Expr {
	return &This{ keyword,  }
}

func (t *This) Accept(ev ExprVisitor) any {
	return ev.VisitThisExpr(t)
}

type Unary struct {
	Operator *Token
	Right Expr
}

func NewUnary(operator *Token, right Expr, ) Expr {
	return &Unary{ operator, right,  }
}

func (u *Unary) Accept(ev ExprVisitor) any {
	return ev.VisitUnaryExpr(u)
}

type Variable struct {
	Name *Token
}

func NewVariable(name *Token, ) Expr {
	return &Variable{ name,  }
}

func (v *Variable) Accept(ev ExprVisitor) any {
	return ev.VisitVariableExpr(v)
}

