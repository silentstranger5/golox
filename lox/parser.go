package lox

import "fmt"

type Parser struct {
	tokens  []*Token
	current int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() []Stmt {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			LoxInstance.Report(r.(error))
			p.synchronize()
		}
	}()

	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect class name.")

	var superclass *Variable
	if p.match(LESS) {
		p.consume(IDENTIFIER, "Expect superclass name.")
		superclass = NewVariable(p.previous()).(*Variable)
	}

	p.consume(LEFT_BRACE, "Expect '{' before class body.")

	methods := make([]*Function, 0)
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		function, _ := p.function("method").(*Function)
		methods = append(methods, function)
	}

	p.consume(RIGHT_BRACE, "Expect '}' after class body.")

	return NewClass(name, superclass, methods)
}

func (p *Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	} else if p.match(IF) {
		return p.ifStatement()
	} else if p.match(PRINT) {
		return p.printStatement()
	} else if p.match(RETURN) {
		return p.returnStatement()
	} else if p.match(WHILE) {
		return p.whileStatement()
	} else if p.match(LEFT_BRACE) {
		return NewBlock(p.block())
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Stmt
	if p.match(SEMICOLON) {
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clause.")

	body := p.statement()

	if increment != nil {
		body = NewBlock([]Stmt{
			body, NewExpression(increment),
		})
	}

	if condition == nil {
		condition = NewLiteral(true)
	}

	body = NewWhile(condition, body)

	if initializer != nil {
		body = NewBlock([]Stmt{
			initializer, body,
		})
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	var thenBranch Stmt = p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return NewIf(condition, thenBranch, elseBranch)
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return NewPrint(value)
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		value = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after return value.")
	return NewReturn(keyword, value)
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return NewVar(name, initializer)
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()

	return NewWhile(condition, body)
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return NewExpression(expr)
}

func (p *Parser) function(kind string) Stmt {
	name := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))
	parameters := make([]*Token, 0)
	if !p.check(RIGHT_PAREN) {
		parameters = append(parameters,
			p.consume(IDENTIFIER, "Expect parameter name."),
		)
		for p.match(COMMA) {
			if len(parameters) >= 255 {
				panic(NewRuntimeError(
					p.peek(), "Can't have more than 255 parameters.",
				))
			}
			parameters = append(parameters,
				p.consume(IDENTIFIER, "Expect parameter name."),
			)
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	body := p.block()
	return NewFunction(name, parameters, body)
}

func (p *Parser) block() []Stmt {
	statements := make([]Stmt, 0)

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if _, ok := expr.(*Variable); ok {
			name := expr.(*Variable).Name
			return NewAssign(name, value)
		} else if val, ok := expr.(*Get); ok {
			get := val
			return NewSet(get.Object, get.Name, value)
		}
		panic(NewParseError(equals, "Invalid assignment target."))
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = NewLogical(expr, operator, right)
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = NewLogical(expr, operator, right)
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return NewUnary(operator, right)
	}
	return p.call()
}

func (p *Parser) finishCall(callee Expr) Expr {
	arguments := make([]Expr, 0)
	if !p.check(RIGHT_PAREN) {
		arguments = append(arguments, p.expression())
		for p.match(COMMA) {
			if len(arguments) >= 255 {
				panic(NewParseError(p.peek(), "Can't have more than 255 arguments."))
			}
			arguments = append(arguments, p.expression())
		}
	}

	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")

	return NewCall(callee, paren, arguments)
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name := p.consume(IDENTIFIER,
				"Expect property name after '.'")
			expr = NewGet(expr, name)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return NewLiteral(false)
	}

	if p.match(TRUE) {
		return NewLiteral(true)
	}

	if p.match(NIL) {
		return NewLiteral(nil)
	}

	if p.match(NUMBER, STRING) {
		return NewLiteral(p.previous().Literal)
	}

	if p.match(SUPER) {
		keyword := p.previous()
		p.consume(DOT, "Expect '.' after 'super'.")
		method := p.consume(
			IDENTIFIER, "Expect superclass method name.",
		)
		return NewSuper(keyword, method)
	}

	if p.match(THIS) {
		return NewThis(p.previous())
	}

	if p.match(IDENTIFIER) {
		return NewVariable(p.previous())
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(expr)
	}

	panic(NewParseError(p.peek(), "Expect expression."))
}

func (p *Parser) match(types ...TokenType) bool {
	for _, ttype := range types {
		if p.check(ttype) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(ttype TokenType, message string) *Token {
	if p.check(ttype) {
		return p.advance()
	}
	panic(NewParseError(p.peek(), message))
}

func (p *Parser) check(ttype TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == ttype
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}
