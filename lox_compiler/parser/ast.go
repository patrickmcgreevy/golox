package parser

type ASTNode interface {
}

type Expr = ASTNode
type Statement = ASTNode


type Assign struct {
	Name  Token
	Value Expr
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type Call struct {
	Callee Expr
	Paren  Token
	Args   []Expr
}

type Get struct {
	Object Expr
	Name   Token
}

type Grouping struct {
	Expr Expr
}

type Literal struct {
	Value any
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type Unary struct {
	Operator Token
	Right    Expr
}

type Set struct {
	Object Expr
	Name   Token
	Value  Expr
}

type Super struct {
	Keyword Token
	Method  Token
}

type This struct {
	Keyword Token
}

type Variable struct {
	Name Token
}

type Class struct {
    Name Token
    Methods []Function
    ParentClass *Variable
}

type Block struct {
	statements []Statement
}

type Expression struct {
	Val Expr
}

type Function struct {
    Name Token
    Params []Token
    Body []Statement
}

type If struct {
	Conditional Expr
	If_stmt     Statement
	Else_stmt   Statement
}

type Print struct {
	Val Expr
}

type Return struct {
    Return_expr Expr
}

type Var struct {
	Initializer Expr
	Name        Token
}

type While struct {
	Conditional Expr
	Stmt        Statement
}

