package parser

import "fmt"

type Parser struct {
	tokens  []Token
	current int
	prev    *Token
}

func NewParser(tokens []Token) Parser {
	return Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() []Statement {
	var statements []Statement
	at_end := p.IsAtEnd()
	for !at_end {
		stmt, err := p.declaration()
		if err != nil {
			// errorhandling.RuntimeError(err)
			fmt.Println("ERRROR")
			return nil
		}
		statements = append(statements, stmt)
		at_end = p.IsAtEnd()
	}

	return statements
}

func (p *Parser) declaration() (Statement, error) {
	var stmt Statement
	var err error
	if p.match(VAR) {
		stmt, err = p.varDeclaration()
	} else if p.match(FUN) {
		stmt, err = p.funcDeclaration()
	} else if p.match(CLASS) {
		stmt, err = p.classDeclaration()
	} else {
		stmt, err = p.statement()
	}
	if err != nil {
		p.syncronize()
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) statement() (Statement, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.peek().Token_type == LEFT_BRACE {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return Block{Statements: stmts}, nil
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(FOR) {
		return p.forStatement()
	}

	if p.peek().Token_type == RETURN {
		return p.returnStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) returnStatement() (Statement, error) {
	_, err := p.consume(RETURN, "expected 'return'")
	if err != nil {
		return nil, err
	}
	if p.match(SEMICOLON) {
		return Return{Return_expr: nil}, nil
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "expected ';'")
	if err != nil {
		return nil, err
	}

	return Return{Return_expr: expr}, nil
}

func (p *Parser) forStatement() (Statement, error) {
	var initializer_stmt Statement
	var conditional_expr Expr
	var increment_expression Expr
	var loop_stmt Statement

	_, err := p.consume(LEFT_PAREN, "expected '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	if p.match(VAR) {
		initializer_stmt, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else if p.match(SEMICOLON) {
		initializer_stmt = nil
	} else {
		initializer_stmt, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}
	if p.match(SEMICOLON) {
		conditional_expr = nil
	} else {
		conditional_expr, err = p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(SEMICOLON, "expected ';' after conditional ")
		if err != nil {
			return nil, err
		}
	}

	if p.match(RIGHT_PAREN) {
		increment_expression = nil
	} else {
		increment_expression, err = p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "expected ')' after expression")
		if err != nil {
			return nil, err
		}
	}

	loop_stmt, err = p.statement()
	if err != nil {
		return nil, err
	}
	if increment_expression != nil {
		body := []Statement{loop_stmt, ExpressionStmt{increment_expression}}
		loop_stmt = Block{Statements: body}
	}
	if conditional_expr == nil {
		conditional_expr = Literal{Value: true}
	}
	var body Statement = While{Conditional: conditional_expr, Stmt: loop_stmt}
	if initializer_stmt != nil {
		tmp := []Statement{initializer_stmt, body}
		body = Block{Statements: tmp}
	}

	return body, nil // for stmt

}

func (p *Parser) whileStatement() (Statement, error) {
	var err error
	_, err = p.consume(LEFT_PAREN, "expected '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	conditional_stmt, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "expected ')' after conditional ")
	if err != nil {
		return nil, err
	}

	while_body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return While{Conditional: conditional_stmt, Stmt: while_body}, nil
}

func (p *Parser) ifStatement() (Statement, error) {
	var else_stmt Statement

	_, err := p.consume(LEFT_PAREN, "expected '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "expected ')' after ")
	if err != nil {
		return nil, err
	}

	if_stmt, err := p.statement()
	if err != nil {
		return nil, err
	}

	if p.match(ELSE) {
		else_stmt, err = p.statement()
	}
	if err != nil {
		return nil, err
	}

	// return NewIfStatement(expr, if_stmt, else_stmt), nil
	return If{Conditional: expr, If_stmt: if_stmt, Else_stmt: else_stmt}, nil
}

func (p *Parser) printStatement() (Statement, error) {
	var stmt Statement
	expr, err := p.expression()
	if err != nil {
		return stmt, err
	}
	_, err = p.consume(SEMICOLON, "expected ';' after value.")
	if err != nil {
		// The semi colon after "false" is being consumed. I expect it's a scanning bug.
		return stmt, err
	}

	return Print{Val: expr}, nil
}

func (p *Parser) expressionStatement() (Statement, error) {
	var stmt Statement
	expr, err := p.expression()
	if err != nil {
		return stmt, err
	}
	_, err = p.consume(SEMICOLON, "expected ';' at end of line")
	if err != nil {
		return stmt, err
	}

	return ExpressionStmt{Val: expr}, nil
}

func (p *Parser) block() ([]Statement, error) {
	var statements []Statement

	_, err := p.consume(LEFT_BRACE, "expected '{'")
	if err != nil {
		return nil, err
	}

	for !p.check(RIGHT_BRACE) && !p.IsAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	_, err = p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) classDeclaration() (Statement, error) {
	var classId Token
	var parentClass *Variable
	var err error
	var funcs []Function

	classId, err = p.consume(IDENTIFIER, "expected an identifier")
	if err != nil {
		return nil, err
	}

	if p.match(LESS) {
		_, err = p.consume(IDENTIFIER, "expected an identifier")
		if err != nil {
			return nil, err
		}
		parentClass = &Variable{Name: p.previous()}
	}

	_, err = p.consume(LEFT_BRACE, "expected '{'")
	if err != nil {
		return nil, err
	}

	if p.match(RIGHT_BRACE) {
		return Class{Name: classId}, nil // return class here
	}

	for !p.match(RIGHT_BRACE) && !p.IsAtEnd() {
		val, err := p.funcDeclaration()
		if err != nil {
			return nil, err
		}
		fun, ok := val.(Function)
		if !ok {
			return nil, NewParseError("expected a function definition")
		}
		funcs = append(funcs, fun)
	}

	return Class{
		Name:        classId,
		Methods:     funcs,
		ParentClass: parentClass,
	}, nil // return class here
}

func (p *Parser) funcDeclaration() (Statement, error) {
	// function
	return p.function()
}

func (p *Parser) function() (Statement, error) {
	var funcId Token
	var err ParseError
	var identifers []Token

	if !p.match(IDENTIFIER) {
		err = NewParseError("expected an identifer")
		return nil, &err
	}
	funcId = p.previous()

	_, pErr := p.consume(LEFT_PAREN, "expected '('.")
	if pErr != nil {
		return nil, pErr
	}

	if p.peek().Token_type != RIGHT_PAREN {
		identifers, pErr = p.identifiers()
		if pErr != nil {
			return nil, pErr
		}
	}

	_, pErr = p.consume(RIGHT_PAREN, "expected ')'.")
	if pErr != nil {
		return nil, pErr
	}

	block, pErr := p.block()
	if pErr != nil {
		return nil, pErr
	}

	return Function{Name: funcId, Params: identifers, Body: block}, nil
}

func (p *Parser) identifiers() ([]Token, error) {
	// parameters     → IDENTIFIER ( "," IDENTIFIER )* ;
	var tokens []Token
	var err ParseError

	if !p.match(IDENTIFIER) {
		err = NewParseError("expected an idenifier.")
		return nil, &err
	}
	tokens = append(tokens, p.previous())

	if p.match(COMMA) {
		val, pErr := p.identifiers()
		if pErr != nil {
			return nil, pErr
		}
		return append(tokens, val...), nil
	}

	return tokens, nil
}

func (p *Parser) varDeclaration() (Statement, error) {
	var initializer Expr
	var err error

	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	if p.match(EQUAL) {
		initializer, err = p.expression()
	} else {
		initializer = Literal{Value: nil}
	}

	if err != nil {
		return nil, err
	}
	_, err = p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return Var{Name: name, Initializer: initializer}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	left, err := p.logic_or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		right, err := p.assignment()
		if err != nil {
			return nil, err
		}

		switch t := left.(type) {
		case Variable:
			return Assign{Name: t.Name, Value: right}, nil
		case Get:
			return Set{Object: t.Object, Name: t.Name, Value: right}, nil // Set expression
		default:
			err := NewParseError("Left side of assignment must be a variable.")
			return nil, &err
		}

		val, ok := left.(Variable)
		if !ok {
			err := NewParseError("Left side of assignment must be a variable.")
			return nil, &err
		}
		return Assign{Name: val.Name, Value: right}, nil
	}

	return left, nil
}

func (p *Parser) logic_or() (Expr, error) {
	left, err := p.logic_and()
	if err != nil {
		return nil, err
	}

	if p.match(OR) {
		op := p.previous()
		right, err := p.logic_or()
		if err != nil {
			return nil, err
		}

		return Logical{Left: left, Operator: op, Right: right}, nil
	}

	return left, nil
}

func (p *Parser) logic_and() (Expr, error) {
	left, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(AND) {
		op := p.previous()
		right, err := p.logic_and()
		if err != nil {
			return nil, err
		}

		return Logical{Left: left, Operator: op, Right: right}, nil
	}

	return left, nil
}

func (p *Parser) equality() (Expr, error) {
	// First binary expression
	prefix, err := p.comparison()
	if err != nil {
		return Unary{}, err
	}
	// Recursive case:
	//  Current token is an equality Operator
	//  Consume the equality operator, increment the token counter, and return a binary expression with Left: prefix, Operator: op, Right: p.equality()
	// p.current += 1
	if m := p.match(EQUAL_EQUAL, BANG_EQUAL); m {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return Unary{}, nil
		}
		return Binary{Left: prefix, Operator: op, Right: right}, nil
	}

	// Base case:
	//  current token is not an equality Operator
	// Return prefix
	return prefix, nil
}

func (p *Parser) comparison() (Expr, error) {
	prefix, err := p.term()
	if err != nil {
		return Unary{}, err
	}

	if p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return Unary{}, nil
		}
		return Binary{Left: prefix, Operator: op, Right: right}, nil
	}

	return prefix, nil
}

// term           → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (Expr, error) {
	prefix, err := p.factor()
	if err != nil {
		return Unary{}, err
	}

	if p.match(MINUS, PLUS) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return Unary{}, err
		}
		return Binary{Left: prefix, Operator: op, Right: right}, nil
	}

	return prefix, nil
}

