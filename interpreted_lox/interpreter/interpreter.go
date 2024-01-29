package interpreter

import (
	"fmt"
    "bufio"
    "os"
	"golox/errorhandling"
	"golox/expression"
	"golox/scanner"
	"golox/statement"
	"reflect"
	"time"
)

type LoxCallable interface {
	Call(interp Interpreter, args []any) (any, *RuntimeError)
	Arity() int
}

type RuntimeError struct {
	error        string
	tok          scanner.Token
	return_value any // This is used to return values up the call stack
}

func (e RuntimeError) Error() string {
    if e.tok.Lexeme != "" {
        return fmt.Sprintf("[line %d]: %s", e.tok.Line, e.error)
    }
	return e.error
}

func (e RuntimeError) GetToken() scanner.Token {
	return e.tok
}

// This is used to return a value up the call stack to the 'call' function
func newReturnError(val any) *RuntimeError {
	return &RuntimeError{error: "'return' statement outside of function", return_value: val}
}

func newRuntimeError(operator scanner.Token, message string) *RuntimeError {
	// msg := fmt.Sprintf("[line %d]: %s", operator.Line, message)
	new_err := RuntimeError{error: message, tok: operator}
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
	locals          map[expression.Expr]int
	globals         Environment
}

func NewInterpreter() Interpreter {
	globals := NewEnvironment()
	env := globals

	globals.Define("clock", BuiltinCallable{params: make([]string, 0), foo: func(a Interpreter, b []any) any {
		return float64(time.Now().UnixMilli()) / 1000
	}})
    globals.Define("input", BuiltinCallable{params: make([]string, 0), foo: func(a Interpreter, b []any) any {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("> ")
        line, err := reader.ReadString('\n')
        if err != nil {
            if err.Error() == "EOF" {
                return nil
            }
            panic(err)
        }
        return line
    }})
	return Interpreter{val: nil, err: nil, pEnvironment: &env, interactiveMode: false, locals: make(map[expression.Expr]int), globals: globals}
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

func (v *Interpreter) resolve(e expression.Expr, depth int) {
	v.locals[e] = depth
}

func (v *Interpreter) execute(stmt statement.Statement) *RuntimeError {
	stmt.Accept(v)
	if v.err != nil {
		return v.err
	}
	return nil
}

func (v *Interpreter) executeBlock(statements []statement.Statement, env Environment) {
	// I think that by calling pushEnvironment, I blow up the closure
	v.pushEnvironment(&env)
	defer v.popEnvironment()
	for _, stmt := range statements {
		v.execute(stmt)
	}
}

func (v *Interpreter) pushEnvironment(env *Environment) {
	// Often, we are passed an env that doesn't have an enclosing scope
	// in that case, we need to supply one.
	if env.enclosing == nil {
		env.SetEnclosing(v.pEnvironment)
	}
	v.pEnvironment = env
}

func (v *Interpreter) popEnvironment() {
	parent := v.pEnvironment.GetEnclosing()
	if parent != nil {
		v.pEnvironment = parent
	}
}

func (v *Interpreter) Evaluate(e expression.Expr) (any, *RuntimeError) {
	if e == nil {
		v.val, v.err = nil, nil
		return nil, nil
	}
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
	val, ok := v.locals[e]
	if !ok {
		// v.err = newRuntimeError(e.Name, "undefined variable")
		// return
		v.globals.Assign(e.Name.Lexeme, right)
	}
	assignment_error := v.pEnvironment.AssignAt(val, e.Name.Lexeme, right)
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
		v.err = newRuntimeError(e.Paren, fmt.Sprint("Expected ", lox_func.Arity(), " arguments but got ", len(args)))
		return
	}

	val, err := lox_func.Call(*v, args)
	if err != nil {
		if err.return_value != nil {
			val, err = err.return_value, nil
		}
	}

	v.val, v.err = val, err
}

