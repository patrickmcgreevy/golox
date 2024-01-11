package interpreter

import "golox/expression"
import "golox/statement"
import "golox/scanner"

type scope = map[string]bool

type resolver struct {
	interp Interpreter
    scopes []scope
}


func (r *resolver) beginScope() {
    r.scopes = append(r.scopes, make(scope))
}

func (r *resolver) endScope() {
    r.scopes = r.scopes[len(r.scopes)-1:]
}

func (r *resolver) declare(name scanner.Token) {
    if len(r.scopes) == 0 {
        return
    }

    r.scopes[len(r.scopes)-1][name.Lexeme] = false
}

func (r *resolver) define(name scanner.Token) {
    if len(r.scopes) == 0 {
        return
    }

    r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *resolver) resolve_statements(statements []statement.Statement) {
	for _, s := range statements {
        r.resolve_statement(s)
	}
}

func (r *resolver) resolve_statement(stmt statement.Statement) {
	stmt.Accept(r)
}

func (r *resolver) resolve_expression(e expression.Expr) {
	e.Accept(r)
}

func (r *resolver) VisitAssign(e expression.Assign) {
}
func (r *resolver) VisitBinary(e expression.Binary) {
}
func (r *resolver) VisitCall(e expression.Call) {
}
func (r *resolver) VisitGrouping(e expression.Grouping) {
}
func (r *resolver) VisitLiteral(e expression.Literal) {
}
func (r *resolver) VisitLogical(e expression.Logical) {
}
func (r *resolver) VisitUnary(e expression.Unary) {
}
func (r *resolver) VisitVariable(e expression.Variable) {
}
func (r *resolver) VisitBlockStmt(stmt statement.Block) {
	r.beginScope()
	r.resolve_statements(stmt.GetStatements())
	r.endScope()
}
func (r *resolver) VisitClassStmt(stmt statement.Class) {
}
func (r *resolver) VisitExpressionStmt(stmt statement.Expression) {
}
func (r *resolver) VisitFunctionStmt(stmt statement.Function) {
}
func (r *resolver) VisitIfStmt(stmt statement.If) {
}
func (r *resolver) VisitPrintStmt(stmt statement.Print) {
}
func (r *resolver) VisitReturnStmt(stmt statement.Return) {
}
func (r *resolver) VisitVarStmt(stmt statement.Var) {
    r.declare(stmt.Name)
    if stmt.Initializer != nil {
        r.resolve_expression(stmt.Initializer)
    }

    r.define(stmt.Name)
}
func (r *resolver) VisitWhileStmt(stmt statement.While) {
}