// factor         → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (Expr, error) {
	prefix, err := p.unary()
	if err != nil {
		return Unary{}, err
	}

	if p.match(STAR, SLASH) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return Unary{}, err
		}
		return Binary{Left: prefix, Operator: op, Right: right}, nil
	}

	return prefix, nil
}

// unary          → ( "!" | "-" ) unary
//
//	| primary ;
func (p *Parser) unary() (Expr, error) {
	var right Expr
	var err error

	// prefix := p.advance()
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err = p.unary()
		if err != nil {
			return Unary{}, err
		}
		return Unary{Operator: operator, Right: right}, nil
	}

	primary, err := p.call()

	if err != nil {
		return Unary{}, err
	}

	return primary, nil
}

func (p *Parser) call() (Expr, error) {
	var expr Expr
	var err error

	expr, err = p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LEFT_PAREN) {
			expr, err = p.add_args(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(DOT) {
			name, err := p.consume(IDENTIFIER, "expected an identifier afer \".\"")
			if err != nil {
				return nil, err
			}
			expr = Get{Object: expr, Name: name}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) add_args(expr Expr) (Expr, error) {
	paren := p.previous()
	if p.match(RIGHT_PAREN) {
		return Call{Callee: expr, Paren: paren, Args: nil}, nil
	}

	args, err := p.arguments()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(RIGHT_PAREN, "expected ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return Call{Callee: expr, Paren: paren, Args: args}, nil // call expr with args, nil
}

func (p *Parser) arguments() ([]Expr, error) {
	var args []Expr
	var err error
	for cur_arg, err := p.expression(); err == nil; cur_arg, err = p.expression() {
		if len(args) >= 255 {
			p.error(p.peek(), "Can't have more than 255 argumens.")
		}
		args = append(args, cur_arg)
		if !p.match(COMMA) {
			return args, nil
		}
	}
	return nil, err
}

// primary        → NUMBER | STRING | "true" | "false" | "nil" | IDENTIFIER | (expression)
//
//	| "(" expression ")"
func (p *Parser) primary() (Expr, error) {
	var err error
	var expr Expr
	if p.match(FALSE) {
		return Literal{Value: false}, nil
	}
	if p.match(TRUE) {
		return Literal{Value: true}, nil
	}
	if p.match(NIL) {
		return Literal{Value: nil}, nil
	}
	if p.match(STRING, NUMBER) {
		return Literal{Value: p.previous().Literal}, nil
	}
	if p.match(LEFT_PAREN) {
		expr, err = p.expression()
		_, err = p.consume(RIGHT_PAREN, "expected right paren!")
		if err != nil {
			return Unary{}, err
		}

		return Grouping{Expr: expr}, nil
	}
	if p.match(IDENTIFIER) {
		return Variable{p.previous()}, nil
	}
	if p.match(THIS) {
		return This{Keyword: p.previous()}, nil
	}
	if p.match(SUPER) {
		keyword := p.previous()
		_, err := p.consume(DOT, "expected \".\" after \"super\"")
		if err != nil {
			return nil, err
		}
		id, err := p.consume(IDENTIFIER, "expected an identifier after 'super.'")
		if err != nil {
			return nil, err
		}

		return Super{Keyword: keyword, Method: id}, nil
	}
	if p.match(ERROR) {
		e := p.error(p.peek(), p.prev.Lexeme)
		return nil, e
	}
	parse_error := p.error(p.peek(), "invalid token")

	return nil, &parse_error
}

// func (p *Parser) identifier() (Expr, error)

func (p *Parser) syncronize() {
	p.advance()

	for !p.IsAtEnd() {
		if p.previous().Token_type == SEMICOLON {
			return
		}
		t := p.peek().Token_type
		switch t {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(token_type ...TokenType) bool {
	for _, v := range token_type {
		if p.check(v) {
			p.advance()
			return true
		}
	}

	return false
}

func (p Parser) check(token_type TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.peek().Token_type == token_type
}

func (p Parser) cur() *Token {
	return &p.tokens[p.current]
}

func (p *Parser) advance() Token {
	for {
		if !p.IsAtEnd() {
			p.prev = p.cur()
			p.current += 1
			if p.peek().Token_type == ERROR {
				fmt.Println(p.peek().Lexeme)
				continue
			}
		}
		break
	}

	ret := p.previous()

	return ret
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	error := p.error(p.peek(), message)

	return Token{}, &error
}

func (p Parser) error(token Token, message string) ParseError {
	if token.Token_type == EOF {
		// errorhandling.Report(token.Line, " at end ", message)
		fmt.Printf("[line %d] Error %s: '%s'\n", token.Line, "at end", message)
	} else {
		// errorhandling.Report(token.Line, " at "+token.Lexeme, message)
		fmt.Printf("[line %d] Error at %s: %s\n", token.Line, token.Lexeme, message)
	}
	return NewParseError(message)
}

func (p Parser) peek() Token {
	return p.tokens[p.current]
}

func (p Parser) previous() Token {
	return *p.prev
}

func (p Parser) IsAtEnd() bool {
	return p.peek().Token_type == EOF
}

func (p Parser) errorAtCurrent() {
	fmt.Println(p.peek().Lexeme)
}

type ParseError struct {
	error string
}

func (e ParseError) Error() string {
	return e.error
}

func NewParseError(message string) ParseError {
	return ParseError{error: "Parser Error: " + message}
}
