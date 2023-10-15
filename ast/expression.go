package expression

import (
	"golox/token"
	"strings"
    "strconv"
)

type Expr interface {
	expand_to_string() string
}

type Assign struct {
	name  token.Token
	value Expr
}

type Binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

type Grouping struct {
	expr Expr
}

type Literal struct {
	value any
}

type Unary struct {
	operator token.Token
	right    Expr
}

func (e Assign) expand_to_string() string {
	sb := strings.Builder{}

	sb.WriteString(e.name.Lexeme)
	sb.WriteString(" = ")
	sb.WriteString(e.value.expand_to_string())

	return sb.String()
}

func (e Binary) expand_to_string() string {
	sb := strings.Builder{}

    sb.WriteString(parenthesize(e.operator.Lexeme, e.left, e.right))

    return sb.String()
}

func (e Grouping) expand_to_string() string {
    return parenthesize("grouping", e.expr)
}

func (e Literal) expand_to_string() string {
    if e.value == nil {
        return "nil"
    }

    switch v:=e.value.(type) {
    case string:
        return v
    case int:
        ret:= strconv.FormatInt(v, 10){
            panic("Couldn't convert to int!")
        }
        return ret
    default:
        panic("Unexpected type")
    }
}

func parenthesize(name string, exprs ...Expr) string {
	sb := strings.Builder{}

	sb.WriteString("(")
	sb.WriteString(name + " ")
	for _, v := range exprs {
		sb.WriteString(v.expand_to_string())
	}
	sb.WriteString(" )")

	return sb.String()
}