func (v *Interpreter) VisitGet(e expression.Get) {
	val, err := v.Evaluate(e.Object)
	if err != nil {
		v.err = err
		return
	}

	obj, ok := val.(LoxInstance)
	if !ok {
		v.err = &RuntimeError{error: "only class instances have properties", tok: e.Name}
		return
	}

	ret, err := obj.Get(e.Name)
	if err != nil {
		v.err = err
		return
	}

	v.val = ret
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

func (v *Interpreter) VisitThis(e expression.This) {
    val, err := v.lookUpVariable(e.Keyword, e)
    if err != nil {
        newError := newRuntimeError(e.Keyword, err.Error())
        v.err = newError
        return
    }

    v.val, v.err = val, nil
}

func (v *Interpreter) VisitSet(e expression.Set) {
	obj, err := v.Evaluate(e.Object)
	if err != nil {
		v.err = err
		return
	}
	instance, ok := obj.(LoxInstance)
	if !ok {
		v.err = &RuntimeError{error: "only instances have fields", tok: e.Name}
		return
	}
	val, err := v.Evaluate(e.Value)
	if err != nil {
		v.err = err
		return
	}
	instance.Fields[e.Name.Lexeme] = val

	v.val = val
}

func (v *Interpreter) VisitSuper(e expression.Super) {
    dist := v.locals[e]
    val, err := v.pEnvironment.GetAt(dist, "super")
    if err != nil {
        v.err = &RuntimeError{error: err.Error(), tok: e.Keyword}
        return
    }

    super, ok := val.(LoxClass)
    if !ok {
        v.err = &RuntimeError{error: "super did not evaluate to a class", tok: e.Keyword}
        return
    }

    val, err = v.pEnvironment.GetAt(dist-1, "this")
    if err != nil {
        v.err = &RuntimeError{error: err.Error(), tok: e.Keyword}
        return
    }
    instance, ok := val.(LoxInstance)
    if !ok {
        v.err = &RuntimeError{error: "this did not evaluate to an object", tok: e.Keyword}
        return
    }

    method, err := super.GetMethod(e.Method.Lexeme)
    if err != nil {
        v.err = &RuntimeError{error: err.Error(), tok: e.Method}
        return
    }
    v.val = method.Bind(instance)
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
	val, err := v.lookUpVariable(e.GetToken(), e)
	if err != nil {
		v.err = newRuntimeError(e.GetToken(), err.Error())
		return
	}

	v.val = val
}

func (v *Interpreter) lookUpVariable(name scanner.Token, expr expression.Expr) (any, error) {
	distance, ok := v.locals[expr]
	if ok {
		return v.pEnvironment.GetAt(distance, name.Lexeme)
	} else {
		return v.globals.Get(name.Lexeme)
	}
}

func (v *Interpreter) VisitBlockStmt(stmt statement.Block) {
	// Declare a new environment
	// Execute all the declarations in the block
	env := NewEnvironment()
	v.executeBlock(stmt.GetStatements(), env)
}
func (v *Interpreter) VisitClassStmt(stmt statement.Class) {
    var parentClass LoxClass
    var ok bool
    methods := make(map[string]UserCallable)

    if stmt.ParentClass != nil {
        parent, err := v.Evaluate(stmt.ParentClass)
        if err != nil {
            v.err = err
            return
        }

        parentClass, ok = parent.(LoxClass)
        if !ok {
            v.err = &RuntimeError{error: "superclass must be a class", tok: stmt.ParentClass.GetToken()}
            return
        }
    }
	v.pEnvironment.Define(stmt.Name.Lexeme, nil)
    if stmt.ParentClass != nil {
        env := NewEnvironment()
        v.pushEnvironment(&env)
        v.pEnvironment.Define("super", parentClass)
    }
    for _, m := range stmt.Methods {
        methods[m.Name.Lexeme] = UserCallable{declaration: m, closure: v.pEnvironment}
    }
    class := LoxClass{Name: stmt.Name.Lexeme, Methods: methods, Parent: &parentClass}
    if stmt.ParentClass != nil {
        v.popEnvironment()
    }
	v.pEnvironment.Assign(stmt.Name.Lexeme, class)
    
}
func (v *Interpreter) VisitExpressionStmt(stmt statement.Expression) {
	val, err := v.Evaluate(stmt.Val)
	if err == nil && v.interactiveMode {
		fmt.Println(val)
	}
}
func (v *Interpreter) VisitFunctionStmt(stmt statement.Function) {
	var funcDef UserCallable = UserCallable{declaration: stmt, closure: v.pEnvironment}

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
	val, err := v.Evaluate(stmt.Return_expr)
	if err == nil {
		err = newReturnError(val)
	}
	v.val, v.err = val, err
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
