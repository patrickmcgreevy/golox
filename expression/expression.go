package expression

import (
	"fmt"
	"golox/scanner"
	"strconv"
	"strings"
)

type Expr interface {
	Expand_to_string() string
	Accept(v Visitor)
}

func parenthesize(name string, exprs ...Expr) string {
	sb := strings.Builder{}

	sb.WriteString("(")
	sb.WriteString(name)
	for _, v := range exprs {
		sb.WriteString(" " + v.Expand_to_string())
	}
	sb.WriteString(")")

	return sb.String()
}

type Assign struct {
	Name  scanner.Token
	Value Expr
}

func (e Assign) Accept(v Visitor) {
	v.VisitAssign(e)
}

type Binary struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (e Binary) Accept(v Visitor) {
	v.VisitBinary(e)
}

type Call struct {
	Callee Expr
	Paren  scanner.Token
	Args   []Expr
}

func NewCall(callee Expr, paren scanner.Token, args []Expr) Call {
	return Call{Callee: callee, Paren: paren, Args: args}
}

func (e Call) Accept(v Visitor) {
	v.VisitCall(e)
}

type Get struct {
	Object Expr
	Name   scanner.Token
}

func (e Get) Accept(v Visitor) {
	v.VisitGet(e)
}

func (e Get) Expand_to_string() string {
	return e.Object.Expand_to_string() + "." + e.Name.Lexeme
}

type Grouping struct {
	Expr Expr
}

func (e Grouping) Accept(v Visitor) {
	v.VisitGrouping(e)
}

type Literal struct {
	Value any
}

func (e Literal) Accept(v Visitor) {
	v.VisitLiteral(e)
}

type Logical struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func NewLogical(left Expr, operator scanner.Token, right Expr) Logical {
	return Logical{Left: left, Operator: operator, Right: right}
}

func (e Logical) Accept(v Visitor) {
	v.VisitLogical(e)
}

func (e Logical) Expand_to_string() string {
	sb := strings.Builder{}

	sb.WriteString(e.Left.Expand_to_string())
	sb.WriteString(e.Operator.Lexeme)
	sb.WriteString(e.Right.Expand_to_string())

	return sb.String()
}

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

func (e Unary) Accept(v Visitor) {
	v.VisitUnary(e)
}

type Set struct {
	Object Expr
	Name   scanner.Token
	Value  Expr
}

func (e Set) Expand_to_string() string {
	return fmt.Sprintf("%s.%s = %s",
		e.Object.Expand_to_string(),
		e.Name.Lexeme,
		e.Value.Expand_to_string(),
	)
}

func (e Set) Accept(v Visitor) {
	v.VisitSet(e)
}

type Variable struct {
	name scanner.Token
}

func NewVariableExpression(name scanner.Token) Variable {
	return Variable{name: name}
}

func (e Variable) Accept(v Visitor) {
	v.VisitVariable(e)
}

func (e Variable) GetName() string {
	return e.name.Lexeme
}

func (e Variable) GetToken() scanner.Token {
	return e.name
}

func (e Assign) Expand_to_string() string {
	sb := strings.Builder{}

	sb.WriteString(e.Name.Lexeme)
	sb.WriteString(" = ")
	sb.WriteString(e.Value.Expand_to_string())

	return sb.String()
}

func (e Binary) Expand_to_string() string {
	sb := strings.Builder{}

	sb.WriteString(parenthesize(e.Operator.Lexeme, e.Left, e.Right))

	return sb.String()
}

func (e Call) Expand_to_string() string {
	sb := strings.Builder{}
	sb.WriteString(e.Callee.Expand_to_string())
	sb.WriteString("(arguments ")
	for _, i := range e.Args {
		sb.WriteString(i.Expand_to_string())
	}
	sb.WriteString(")")

	return sb.String()
}

func (e Grouping) Expand_to_string() string {
	return parenthesize("grouping", e.Expr)
}

