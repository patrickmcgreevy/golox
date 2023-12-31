package interpreter

import (
	"fmt"
	"golox/errorhandling"
	"golox/expression"
	"golox/scanner"
	"golox/statement"
	"reflect"
	"time"
)

type LoxCallable interface {
	Call(interp Interpreter, args []any) any
	Arity() int
}

type RuntimeError struct {
	error string
	tok   scanner.Token
}

func (e RuntimeError) Error() string {
	return e.error
}

func (e RuntimeError) GetToken() scanner.Token {
	return e.tok
}

func newRuntimeError(operator scanner.Token, message string) *RuntimeError {
	msg := fmt.Sprintf("[line %d]: %s", operator.Line, message)
	new_err := RuntimeError{error: msg, tok: operator}
	return &new_err
}

func newNumberError(operator scanner.Token) *RuntimeError {
	return newRuntimeError(operator, "Operand must be a number.")
}

func newOperandsError(operator scanner.Token) *RuntimeError {
	return newRuntimeError(operator, "Operands must be numbers.")
}

type Interpreter struct {
	val             any
	err             *RuntimeError
	pEnvironment    *Environment
	interactiveMode bool
}

func NewInterpreter() Interpreter {
	env := NewEnvironment()

	env.Define("clock", BuiltinCallable{params: make([]string, 0), foo: func(a Interpreter, b []any) any {
		return float64(time.Now().UnixMilli())/1000
	}})
	return Interpreter{val: nil, err: nil, pEnvironment: &env, interactiveMode: false}
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

func (v *Interpreter) EnableInteractiveMode() {
	v.interactiveMode = true
}

func (v *Interpreter) DisableInteractiveMode() {
	v.interactiveMode = false
}

func (v *Interpreter) execute(stmt statement.Statement) *RuntimeError {
	stmt.Accept(v)
	if v.err != nil {
		return v.err
	}
	return nil
}

func (v *Interpreter) executeBlock(statements []statement.Statement, env Environment) {
	v.pushEnvironment(&env)
	defer v.popEnvironment()
	for _, stmt := range statements {
		v.execute(stmt)
	}
}

func (v *Interpreter) pushEnvironment(env *Environment) {
	env.SetEnclosing(v.pEnvironment)
	v.pEnvironment = env
}

func (v *Interpreter) popEnvironment() {
	parent := v.pEnvironment.GetEnclosing()
	if parent != nil {
		v.pEnvironment = parent
	}
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
	right, err := v.Evaluate(e.Value)
	if err != nil {
		v.err = err
		return
	}
	assignment_error := v.pEnvironment.Assign(e.Name.Lexeme, right)
	if assignment_error != nil {
		err = newRuntimeError(e.Name, assignment_error.Error())
		v.err = err
	}
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
	case scanner.MINUS:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l - r
	case scanner.PLUS:
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

	case scanner.SLASH:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l / r

	case scanner.STAR:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l * r

	case scanner.GREATER:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l > r

	case scanner.GREATER_EQUAL:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l >= r

	case scanner.LESS:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l < r

	case scanner.LESS_EQUAL:
		l, l_ok := left.(float64)
		r, r_ok := right.(float64)
		if !(l_ok && r_ok) {
			v.err = newOperandsError(e.Operator)
		}
		v.val = l <= r

	case scanner.BANG_EQUAL:
		v.val = !v.isEqual(left, right)

	case scanner.EQUAL_EQUAL:
		v.val = v.isEqual(left, right)
	}
}

func (v *Interpreter) VisitCall(e expression.Call) {
	var args []any
	callee, err := v.Evaluate(e.Callee)
	if err != nil {
		v.err = err
		return
	}

	for _, arg := range e.Args {
		val, err := v.Evaluate(arg)
		if err != nil {
			v.err = err
			return
		}
		args = append(args, val)
	}

	// Call a loxcallable
	lox_func, ok := callee.(LoxCallable)
	if !ok {
		v.err = newRuntimeError(e.Paren, "Can only call functions and classes.")
		return
	}

	if len(args) != lox_func.Arity() {
		v.err = newRuntimeError(e.Paren, fmt.Sprint("Expected", lox_func.Arity(), "arguments but got", len(args)))
		return
	}

	v.val = lox_func.Call(*v, args)
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

func (v *Interpreter) VisitLogical(e expression.Logical) {
	left, err := v.Evaluate(e.Left)
	if err != nil {
		v.err = err
		return
	}

	left_truth_value := v.isTruthy(left)

	switch e.Operator.Token_type {
	case scanner.OR:
		if left_truth_value {
			v.val = left
			v.err = nil
			return
		} else {
			right, err := v.Evaluate(e.Right)
			if err != nil {
				v.err = err
				return
			}

			v.val = right
			v.err = nil
			return
		}
	case scanner.AND:
		if !left_truth_value {
			v.err = nil
			v.val = left
			return
		} else {
			right, err := v.Evaluate(e.Right)
			if err != nil {
				v.err = err
				return
			}
			v.val = right
			v.err = nil
			return
		}
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
	case scanner.MINUS:
		r, ok := right.(float64)
		if !ok {
			v.err = newNumberError(e.Operator)
		}
		v.val = -r
	case scanner.BANG:
		v.val = v.isTruthy(right)
	}
}

func (v *Interpreter) VisitVariable(e expression.Variable) {
	val, err := v.pEnvironment.Get(e.GetToken())
	if err != nil {
		v.err = newRuntimeError(e.GetToken(), err.Error())
		return
	}

	v.val = val
}

func (v *Interpreter) VisitBlockStmt(stmt statement.Block) {
	// Declare a new environment
	// Execute all the declarations in the block
	env := NewEnvironment()
	v.executeBlock(stmt.GetStatements(), env)
}
func (v *Interpreter) VisitClassStmt(stmt statement.Class) {
}
func (v *Interpreter) VisitExpressionStmt(stmt statement.Expression) {
	val, err := v.Evaluate(stmt.Val)
	if err == nil && v.interactiveMode {
		fmt.Println(val)
	}
}
func (v *Interpreter) VisitFunctionStmt(stmt statement.Function) {
    var funcDef UserCallable
    funcDef.declaration = stmt

    v.pEnvironment.Define(stmt.Name.Lexeme, funcDef)
}
func (v *Interpreter) VisitIfStmt(stmt statement.If) {
	val, err := v.Evaluate(stmt.Conditional)
	if err != nil {
		v.err = err
		return
	}
	if v.isTruthy(val) {
		err := v.execute(stmt.If_stmt)
		if err != nil {
			v.err = err
			return
		}
	} else {
		if stmt.Else_stmt != nil {
			err := v.execute(stmt.Else_stmt)
			if err != nil {
				v.err = err
				return
			}
		}
	}
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
	v.pEnvironment.Define(stmt.Name.Lexeme, val)
}

func (v *Interpreter) VisitWhileStmt(stmt statement.While) {
	var err *RuntimeError
	var val any
	for val, err = v.Evaluate(stmt.Conditional); err == nil && v.isTruthy(val); val, err = v.Evaluate(stmt.Conditional) {
		err = v.execute(stmt.Stmt)
		if err != nil {
			v.err = err
			return
		}
	}
	v.err = err
}
