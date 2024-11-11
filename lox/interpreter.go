package lox

import (
	"fmt"
	"strings"
	"time"
)

type Interpreter struct {
	environment *Environment
	globals     *Environment
	locals      map[Expr]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)
	arity := func() int {
		return 0
	}
	call := func(i *Interpreter, args []any) any {
		return time.Now().UnixMilli()
	}
	globals.Define("clock", NewNativeFunc(arity, call))
	return &Interpreter{globals, globals, make(map[Expr]int)}
}

func (i *Interpreter) Interpret(statements []Stmt) {
	for _, stmt := range statements {
		i.execute(stmt)
	}
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) Resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) {
	previous := i.environment
	i.environment = environment

	for _, statement := range statements {
		i.execute(statement)
	}

	i.environment = previous
}

func (i *Interpreter) VisitBlockStmt(stmt *Block) any {
	i.executeBlock(stmt.Statements, NewEnvironment(i.environment))
	return nil
}

func (i *Interpreter) VisitClassStmt(stmt *Class) any {
	var superclass *LoxClass

	if stmt.Superclass != nil {
		var ok bool
		superclass, ok = i.evaluate(stmt.Superclass).(*LoxClass)
		if !ok {
			panic(NewRuntimeError(
				stmt.Superclass.Name,
				"Superclass must be a class.",
			))
		}
	}

	i.environment.Define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		i.environment = NewEnvironment(i.environment)
		i.environment.Define("super", superclass)
	}

	methods := make(map[string]*LoxFunction)
	for _, method := range stmt.Methods {
		function := NewLoxFunction(method, i.environment,
			method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}

	class := NewLoxClass(stmt.Name.Lexeme, superclass, methods)

	if superclass != nil {
		i.environment = i.environment.Enclosing
	}

	i.environment.Assign(stmt.Name, class)
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) any {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *Function) any {
	function := NewLoxFunction(stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *If) any {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) any {
	value := i.evaluate(stmt.Expression)
	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *Return) any {
	var value any
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	panic(NewReturnValue(value))
}

func (i *Interpreter) VisitVarStmt(stmt *Var) any {
	var value any
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *While) any {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i *Interpreter) VisitAssignExpr(expr *Assign) any {
	value := i.evaluate(expr.Value)

	distance, ok := i.locals[expr]
	if ok {
		i.environment.AssignAt(distance, expr.Name, value)
	} else {
		i.globals.Assign(expr.Name, value)
	}

	return value
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case GREATER:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case MINUS:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	case PLUS:
		lval, lok := left.(float64)
		rval, rok := right.(float64)
		if lok && rok {
			return lval + rval
		}

		lstr, lok := left.(string)
		rstr, rok := right.(string)
		if lok && rok {
			return lstr + rstr
		}

		panic(NewRuntimeError(
			expr.Operator,
			"Operands must be two numbers or strings.",
		))
	case SLASH:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	}
	return nil
}

func (i *Interpreter) VisitCallExpr(expr *Call) any {
	callee := i.evaluate(expr.Callee)

	arguments := make([]any, 0)
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	function, ok := callee.(Callable)
	if !ok {
		panic(NewRuntimeError(
			expr.Paren, "Can only call functions and classes.",
		))
	}

	if len(arguments) != function.Arity() {
		panic(NewRuntimeError(
			expr.Paren,
			fmt.Sprintf("Expected %d arguments but got %d.",
				function.Arity(), len(arguments)),
		))
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitGetExpr(expr *Get) any {
	object := i.evaluate(expr.Object)
	if val, ok := object.(*Instance); ok {
		return val.Get(expr.Name)
	}
	panic(NewRuntimeError(expr.Name,
		"Only instances have properties."))
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) any {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr *Logical) any {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitSetExpr(expr *Set) any {
	object := i.evaluate(expr.Object)

	if _, ok := object.(*Instance); !ok {
		panic(NewRuntimeError(expr.Name,
			"Only instances have fields."))
	}

	value := i.evaluate(expr.Value)
	object.(*Instance).Set(expr.Name, value)
	return value
}

func (i *Interpreter) VisitSuperExpr(expr *Super) any {
	distance := i.locals[expr]
	superclass := i.environment.GetAt(distance, "super").(*LoxClass)
	object := i.environment.GetAt(distance-1, "this").(*Instance)
	method := superclass.FindMethod(expr.Method.Lexeme)

	if method == nil {
		panic(NewResolveError(
			expr.Method,
			fmt.Sprint("Undefined property '%s'.", expr.Method.Lexeme),
		))
	}

	return method.Bind(object)
}

func (i *Interpreter) VisitThisExpr(expr *This) any {
	return i.lookupVariable(expr.Keyword, expr)
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case MINUS:
		i.checkNumberOperand(expr.Operator, right)
		return right.(float64)
	case BANG:
		return !i.isTruthy(right)
	}

	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *Variable) any {
	return i.lookupVariable(expr.Name, expr)
}

func (i *Interpreter) lookupVariable(name *Token, expr Expr) any {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.GetAt(distance, name.Lexeme)
	} else {
		return i.globals.Get(name)
	}
}

func (i *Interpreter) checkNumberOperand(operator *Token, operand any) {
	if _, ok := operand.(float64); ok {
		return
	}
	panic(NewRuntimeError(operator, "Operand must be a number."))
}

func (i *Interpreter) checkNumberOperands(operator *Token, left, right any) {
	_, lok := left.(float64)
	_, rok := right.(float64)
	if lok && rok {
		return
	}

	panic(NewRuntimeError(operator, "Operands must be numbers"))
}

func (i *Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}
	if val, ok := object.(bool); ok {
		return val
	}
	return true
}

func (i *Interpreter) isEqual(a, b any) bool {
	return a == b
}

func (i *Interpreter) stringify(object any) string {
	if object == nil {
		return "nil"
	}

	if val, ok := object.(float64); ok {
		text := fmt.Sprintf("%v", val)
		if strings.HasSuffix(text, ".0") {
			text = text[0 : len(text)-2]
		}
		return text
	}

	return fmt.Sprintf("%v", object)
}
