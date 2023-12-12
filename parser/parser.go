package parser

import (
	"golox/errorhandling"
	"golox/expression"
	"golox/statement"
	"golox/token"

	"golang.org/x/exp/slices"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) Parser {
	return Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse()  []statement.Statement {
    statements := []statement.Statement{}
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
    if p.match(token.VAR) {
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
    if p.match(token.PRINT) {
        return p.printStatement()
    } 
    if p.match(token.LEFT_BRACE) {
        stmts, err := p.block()
        if err != nil {
            return nil, err
        }
        return statement.NewBlockStmt(stmts), nil
    }

    if p.match(token.IF) {
        return p.ifStatement()
    }

    if p.match(token.WHILE) {
        return p.whileStatement()
    }

    return p.expressionStatement()
}

func (p *Parser) forStatement() (statement.Statement, *ParseError) {
    var initializer_stmt statement.Statement
    var conditional_expr expression.Expr
    var increment_expression expression.Expr
    var loop_stmt statement.Statement

    _, err := p.consume(token.LEFT_PAREN, "Expected '(' after 'for'.")
    if err != nil {
        return nil , err
    }

    if p.match(token.VAR) {
        initializer_stmt, err = p.varDeclaration()
        if err != nil {
            return nil, err
        }
    } else if p.match(token.SEMICOLON) {
        initializer_stmt = nil
    } else {
        initializer_stmt, err = p.expressionStatement()
        if err != nil {
            return nil, err
        }
    }
    if p.match(token.SEMICOLON) {
        conditional_expr = nil
    } else {
        conditional_expr, err = p.expression()
        if err != nil {
            return nil, err
        }
        _, err = p.consume(token.SEMICOLON, "Expected ';' after conditional expression.")
        if err != nil {
            return nil, err
        }
    }

    if p.match(token.RIGHT_PAREN) {
        increment_expression = nil
    } else {
        increment_expression, err = p.expression()
        if err != nil {
            return nil, err
        }
        _, err = p.consume(token.RIGHT_PAREN, "Expected ')' after expression")
        return nil, err
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
    body := statement.NewWhileStmt(conditional_expr, loop_stmt)
    if initializer_stmt != nil {
        tmp := []statement.Statement{initializer_stmt, body}
        body = statement.NewBlockStmt(tmp)
    }

    return // for stmt


}

func (p *Parser) whileStatement() (statement.Statement, *ParseError) {
    var err *ParseError
    _, err = p.consume(token.LEFT_PAREN, "Expected '(' after 'while'.")
    if err != nil {
        return nil, err
    }

    conditional_stmt, err := p.expression()
    if err != nil {
        return nil, err
    }

    _, err = p.consume(token.RIGHT_PAREN, "Expected ')' after conditional expression.")
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

    _, err := p.consume(token.LEFT_PAREN, "Expected '(' after 'if'.")
    if err != nil {
        return nil, err
    }

    expr, err := p.expression()
    if err != nil {
        return nil, err
    }

    _, err = p.consume(token.RIGHT_PAREN, "Expected ')' after expression.")
    if err != nil {
        return nil, err
    }

    if_stmt, err := p.statement()
    if err != nil {
        return nil, err
    }

    if p.match(token.ELSE) {
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
    _, err = p.consume(token.SEMICOLON, "Expected ';' after value.")
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
    _, err = p.consume(token.SEMICOLON, "Expected ';' after statement.")
    if err != nil {
        return stmt, err
    }

    return statement.NewExpressionStmt(expr), nil
}

func (p *Parser) block() ([]statement.Statement, *ParseError) {
    statements := []statement.Statement{}

    for !p.check(token.RIGHT_BRACE) && !p.IsAtEnd() {
        stmt, err := p.declaration()
        if err != nil {
            return nil, err
        }
        statements = append(statements, stmt)
    }
    _, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
    if err != nil {
        return nil, err
    }

    return statements, nil
}

func (p *Parser) varDeclaration() (statement.Statement, *ParseError) {
    var initializer expression.Expr
    var err *ParseError

    name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
    if err != nil {
        return nil, err
    }
    if p.match(token.EQUAL) {
        initializer, err = p.expression()
    } else {
        initializer = nil
    }

    if err != nil {
        return nil, err
    }
    _, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
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

    if p.match(token.EQUAL) {
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

    if p.match(token.OR) {
        op := p.previous()
        right, err := p.logic_or()
        if err != nil {
            return nil, err
        }

        return expression.NewLogical(left, op,  right), nil
    }

    return left, nil
}

func (p *Parser) logic_and() (expression.Expr, *ParseError) {
    left, err := p.equality()
    if err != nil {
        return nil, err
    }

    if p.match(token.AND) {
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
	if m := p.match(token.EQUAL_EQUAL, token.BANG_EQUAL); m {
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

	if p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
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
func (p *Parser) term() (expression.Expr, *ParseError){
	prefix, err := p.factor()
    if err != nil {
        return expression.Unary{}, err
    }

	if p.match(token.MINUS, token.PLUS) {
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
func (p *Parser) factor() (expression.Expr, *ParseError ){
	prefix, err := p.unary()
    if err != nil {
        return expression.Unary{}, err
    }

	if p.match(token.STAR, token.SLASH) {
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
//                | primary ;
func (p *Parser) unary() (expression.Expr, *ParseError) {
	var right expression.Expr
	var err *ParseError

	// prefix := p.advance()
	if p.match(token.BANG, token.MINUS) {
        operator := p.previous()
		right, err = p.unary()
        if err != nil {
            return expression.Unary{}, err
        }
        return expression.Unary{Operator: operator, Right: right}, nil
    }

    primary, err := p.primary()

	if err != nil {
		return expression.Unary{}, err
	}

    return primary, nil
}

// primary        → NUMBER | STRING | "true" | "false" | "nil" | IDENTIFIER | (expression)
//
//	| "(" expression ")"
func (p *Parser) primary() (expression.Expr, *ParseError) {
    var err *ParseError
    var expr expression.Expr
	if p.match(token.FALSE) {
		return expression.Literal{Value: false}, nil
	}
	if p.match(token.TRUE) {
		return expression.Literal{Value: true}, nil
	}
	if p.match(token.NIL) {
		return expression.Literal{Value: nil}, nil
	}
	if p.match(token.STRING, token.NUMBER) {
		return expression.Literal{Value: p.previous().Literal}, nil
	}
	if p.match(token.LEFT_PAREN) {
		expr, err = p.expression()
		_, err = p.consume(token.RIGHT_PAREN, "Expected right paren!")
		if err != nil {
			return expression.Unary{}, err
		}

		return expression.Grouping{Expr: expr}, nil
	}
    if p.match(token.IDENTIFIER) {
        return expression.NewVariableExpression(p.previous()), nil
    }
	parse_error := p.error(p.peek(), "Expect expression.")

	return expression.Unary{}, &parse_error
}

// func (p *Parser) identifier() (expression.Expr, *ParseError)

func (p *Parser) syncronize() {
	p.advance()

	for !p.IsAtEnd() {
		if p.previous().Token_type == token.SEMICOLON {
			return
		}
		t := p.peek().Token_type
		switch t {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(token_type ...token.TokenType) bool {
	for _, v := range token_type {
		if p.check(v) {
			p.advance()
			return true
		}
	}

	return false
}

func (p Parser) check(token_type token.TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.peek().Token_type == token_type
}

func (p *Parser) advance() token.Token {
	if !p.IsAtEnd() {
		p.current += 1
	}

	ret := p.previous()

	return ret
}

func (p *Parser) consume(tokenType token.TokenType, message string) (token.Token, *ParseError) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	error := p.error(p.peek(), message)

	return token.Token{}, &error
}

func (p Parser) error(token token.Token, message string) ParseError {
	errorhandling.Error(token, message)

	return NewParseError(message)
}

func (p Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p Parser) IsAtEnd() bool {
	return p.peek().Token_type == token.EOF
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
