package expression

import (
	"fmt"
	"golox/token"
	"strconv"
	"strings"
)

type Expr interface {
	Expand_to_string() string
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
	Name  token.Token
	Value Expr
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

type Grouping struct {
	Expr Expr
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator token.Token
	Right    Expr
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

