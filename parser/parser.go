package parser

import (
	"errors"
	"golox/expression"
	"golox/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func newParser(tokens []token.Token) Parser {
	return Parser{tokens: tokens, current: 0}
}

func (p *Parser) expression() expression.Expr {
	return p.equality()
}

func (p *Parser) equality() expression.Expr {
	// First binary expression
	prefix := p.comparison()
	// Recursive case:
	//  Current token is an equality Operator
	//  Consume the equality operator, increment the token counter, and return a binary expression with Left: prefix, Operator: op, Right: p.equality()
	// p.current += 1
	if m := p.match(token.EQUAL_EQUAL, token.BANG_EQUAL); m {
		return expression.Binary{Left: prefix, Operator: p.previous(), Right: p.equality()}
	}

	// Base case:
	//  current token is not an equality Operator
	// Return prefix
	return prefix
}

func (p *Parser) comparison() expression.Expr {
    prefix := p.term()

    if p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
        return expression.Binary{Left:prefix, Operator: p.previous(), Right: p.comparison()}
    }

    return prefix
}

// term           → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() expression.Expr {
    prefix := p.factor()

    if p.match(token.MINUS, token.BANG, token.PLUS) {
        return expression.Binary{Left: prefix, Operator: p.previous(), Right: p.term()}
    }

    return prefix
}

//factor         → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() expression.Expr {
    prefix := p.unary()

    if p.match(token.STAR, token.SLASH) {
        return expression.Binary{Left: prefix, Operator: p.previous(), Right: p.factor()}
    }

    return prefix
}

func (p *Parser) unary() expression.Expr {
    var right expression.Expr

    prefix := p.advance()
    if p.match(token.BANG, token.MINUS) {
        right = p.unary()
    } else {
        right = p.primary()
    }

    return expression.Unary{Operator: prefix, Right: right}
}

// primary        → NUMBER | STRING | "true" | "false" | "nil"
//               | "(" expression ")"
func (p *Parser) primary() expression.Expr {
    if p.match(token.LEFT_PAREN) {
        expr := p.expression()
        p.advance()

        return expression.Grouping{Expr: expr}
    }

    return expression.Literal{Value: p.previous().Literal}

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
	if p.isAtEnd() {
		return false
	}
	return p.peek().Token_type == token_type
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current += 1
	}

	ret := p.previous()

	return ret
}

func (p Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p Parser) isAtEnd() bool {
	return p.peek().Token_type == token.EOF
}
