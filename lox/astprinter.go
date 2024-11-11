package lox

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) PrintExpr(expr Expr) any {
	return expr.Accept(a)
}

func (a *AstPrinter) PrintStmt(stmt Stmt) any {
	return stmt.Accept(a)
}

func (a *AstPrinter) VisitBlockStmt(stmt *Block) any {
	var builder strings.Builder
	builder.WriteString("(block ")

	for _, stmt := range stmt.Statements {
		builder.WriteString(stmt.Accept(a).(string))
	}

	builder.WriteRune(')')
	return builder.String()
}

func (a *AstPrinter) VisitClassStmt(stmt *Class) any {
	return a.parenthesizeAny("class", stmt.Name)
}

func (a *AstPrinter) VisitExpressionStmt(stmt *Expression) any {
	return a.parenthesize("expr", stmt.Expression)
}

func (a *AstPrinter) VisitFunctionStmt(stmt *Function) any {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("(%s (", stmt.Name.Lexeme))

	for _, param := range stmt.Params {
		if param != stmt.Params[0] {
			builder.WriteRune(' ')
		}
		builder.WriteString(param.Lexeme)
	}

	builder.WriteString(") ")

	for _, stmt := range stmt.Body {
		builder.WriteString(stmt.Accept(a).(string))
	}

	builder.WriteRune(')')
	return builder.String()
}

func (a *AstPrinter) VisitIfStmt(stmt *If) any {
	return a.parenthesizeAny("if", stmt.Condition, stmt.ThenBranch, "else", stmt.ElseBranch)
}

func (a *AstPrinter) VisitPrintStmt(stmt *Print) any {
	return a.parenthesize("print", stmt.Expression)
}

func (a *AstPrinter) VisitReturnStmt(stmt *Return) any {
	return a.parenthesize("return", stmt.Value)
}

func (a *AstPrinter) VisitVarStmt(stmt *Var) any {
	return a.parenthesizeAny("var", stmt.Name.Lexeme, stmt.Initializer) + "\n"
}

func (a *AstPrinter) VisitWhileStmt(stmt *While) any {
	return a.parenthesizeAny("while", stmt.Condition, stmt.Body)
}

func (a *AstPrinter) VisitAssignExpr(expr *Assign) any {
	return a.parenthesizeAny("assign", expr.Name.Lexeme, expr.Value)
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitCallExpr(expr *Call) any {
	return a.parenthesizeAny("call", expr.Callee, expr.Arguments)
}

func (a *AstPrinter) VisitGetExpr(expr *Get) any {
	return a.parenthesizeAny("get", expr.Object, expr.Name.Lexeme)
}

func (a *AstPrinter) VisitGroupingExpr(expr *Grouping) any {
	return a.parenthesize("group", expr.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.Value)
}

func (a *AstPrinter) VisitLogicalExpr(expr *Logical) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a *AstPrinter) VisitSetExpr(expr *Set) any {
	return a.parenthesizeAny("set", expr.Object, expr.Name.Lexeme, expr.Value)
}

func (a *AstPrinter) VisitSuperExpr(expr *Super) any {
	return a.parenthesize("super", expr)
}

func (a *AstPrinter) VisitThisExpr(expr *This) any {
	return a.parenthesize("this", expr)
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) any {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a *AstPrinter) VisitVariableExpr(expr *Variable) any {
	return expr.Name.Lexeme
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("(%s", name))
	for _, expr := range exprs {
		builder.WriteString(fmt.Sprintf(" %s", expr.Accept(a)))
	}
	builder.WriteRune(')')
	return builder.String()
}

func (a *AstPrinter) parenthesizeAny(name string, vals ...any) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("(%s", name))
	for _, val := range vals {
		a.stringify(&builder, val)
	}
	builder.WriteRune(')')
	return builder.String()
}

func (a *AstPrinter) stringify(builder *strings.Builder, val any) {
	builder.WriteRune(' ')
	switch v := val.(type) {
	case Expr:
		builder.WriteString(v.Accept(a).(string))
	case Stmt:
		builder.WriteString(v.Accept(a).(string))
	case *Token:
		builder.WriteString(v.Lexeme)
	case []Expr:
		for _, arg := range v {
			builder.WriteString(arg.Accept(a).(string))
		}
	default:
		builder.WriteString(fmt.Sprintf("%v", v))
	}
}
