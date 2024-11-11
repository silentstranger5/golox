package lox

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction int
	currentClass    int
}

const (
	FN_NONE = iota
	FN_FUNCTION
	FN_INITIALIZER
	FN_METHOD
)

const (
	CLS_NONE = iota
	CLS_CLASS
	CLS_SUBCLASS
)

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{interpreter, make([]map[string]bool, 0), FN_NONE, CLS_NONE}
}

func (r *Resolver) ResolveStatements(statements []Stmt) {
	defer func() {
		if r := recover(); r != nil {
			LoxInstance.Report(r.(error))
		}
	}()

	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) VisitBlockStmt(stmt *Block) any {
	r.beginScope()
	r.ResolveStatements(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *Class) any {
	enclosingClass := r.currentClass
	r.currentClass = CLS_CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil &&
		stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		panic(NewResolveError(
			stmt.Superclass.Name,
			"A class can't inherit from itself.",
		))
	}

	if stmt.Superclass != nil {
		r.currentClass = CLS_SUBCLASS
		r.resolveExpr(stmt.Superclass)
	}

	if stmt.Superclass != nil {
		r.beginScope()
		r.scopes[len(r.scopes)-1]["super"] = true
	}

	r.beginScope()
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, method := range stmt.Methods {
		declaration := FN_METHOD
		if method.Name.Lexeme == "init" {
			declaration = FN_INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}

	r.endScope()

	if stmt.Superclass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *Expression) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *Function) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt, FN_FUNCTION)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *If) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStatement(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStatement(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *Print) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *Return) any {
	if r.currentFunction == FN_NONE {
		panic(NewResolveError(
			stmt.Keyword, "Can't return from top-level code.",
		))
	}
	if stmt.Value != nil {
		if r.currentFunction == FN_INITIALIZER {
			panic(NewResolveError(
				stmt.Keyword, "Can't return a value from an initializer.",
			))
		}
		r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitVarStmt(stmt *Var) any {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *While) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStatement(stmt.Body)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *Assign) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *Binary) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *Call) any {
	r.resolveExpr(expr.Callee)

	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}

	return nil
}

func (r *Resolver) VisitGetExpr(expr *Get) any {
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *Grouping) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *Literal) any {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *Logical) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitSetExpr(expr *Set) any {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitSuperExpr(expr *Super) any {
	if r.currentClass == CLS_NONE {
		panic(NewResolveError(
			expr.Keyword, "Can't use 'super' outside of a class.",
		))
	} else if r.currentClass != CLS_SUBCLASS {
		panic(NewResolveError(
			expr.Keyword, "Can't use 'super' in a class with no superclass.",
		))
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitThisExpr(expr *This) any {
	if r.currentClass == CLS_NONE {
		panic(NewResolveError(
			expr.Keyword, "Can't use 'this' outside a class."))
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *Variable) any {
	if len(r.scopes) > 0 {
		val, ok := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]
		if ok && !val {
			panic(NewResolveError(
				expr.Name, "Can't read local variable in it's own initializer.",
			))
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) resolveStatement(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveFunction(function *Function, type_ int) {
	enclosingFunction := r.currentFunction
	r.currentFunction = type_
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStatements(function.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) resolveLocal(expr Expr, name *Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name *Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.Lexeme]; ok {
		panic(NewResolveError(
			name, "Already a variable with this name in this scope.",
		))
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = false
}

func (r *Resolver) define(name *Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}
