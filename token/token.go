package token


type TokenType int

const (
    // Single-character tokens
    LEFT_PAREN TokenType = iota
    RIGHT_PAREN
    LEFT_BRACE
    RIGHT_BRACE
    COMMA
    DOT
    MINUS
    PLUS
    SEMICOLON
    SLASH
    STAR

    // One or two character tokens
    BANG
    BANG_EQUAL
    EQUAL
    EQUAL_EQUAL
    GREATER
    GREATER_EQUAL
    LESS
    LESS_EQUAL

    // Literals
    IDENTIFIER
    STRING
    NUMBER

    // Keywords
    AND
    CLASS
    ELSE
    FALSE
    FUN
    IF
    NIL
    OR
    PRINT
    RETURN
    SUPER
    THIS
    TRUE
    VAR
    WHILE

    EOF
)

type Literal interface {
    get_literal_as_string() string
}

type Token struct {
    Token_type TokenType
    Lexeme string
    Literal *string
    Line int
}

func (t Token) toString() string {
    return string(t.Token_type) + " " + t.Lexeme + " " + *t.Literal
}

var KeywordMap = map[string]TokenType{
    "and": AND,
    "class": CLASS,
    "else": ELSE,
    "false": FALSE,
    "fun": FUN,
    "if": IF,
    "nil": NIL,
    "or": OR,
    "print": PRINT,
    "return": RETURN,
    "super": SUPER,
    "this": THIS,
    "true": TRUE,
    "var": VAR,
    "while": WHILE,
}



