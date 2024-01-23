package interpreter

import (
	"fmt"
	"golox/expression"
	"golox/scanner"
	"golox/statement"
)

type scope = map[string]bool
type functionType int
type classType int

const (
	notFunction functionType = iota
	function
	method
	initializer
)

const (
	notClass classType = iota
	class
)

type Resolver struct {
	Interp          *Interpreter
	scopes          []scope
	err             error
	currentFunction functionType
	currentClass    classType
}

type resolver_error struct {
	prefix string
	msg    string
}

func (e resolver_error) Error() string {
	return fmt.Sprintf("error in '%s': %s", e.prefix, e.msg)
}

func (r *Resolver) Resolve(stmts []statement.Statement) error {
	return r.resolve_statements(stmts)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(scope))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.Lexeme] = false
}

func (r *Resolver) define(name scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}

	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) resolve_statements(statements []statement.Statement) error {
	for _, s := range statements {
		r.resolve_statement(s)
		if r.err != nil {
			return r.err
		}
	}

	return nil
}

func (r *Resolver) resolve_statement(stmt statement.Statement) error {
	stmt.Accept(r)
	return r.err
}

func (r *Resolver) resolve_expression(e expression.Expr) error {
	e.Accept(r)
	return r.err
}

func (r *Resolver) resolveLocal(e expression.Expr, name scanner.Token) error {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		_, ok := r.scopes[i][name.Lexeme]
		if ok {
			r.Interp.resolve(e, len(r.scopes)-1-i)
			return nil
		}
	}
	return nil
}

func (r *Resolver) resolveFunction(stmt statement.Function, t functionType) error {
	r.beginScope()
	defer r.endScope()
	defer r.setFunctionStatus(r.currentFunction)
	r.setFunctionStatus(t)
	for _, param := range stmt.Params {
		r.declare(param)
		r.define(param)
	}

	err := r.resolve_statements(stmt.Body)
	if err != nil {
		return err
	}

	return nil
}

func (r *Resolver) setFunctionStatus(status functionType) {
	r.currentFunction = status
}

func (r *Resolver) VisitAssign(e expression.Assign) {
	r.err = r.resolve_expression(e.Value)
	if r.err != nil {
		return
	}

	r.err = r.resolveLocal(e, e.Name)
}
func (r *Resolver) VisitBinary(e expression.Binary) {
	r.err = r.resolve_expression(e.Left)
	if r.err != nil {
		return
	}

	r.err = r.resolve_expression(e.Right)
}
func (r *Resolver) VisitCall(e expression.Call) {
	r.err = r.resolve_expression(e.Callee)
	if r.err != nil {
		return
	}

	for _, arg := range e.Args {
		r.err = r.resolve_expression(arg)
		if r.err != nil {
			return
		}
	}
}

func (r *Resolver) VisitGet(e expression.Get) {
	r.err = r.resolve_expression(e.Object)
}

func (r *Resolver) VisitGrouping(e expression.Grouping) {
	r.err = r.resolve_expression(e.Expr)
}
func (r *Resolver) VisitLiteral(e expression.Literal) {
	return
}
func (r *Resolver) VisitLogical(e expression.Logical) {
	r.err = r.resolve_expression(e.Left)
	if r.err != nil {
		return
	}

	r.err = r.resolve_expression(e.Right)
}

func (r *Resolver) VisitThis(e expression.This) {
	if r.currentClass != class {
		r.err = resolver_error{prefix: fmt.Sprintf("[line %d]:", e.Keyword.Line), msg: "\"this\" can't be used outside of a class declaration"}
		return
	}
	r.resolveLocal(e, e.Keyword)
}

func (r *Resolver) VisitSet(e expression.Set) {
	r.resolve_expression(e.Object)
	r.resolve_expression(e.Value)
}

func (r *Resolver) VisitUnary(e expression.Unary) {
	r.err = r.resolve_expression(e.Right)
}
func (r *Resolver) VisitVariable(e expression.Variable) {
	if len(r.scopes) > 0 {
		res, ok := r.scopes[len(r.scopes)-1][e.GetName()]
		if ok && res == false {
			r.err = resolver_error{prefix: e.GetName(), msg: "can't read a local variable in its own initializer"}
			return
		}
	}

	r.resolveLocal(e, e.GetToken())

}
func (r *Resolver) VisitBlockStmt(stmt statement.Block) {
	r.beginScope()
	defer r.endScope()
	r.resolve_statements(stmt.GetStatements())
}

func (r *Resolver) changeClassType(t classType) {
	r.currentClass = t
}

func (r *Resolver) VisitClassStmt(stmt statement.Class) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.beginScope()
	defer r.endScope()
	defer r.changeClassType(r.currentClass)
	r.changeClassType(class)
	r.scopes[len(r.scopes)-1]["this"] = true

	for _, m := range stmt.Methods {
		if m.Name.Lexeme == constructor_name {
			r.resolveFunction(m, initializer)
		} else {
			r.resolveFunction(m, method)
        }
	}
}
func (r *Resolver) VisitExpressionStmt(stmt statement.Expression) {
	r.err = r.resolve_expression(stmt.Val)
}
func (r *Resolver) VisitFunctionStmt(stmt statement.Function) {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.err = r.resolveFunction(stmt, function)
}
func (r *Resolver) VisitIfStmt(stmt statement.If) {
	r.err = r.resolve_expression(stmt.Conditional)
	if r.err != nil {
		return
	}
	r.err = r.resolve_statement(stmt.If_stmt)
	if r.err != nil {
		return
	}
	if stmt.Else_stmt != nil {
		r.err = r.resolve_statement(stmt.Else_stmt)
	}
}
func (r *Resolver) VisitPrintStmt(stmt statement.Print) {
	r.err = r.resolve_expression(stmt.Val)
}
func (r *Resolver) VisitReturnStmt(stmt statement.Return) {
	if r.currentFunction == notFunction {
		r.err = resolver_error{prefix: "return statement", msg: "cannot call \"return\" outside of a function or method"}
		return
	} else if r.currentFunction == initializer {
		r.err = resolver_error{prefix: "return statement", msg: "cannot call \"return\" inside of an initializer"}
        return
    }
	r.err = r.resolve_expression(stmt.Return_expr)
}
func (r *Resolver) VisitVarStmt(stmt statement.Var) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolve_expression(stmt.Initializer)
	}

	r.define(stmt.Name)
}
func (r *Resolver) VisitWhileStmt(stmt statement.While) {
	r.err = r.resolve_expression(stmt.Conditional)
	if r.err != nil {
		return
	}

	r.err = r.resolve_statement(stmt.Stmt)
}
