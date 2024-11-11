package lox

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) any
	VisitClassStmt(stmt *Class) any
	VisitExpressionStmt(stmt *Expression) any
	VisitFunctionStmt(stmt *Function) any
	VisitIfStmt(stmt *If) any
	VisitPrintStmt(stmt *Print) any
	VisitReturnStmt(stmt *Return) any
	VisitVarStmt(stmt *Var) any
	VisitWhileStmt(stmt *While) any
}

type Stmt interface {
	Accept(v StmtVisitor) any
}

type Block struct {
	Statements []Stmt
}

func NewBlock(statements []Stmt, ) Stmt {
	return &Block{ statements,  }
}

func (b *Block) Accept(sv StmtVisitor) any {
	return sv.VisitBlockStmt(b)
}

type Class struct {
	Name *Token
	Superclass *Variable
	Methods []*Function
}

func NewClass(name *Token, superclass *Variable, methods []*Function, ) Stmt {
	return &Class{ name, superclass, methods,  }
}

func (c *Class) Accept(sv StmtVisitor) any {
	return sv.VisitClassStmt(c)
}

type Expression struct {
	Expression Expr
}

func NewExpression(expression Expr, ) Stmt {
	return &Expression{ expression,  }
}

func (e *Expression) Accept(sv StmtVisitor) any {
	return sv.VisitExpressionStmt(e)
}

type Function struct {
	Name *Token
	Params []*Token
	Body []Stmt
}

func NewFunction(name *Token, params []*Token, body []Stmt, ) Stmt {
	return &Function{ name, params, body,  }
}

func (f *Function) Accept(sv StmtVisitor) any {
	return sv.VisitFunctionStmt(f)
}

type If struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt, ) Stmt {
	return &If{ condition, thenBranch, elseBranch,  }
}

func (i *If) Accept(sv StmtVisitor) any {
	return sv.VisitIfStmt(i)
}

type Print struct {
	Expression Expr
}

func NewPrint(expression Expr, ) Stmt {
	return &Print{ expression,  }
}

func (p *Print) Accept(sv StmtVisitor) any {
	return sv.VisitPrintStmt(p)
}

type Return struct {
	Keyword *Token
	Value Expr
}

func NewReturn(keyword *Token, value Expr, ) Stmt {
	return &Return{ keyword, value,  }
}

func (r *Return) Accept(sv StmtVisitor) any {
	return sv.VisitReturnStmt(r)
}

type Var struct {
	Name *Token
	Initializer Expr
}

func NewVar(name *Token, initializer Expr, ) Stmt {
	return &Var{ name, initializer,  }
}

func (v *Var) Accept(sv StmtVisitor) any {
	return sv.VisitVarStmt(v)
}

type While struct {
	Condition Expr
	Body Stmt
}

func NewWhile(condition Expr, body Stmt, ) Stmt {
	return &While{ condition, body,  }
}

func (w *While) Accept(sv StmtVisitor) any {
	return sv.VisitWhileStmt(w)
}

