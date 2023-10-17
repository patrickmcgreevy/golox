package parser

import (
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
	p.equality()
}

func (p *Parser) equality() expression.Expr {
	prefix := p.comparison()
    comparisons := []expression.Expr {}

	for next_token := p.tokens[p.current].Token_type; next_token == token.BANG_EQUAL || next_token == token.EQUAL_EQUAL; {

        
	}
}

func (p *Parser) comparison() expression.Expr {
}

func (p *Parser) term() expression.Expr {
}

func (p *Parser) factor() expression.Expr {
}

func (p *Parser) unary() expression.Expr {
}

func (p *Parser) primary() expression.Expr {
}
