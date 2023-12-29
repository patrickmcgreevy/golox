package parser

import (
	"golox/errorhandling"
	"golox/expression"
	"golox/statement"
	"golox/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) Parser {
	return Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() []statement.Statement {
    var statements []statement.Statement
	at_end := p.IsAtEnd()
	for !at_end {
		stmt, err := p.declaration()
		if err != nil {
			errorhandling.RuntimeError(err)
			return nil
		}
		statements = append(statements, stmt)
		at_end = p.IsAtEnd()
	}

	return statements
}

func (p *Parser) declaration() (statement.Statement, *ParseError) {
	var stmt statement.Statement
	var err *ParseError
	if p.match(scanner.VAR) {
		stmt, err = p.varDeclaration()
	} else {
		stmt, err = p.statement()
	}
	if err != nil {
		p.syncronize()
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) statement() (statement.Statement, *ParseError) {
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}
	if p.match(scanner.LEFT_BRACE) {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return statement.NewBlockStmt(stmts), nil
	}

	if p.match(scanner.IF) {
		return p.ifStatement()
	}

	if p.match(scanner.WHILE) {
		return p.whileStatement()
	}

	if p.match(scanner.FOR) {
		return p.forStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() (statement.Statement, *ParseError) {
	var initializer_stmt statement.Statement
	var conditional_expr expression.Expr
	var increment_expression expression.Expr
	var loop_stmt statement.Statement

	_, err := p.consume(scanner.LEFT_PAREN, "Expected '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	if p.match(scanner.VAR) {
		initializer_stmt, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else if p.match(scanner.SEMICOLON) {
		initializer_stmt = nil
	} else {
		initializer_stmt, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}
	if p.match(scanner.SEMICOLON) {
		conditional_expr = nil
	} else {
		conditional_expr, err = p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(scanner.SEMICOLON, "Expected ';' after conditional expression.")
		if err != nil {
			return nil, err
		}
	}

	if p.match(scanner.RIGHT_PAREN) {
		increment_expression = nil
	} else {
		increment_expression, err = p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(scanner.RIGHT_PAREN, "Expected ')' after expression")
		if err != nil {
			return nil, err
		}
	}

	loop_stmt, err = p.statement()
	if err != nil {
		return nil, err
	}
	if increment_expression != nil {
		body := []statement.Statement{loop_stmt, statement.NewExpressionStmt(increment_expression)}
		loop_stmt = statement.NewBlockStmt(body)
	}
	if conditional_expr == nil {
		conditional_expr = expression.Literal{Value: true}
	}
	var body statement.Statement = statement.NewWhileStmt(conditional_expr, loop_stmt)
	if initializer_stmt != nil {
		tmp := []statement.Statement{initializer_stmt, body}
		body = statement.NewBlockStmt(tmp)
	}

	return body, nil // for stmt

}

func (p *Parser) whileStatement() (statement.Statement, *ParseError) {
	var err *ParseError
	_, err = p.consume(scanner.LEFT_PAREN, "Expected '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	conditional_stmt, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.RIGHT_PAREN, "Expected ')' after conditional expression.")
	if err != nil {
		return nil, err
	}

	while_body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return statement.NewWhileStmt(conditional_stmt, while_body), nil
}

func (p *Parser) ifStatement() (statement.Statement, *ParseError) {
	var else_stmt statement.Statement

	_, err := p.consume(scanner.LEFT_PAREN, "Expected '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.RIGHT_PAREN, "Expected ')' after expression.")
	if err != nil {
		return nil, err
	}

	if_stmt, err := p.statement()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.ELSE) {
		else_stmt, err = p.statement()
	}
	if err != nil {
		return nil, err
	}

	return statement.NewIfStatement(expr, if_stmt, else_stmt), nil
}

func (p *Parser) printStatement() (statement.Statement, *ParseError) {
	var stmt statement.Statement
	expr, err := p.expression()
	if err != nil {
		return stmt, err
	}
	_, err = p.consume(scanner.SEMICOLON, "Expected ';' after value.")
	if err != nil {
		// The semi colon after "false" is being consumed. I expect it's a scanning bug.
		return stmt, err
	}

	return statement.NewPrintStmt(expr), nil
}

func (p *Parser) expressionStatement() (statement.Statement, *ParseError) {
	var stmt statement.Statement
	expr, err := p.expression()
	if err != nil {
		return stmt, err
	}
	_, err = p.consume(scanner.SEMICOLON, "Expected ';' after statement.")
	if err != nil {
		return stmt, err
	}

	return statement.NewExpressionStmt(expr), nil
}

func (p *Parser) block() ([]statement.Statement, *ParseError) {
	var statements []statement.Statement

	for !p.check(scanner.RIGHT_BRACE) && !p.IsAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	_, err := p.consume(scanner.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) varDeclaration() (statement.Statement, *ParseError) {
	var initializer expression.Expr
	var err *ParseError

	name, err := p.consume(scanner.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	if p.match(scanner.EQUAL) {
		initializer, err = p.expression()
	} else {
		initializer = nil
	}

	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return statement.NewVarStmt(name, initializer), nil
}

func (p *Parser) expression() (expression.Expr, *ParseError) {
	return p.assignment()
	// expr, err := p.assignment()
	// if error
	// p.equality

	//    expr, err := p.equality()
	//    if err != nil {
	//        return expression.Unary{}, err
	//    }
	// return expr, nil
}

func (p *Parser) assignment() (expression.Expr, *ParseError) {
	left, err := p.logic_or()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.EQUAL) {
		right, err := p.assignment()
		if err != nil {
			return nil, err
		}

		val, ok := left.(expression.Variable)
		if !ok {
			err := NewParseError("Left side of assignment must be a variable.")
			return nil, &err
		}
		return expression.Assign{Name: val.GetToken(), Value: right}, nil
	}

	return left, nil
}

func (p *Parser) logic_or() (expression.Expr, *ParseError) {
	left, err := p.logic_and()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.OR) {
		op := p.previous()
		right, err := p.logic_or()
		if err != nil {
			return nil, err
		}

		return expression.NewLogical(left, op, right), nil
	}

	return left, nil
}

func (p *Parser) logic_and() (expression.Expr, *ParseError) {
	left, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.AND) {
		op := p.previous()
		right, err := p.logic_and()
		if err != nil {
			return nil, err
		}

		return expression.NewLogical(left, op, right), nil
	}

	return left, nil
}

func (p *Parser) equality() (expression.Expr, *ParseError) {
	// First binary expression
	prefix, err := p.comparison()
	if err != nil {
		return expression.Unary{}, err
	}
	// Recursive case:
	//  Current token is an equality Operator
	//  Consume the equality operator, increment the token counter, and return a binary expression with Left: prefix, Operator: op, Right: p.equality()
	// p.current += 1
	if m := p.match(scanner.EQUAL_EQUAL, scanner.BANG_EQUAL); m {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return expression.Unary{}, nil
		}
		return expression.Binary{Left: prefix, Operator: op, Right: right}, nil
	}

	// Base case:
	//  current token is not an equality Operator
	// Return prefix
	return prefix, nil
}

func (p *Parser) comparison() (expression.Expr, *ParseError) {
	prefix, err := p.term()
	if err != nil {
		return expression.Unary{}, err
	}

	if p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		op := p.previous()
		right, err := p.comparison()
		if err != nil {
			return expression.Unary{}, nil
		}
		return expression.Binary{Left: prefix, Operator: op, Right: right}, nil
	}

	return prefix, nil
}

// term           → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() (expression.Expr, *ParseError) {
	prefix, err := p.factor()
	if err != nil {
		return expression.Unary{}, err
	}

	if p.match(scanner.MINUS, scanner.PLUS) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return expression.Unary{}, err
		}
		return expression.Binary{Left: prefix, Operator: op, Right: right}, nil
	}

	return prefix, nil
}

