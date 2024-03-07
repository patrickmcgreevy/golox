package parser

import (
	"fmt"
	"strings"
)

type ASTNode interface {
	String() string
}

type Expr = ASTNode
type Statement = ASTNode

type Assign struct {
	Name  Token
	Value Expr
}

func (e Assign) String() string {
	return fmt.Sprintf("%s = %s", e.Name.Lexeme, e.Value.String())
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (s Binary) String() string {
	return fmt.Sprintf("%s %v %v", s.Operator.Token_type.String(), s.Left, s.Right)
}

type Call struct {
	Callee Expr
	Paren  Token
	Args   []Expr
}

func (e Call) String() string {
	return fmt.Sprintf("CALL %v (%v)", e.Callee, e.Args)
}

type Get struct {
	Object Expr
	Name   Token
}

func (e Get) String() string {
	return fmt.Sprintf("GET %v.%s", e.Object, e.Name.Lexeme)
}

type Grouping struct {
	Expr Expr
}

func (e Grouping) String() string {
	return fmt.Sprintf("(%s)", e.Expr.String())
}

type Literal struct {
	Value any
}

func (e Literal) String() string {
	return fmt.Sprintf("%v", e.Value)
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e Logical) String() string {
	return fmt.Sprintf("%s %s %s",
		e.Left.String(),
		e.Operator.Token_type.String(),
		e.Right.String(),
	)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (e Unary) String() string {
	return fmt.Sprintf("%s %s", e.Operator.Token_type.String(), e.Right.String())
}

type Set struct {
	Object Expr
	Name   Token
	Value  Expr
}

func (e Set) String() string {
	return fmt.Sprintf("%s.%s = %s",
		e.Object.String(),
		e.Name.Token_type.String(),
		e.Value.String(),
	)
}

type Super struct {
	Keyword Token
	Method  Token
}

func (e Super) String() string {
	return fmt.Sprintf("%s.%s", e.Keyword.Token_type.String(), e.Method.Lexeme)
}

type This struct {
	Keyword Token
}

func (e This) String() string {
	return fmt.Sprint(e.Keyword.Token_type.String())
}

type Variable struct {
	Name Token
}

func (e Variable) String() string {
	return fmt.Sprint(e.Name.Lexeme)
}

type Class struct {
	Name        Token
	Methods     []Function
	ParentClass *Variable
}

func (e Class) String() string {
	str := strings.Builder{}
	for _, v := range e.Methods {
		str.WriteString("\n")
		str.WriteString(v.String())
	}
	return fmt.Sprintf("CLASS %s(%s) {\n%s",
		e.Name.Lexeme,
		e.ParentClass.String(),
		str.String(),
	)
}

type Block struct {
	Statements []Statement
}

func (s Block) String() string {
	str := strings.Builder{}
	for _, v := range s.Statements {
		str.WriteString("\n")
		str.WriteString(v.String())
	}

	return fmt.Sprintf("[%s]", str.String())
}

type ExpressionStmt struct {
	Val Expr
}

func (s ExpressionStmt) String() string {
	return fmt.Sprintf("ExpressionStmt(%v)", s.Val)
}

type Function struct {
	Name   Token
	Params []Token
	Body   []Statement
}

func (e Function) String() string {
	str := strings.Builder{}
	for _, v := range e.Body {
		str.WriteString("\n")
		str.WriteString(v.String())
	}

	return fmt.Sprintf("%s (%v){%s}", e.Name.Lexeme, e.Params, str.String())
}

type If struct {
	Conditional Expr
	If_stmt     Statement
	Else_stmt   Statement
}

func (s If) String() string {
	str := strings.Builder{}
	if s.Else_stmt != nil {
		str.WriteString(fmt.Sprintf(" ELSE {\n%s\n}", s.Else_stmt.String()))
	}

	return fmt.Sprintf(
		"IF %s {\n%s\n}%s",
		s.Conditional.String(),
		s.If_stmt.String(),
		str.String(),
	)
}

type Print struct {
	Val Expr
}

func (s Print) String() string {
	return fmt.Sprintf("PRINT %s", s.Val.String())
}

type Return struct {
	Return_expr Expr
}

func (s Return) String() string {
	return fmt.Sprintf("RETURN %s", s.Return_expr.String())
}

type Var struct {
	Initializer Expr
	Name        Token
}

func (s Var) String() string {
	return fmt.Sprintf("VAR %s = %s", s.Name.Lexeme, s.Initializer.String())
}

type While struct {
	Conditional Expr
	Stmt        Statement
}

func (s While) String() string {
	return fmt.Sprintf("WHILE (%s) %s", s.Conditional.String(), s.Stmt.String())
}
