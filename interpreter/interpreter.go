package interpreter

import (
	"fmt"
	"golox/environment"
	"golox/errorhandling"
	"golox/expression"
	"golox/statement"
	"golox/token"
	"reflect"
)

type RuntimeError struct {
	error string
	tok   token.Token
}

func (e RuntimeError) Error() string {
	return e.error
}

func (e RuntimeError) GetToken() token.Token {
	return e.tok
}

func newRuntimeError(operator token.Token, message string) *RuntimeError {
	msg := fmt.Sprintf("[line %d]: %s", operator.Line, message)
	new_err := RuntimeError{error: msg, tok: operator}
	return &new_err
}

func newNumberError(operator token.Token) *RuntimeError {
	return newRuntimeError(operator, "Operand must be a number.")
}

func newOperandsError(operator token.Token) *RuntimeError {
	return newRuntimeError(operator, "Operands must be numbers.")
}

type Interpreter struct {
	val         any
	err         *RuntimeError
	environment environment.Environment
}

func NewInterpreter() Interpreter {
	return Interpreter{val: nil, err: nil, environment: environment.NewEnvironment()}
}

func (v *Interpreter) Interpret(statements []statement.Statement) {
	v.err = nil
	for _, stmt := range statements {
		err := v.execute(stmt)
		if err != nil {
			errorhandling.RuntimeError(err)
			return
		}
	}
}

func (v *Interpreter) execute(stmt statement.Statement) *RuntimeError {
	stmt.Accept(v)
	if v.err != nil {
		return v.err
	}
	return nil
}

func (v *Interpreter) Evaluate(e expression.Expr) (any, *RuntimeError) {
	e.Accept(v)

	return v.val, v.err
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
	left, err := v.Evaluate(e.Left)
	if err != nil {
		v.err = err
		return
	}
	right, err := v.Evaluate(e.Right)
	if err != nil {
		v.err = err
		return
	}
	switch e.Operator.Token_type {
	case token.MINUS:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l - r
	case token.PLUS:
		l_float, l_ok := left.(float64)
		r_float, r_ok := right.(float64)
		if l_ok && r_ok {
			v.val = l_float + r_float
			return
		}

		l_str, l_ok := left.(string)
		r_str, r_ok := right.(string)

		if l_ok && r_ok {
			v.val = l_str + r_str
			return
		}
		v.err = newRuntimeError(e.Operator, "Operands must be two numbers or two strings")

	case token.SLASH:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l / r

	case token.STAR:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l * r

	case token.GREATER:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l > r

	case token.GREATER_EQUAL:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l >= r

	case token.LESS:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l < r

	case token.LESS_EQUAL:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l <= r

	case token.BANG_EQUAL:
		v.val = !v.isEqual(left, right)

	case token.EQUAL_EQUAL:
		v.val = v.isEqual(left, right)
	}
}

func (v *Interpreter) VisitGrouping(e expression.Grouping) {
	v.val, v.err = v.Evaluate(e.Expr)
}

func (v *Interpreter) VisitLiteral(e expression.Literal) {
	str, ok := e.Value.(*string)
	if ok {
		v.val = *str
	} else {
		v.val = e.Value
	}
}

func (v *Interpreter) VisitUnary(e expression.Unary) {
	right, err := v.Evaluate(e.Right)
	if err != nil {
		v.err = err
		return
	}
	t := e.Operator.Token_type
	switch t {
	case token.MINUS:
		r, ok := right.(float64)
		if !ok {
			v.err = newNumberError(e.Operator)
		}
		v.val = -r
	case token.BANG:
		v.val = v.isTruthy(right)
	}
}

func (v *Interpreter) VisitVariable(e expression.Variable) {
	// TODO: We need global state now!!
}

func (v *Interpreter) VisitBlockStmt(stmt statement.Block) {
}
func (v *Interpreter) VisitClassStmt(stmt statement.Class) {
}
func (v *Interpreter) VisitExpressionStmt(stmt statement.Expression) {
	v.Evaluate(stmt.Val)
}
func (v *Interpreter) VisitFunctionStmt(stmt statement.Function) {
}
func (v *Interpreter) VisitIfStmt(stmt statement.If) {
}
func (v *Interpreter) VisitPrintStmt(stmt statement.Print) {
	val, err := v.Evaluate(stmt.Val)
	if err != nil {
		v.err = err
		return
	}

	fmt.Println(val)
}
func (v *Interpreter) VisitReturnStmt(stmt statement.Return) {
}

func (v *Interpreter) VisitVarStmt(stmt statement.Var) {
	var val any
    var err *RuntimeError
	if stmt.Initializer != nil {
		val, err = v.Evaluate(stmt.Initializer)
		if err != nil {
			v.err = err
			return
		}
	}
	// Create a variable and assign it to val
	v.environment.Define(stmt.Name.Lexeme, val)
}

func (v *Interpreter) VisitWhileStmt(stmt statement.While) {
}