// factor         → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() (expression.Expr, *ParseError) {
	prefix, err := p.unary()
	if err != nil {
		return expression.Unary{}, err
	}

	if p.match(scanner.STAR, scanner.SLASH) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return expression.Unary{}, err
		}
		return expression.Binary{Left: prefix, Operator: op, Right: right}, nil
	}

	return prefix, nil
}

// unary          → ( "!" | "-" ) unary
//
//	| primary ;
func (p *Parser) unary() (expression.Expr, *ParseError) {
	var right expression.Expr
	var err *ParseError

	// prefix := p.advance()
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right, err = p.unary()
		if err != nil {
			return expression.Unary{}, err
		}
		return expression.Unary{Operator: operator, Right: right}, nil
	}

	primary, err := p.call()

	if err != nil {
		return expression.Unary{}, err
	}

	return primary, nil
}

func (p *Parser) call() (expression.Expr, *ParseError) {
	var expr expression.Expr
	var err *ParseError

	expr, err = p.primary()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.LEFT_PAREN) {
		expr, err = p.add_args(expr)
		if err != nil {
			return nil, err
		}
	}

	return expr, nil
}

func (p *Parser) add_args(expr expression.Expr) (expression.Expr, *ParseError) {
	paren := p.previous()
	if p.match(scanner.RIGHT_PAREN) {
		return expression.NewCall(expr, paren, nil), nil
	}

	args, err := p.arguments()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(scanner.RIGHT_PAREN, "Expected ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return expression.NewCall(expr, paren, args), nil // call expr with args, nil
}

func (p *Parser) arguments() ([]expression.Expr, *ParseError) {
	var args []expression.Expr
	var err *ParseError
	for cur_arg, err := p.expression(); err == nil; cur_arg, err = p.expression() {
		if len(args) >= 255 {
			p.error(p.peek(), "Can't have more than 255 argumens.")
		}
		args = append(args, cur_arg)
		if !p.match(scanner.COMMA) {
			return args, nil
		}
	}
	return nil, err
}

// primary        → NUMBER | STRING | "true" | "false" | "nil" | IDENTIFIER | (expression)
//
//	| "(" expression ")"
func (p *Parser) primary() (expression.Expr, *ParseError) {
	var err *ParseError
	var expr expression.Expr
	if p.match(scanner.FALSE) {
		return expression.Literal{Value: false}, nil
	}
	if p.match(scanner.TRUE) {
		return expression.Literal{Value: true}, nil
	}
	if p.match(scanner.NIL) {
		return expression.Literal{Value: nil}, nil
	}
	if p.match(scanner.STRING, scanner.NUMBER) {
		return expression.Literal{Value: p.previous().Literal}, nil
	}
	if p.match(scanner.LEFT_PAREN) {
		expr, err = p.expression()
		_, err = p.consume(scanner.RIGHT_PAREN, "Expected right paren!")
		if err != nil {
			return expression.Unary{}, err
		}

		return expression.Grouping{Expr: expr}, nil
	}
	if p.match(scanner.IDENTIFIER) {
		return expression.NewVariableExpression(p.previous()), nil
	}
	parse_error := p.error(p.peek(), "Expect expression.")

	return expression.Unary{}, &parse_error
}

// func (p *Parser) identifier() (expression.Expr, *ParseError)

func (p *Parser) syncronize() {
	p.advance()

	for !p.IsAtEnd() {
		if p.previous().Token_type == scanner.SEMICOLON {
			return
		}
		t := p.peek().Token_type
		switch t {
		case scanner.CLASS, scanner.FUN, scanner.VAR, scanner.FOR, scanner.IF, scanner.WHILE, scanner.PRINT, scanner.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(token_type ...scanner.TokenType) bool {
	for _, v := range token_type {
		if p.check(v) {
			p.advance()
			return true
		}
	}

	return false
}

func (p Parser) check(token_type scanner.TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.peek().Token_type == token_type
}

func (p *Parser) advance() scanner.Token {
	if !p.IsAtEnd() {
		p.current += 1
	}

	ret := p.previous()

	return ret
}

func (p *Parser) consume(tokenType scanner.TokenType, message string) (scanner.Token, *ParseError) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	error := p.error(p.peek(), message)

	return scanner.Token{}, &error
}

func (p Parser) error(token scanner.Token, message string) ParseError {
	if token.Token_type == scanner.EOF {
		errorhandling.Report(token.Line, " at end ", message)
	} else {
		errorhandling.Report(token.Line, " at "+token.Lexeme, message)
	}
	return NewParseError(message)
}

func (p Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

func (p Parser) IsAtEnd() bool {
	return p.peek().Token_type == scanner.EOF
}

type ParseError struct {
	error string
}

func (e *ParseError) Error() string {
	return e.error
}

func NewParseError(message string) ParseError {
	return ParseError{error: "Parser Error: " + message}
}
