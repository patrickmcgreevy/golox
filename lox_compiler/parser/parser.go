package parser

import "fmt"

type errorType string
const (
    endOfTokens errorType = "end of tokens"
    consumeFailed = ""
)

type ParseError struct {
    // errorType errorType
    msg string
}

func (e ParseError) Error() string {
    return fmt.Sprintf("parsing error: %s", e.msg)
}

type ASTType int

const (
    BinaryExpr ASTType = iota
)

type ASTNode struct {
}

// Parses tokens and returns the AST
type Parser struct {
    hadError bool
    panicMode bool
    tokens []Token
    curToken uint
    prevToken *Token
}

func NewParser(t []Token) Parser {
    return Parser{tokens: t}
}

func (p *Parser) Parse() *ParseError {
    return nil
}

func (p *Parser) expression() *ParseError {
    return nil
}

func (p *Parser) consume(t TokenType, msg string) (Token, *ParseError) {
    if p.peek().Token_type != t {
        p.errorAtCurrent()
        return Token{}, &ParseError{msg: msg}
    }

    p.advance()
    return p.previous()
}

func (p Parser) peek() Token {
    return p.tokens[p.curToken]
}

func (p *Parser) advance() {
    p.prevToken = p.current()
    for ;; {
        p.curToken += 1
        if (*p.current()).Token_type != ERROR {
            p.errorAtCurrent()
        }
        return
    }
}

func (p *Parser) current() *Token {
    return &p.tokens[p.curToken]
}

func (p Parser) previous() (Token, *ParseError) {
    if p.curToken == 0 {
        return Token{}, &ParseError{msg: "no previous token, already at first token"}
    }

    return p.tokens[p.curToken-1], nil
}

// Display an error at the current token and enter panic mode.
func (p *Parser) errorAtCurrent() {
    p.hadError = true
    if p.panicMode == true {
        return
    }
    p.panicMode = true
    fmt.Println(p.current().Lexeme)
}
