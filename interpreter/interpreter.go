package interpreter

import (
	"golox/expression"
	"golox/token"
	"reflect"
)

type Interpreter struct {
	val any
}

func (v *Interpreter) Evaluate(e expression.Expr) any {
	e.Accept(v)

	return v.val
}

func (v Interpreter) isTruthy(val any) bool {
	if val == nil {
		return false
	}

	b, ok := val.(bool)
	if ok {
		return b
	}

	return true
}

func (v Interpreter) isEqual(left, right any) bool {
	return reflect.DeepEqual(left, right)
}

func (v *Interpreter) VisitAssign(e expression.Assign) {
}

func (v *Interpreter) VisitBinary(e expression.Binary) {
	left := v.Evaluate(e.Left)
	right := v.Evaluate(e.Right)
	switch e.Operator.Token_type {
	case token.MINUS:
		l, r := left.(float64), right.(float64)
		v.val = l - r
	case token.PLUS:
		l_float, l_ok := left.(float64)
		r_float, r_ok := right.(float64)
		if l_ok && r_ok {
			v.val = l_float + r_float
		}

		l_str, l_ok := left.(*string)
		r_str, r_ok := right.(*string)

		if l_ok && r_ok {
			v.val = *l_str + *r_str
		}

	case token.SLASH:
		l, r := left.(float64), right.(float64)
		v.val = l / r

	case token.STAR:
		l, r := left.(float64), right.(float64)
		v.val = l * r

	case token.GREATER:
		l, r := left.(float64), right.(float64)
		v.val = l > r

	case token.GREATER_EQUAL:
		l, r := left.(float64), right.(float64)
		v.val = l >= r

	case token.LESS:
		l, r := left.(float64), right.(float64)
		v.val = l < r

	case token.LESS_EQUAL:
		l, r := left.(float64), right.(float64)
		v.val = l <= r

	case token.BANG_EQUAL:
		v.val = !v.isEqual(left, right)

	case token.EQUAL_EQUAL:
		v.val = v.isEqual(left, right)
	}
}

func (v *Interpreter) VisitGrouping(e expression.Grouping) {
	v.val = v.Evaluate(e.Expr)
}

func (v *Interpreter) VisitLiteral(e expression.Literal) {
	v.val = e.Value
}

func (v *Interpreter) VisitUnary(e expression.Unary) {
	right := v.Evaluate(e.Right)
	t := e.Operator.Token_type
	switch t {
	case token.MINUS:
		r := right.(float64)
		v.val = -r
	case token.BANG:
		v.val = v.isTruthy(right)
	}
}