func (e Literal) Expand_to_string() string {
	if e.Value == nil {
		return "nil"
	}

	switch v := e.Value.(type) {
	case string:
		return v
	case *string:
		return *v
	case int:
		return strconv.Itoa(v)
	case float64:
		// return strconv.FormatFloat(v, 'f', 32, 64)
		s := fmt.Sprintf("%f", v)
		return s
	default:
		panic("Unexpected type")
	}
}

func (e Unary) Expand_to_string() string {
	return parenthesize(e.Operator.Lexeme, e.Right)
}

func (e Variable) Expand_to_string() string {
	return e.name.Lexeme
}

type Visitor interface {
	VisitAssign(e Assign)
	VisitBinary(e Binary)
	VisitCall(e Call)
	VisitGet(e Get)
	VisitGrouping(e Grouping)
	VisitLiteral(e Literal)
	VisitLogical(e Logical)
	VisitSet(e Set)
	VisitUnary(e Unary)
	VisitVariable(e Variable)
}

type ExpressionStringVisitor struct {
	expr_string_builder strings.Builder
}

func (v *ExpressionStringVisitor) As_string() string {
	return v.expr_string_builder.String()
}

func (v *ExpressionStringVisitor) Reset() {
	v.expr_string_builder.Reset()
}

func (v *ExpressionStringVisitor) parenthesize(name string, exprs ...Expr) {
	v.expr_string_builder.WriteString("(")
	v.expr_string_builder.WriteString(name)
	for _, val := range exprs {
		v.expr_string_builder.WriteString(" ")
		val.Accept(v)
	}
	v.expr_string_builder.WriteString(")")
}

func (v *ExpressionStringVisitor) VisitAssign(e Assign) {
	v.expr_string_builder.WriteString(e.Name.Lexeme)
	v.expr_string_builder.WriteString(" = ")
	e.Value.Accept(v)
}

func (v *ExpressionStringVisitor) VisitBinary(e Binary) {
	v.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (v *ExpressionStringVisitor) VisitCall(e Call) {
	e.Callee.Accept(v)
	v.parenthesize("arguments", e.Args...)
}

func (v *ExpressionStringVisitor) VisitGet(e Get) {
	e.Object.Accept(v)
	v.expr_string_builder.WriteString(".")
	v.expr_string_builder.WriteString(e.Name.Lexeme)
}

func (v *ExpressionStringVisitor) VisitGrouping(e Grouping) {
	v.parenthesize("grouping", e.Expr)
}

func (v *ExpressionStringVisitor) VisitLiteral(e Literal) {
	if e.Value == nil {
		v.expr_string_builder.WriteString("nil")
	}

	switch val := e.Value.(type) {
	case string:
		v.expr_string_builder.WriteString(val)
	case *string:
		v.expr_string_builder.WriteString(*val)
	case int:
		v.expr_string_builder.WriteString(strconv.Itoa(val))
	case float64:
		// return strconv.FormatFloat(v, 'f', 32, 64)
		s := fmt.Sprintf("%f", val)
		v.expr_string_builder.WriteString(s)
	default:
		panic("Unexpected type")
	}
}

func (v *ExpressionStringVisitor) VisitLogical(e Logical) {
	e.Left.Accept(v)
	v.expr_string_builder.WriteString(e.Operator.Lexeme)
	e.Right.Accept(v)
}

func (v *ExpressionStringVisitor) VisitSet(e Set) {
    e.Object.Accept(v)
    v.expr_string_builder.WriteString(".")
    v.expr_string_builder.WriteString(e.Name.Lexeme)
    v.expr_string_builder.WriteString(" = ")
    e.Value.Accept(v)
}

func (v *ExpressionStringVisitor) VisitUnary(e Unary) {
	v.parenthesize(e.Operator.Lexeme, e.Right)
}

func (v *ExpressionStringVisitor) VisitVariable(e Variable) {
	v.expr_string_builder.WriteString(e.name.Lexeme)
}
