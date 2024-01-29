package statement

import (
	"golox/expression"
	"golox/scanner"
)

type StatementVisitor interface {
	VisitBlockStmt(stmt Block)
	VisitClassStmt(stmt Class)
	VisitExpressionStmt(stmt Expression)
	VisitFunctionStmt(stmt Function)
	VisitIfStmt(stmt If)
	VisitPrintStmt(stmt Print)
	VisitReturnStmt(stmt Return)
	VisitVarStmt(stmt Var)
	VisitWhileStmt(stmt While)
}

type Statement interface {
	Accept(StatementVisitor)
}

type Block struct {
	statements []Statement
}

func (s Block) Accept(v StatementVisitor) {
	v.VisitBlockStmt(s)
}

func NewBlockStmt(statments []Statement) Block {
	return Block{statements: statments}
}

func (s Block) GetStatements() []Statement {
	return s.statements
}

type Class struct {
    Name scanner.Token
    Methods []Function
    ParentClass *expression.Variable
}

func (s Class) Accept(v StatementVisitor) {
    v.VisitClassStmt(s)
}

type Expression struct {
	Val expression.Expr
}

func NewExpressionStmt(val expression.Expr) Expression {
	return Expression{Val: val}
}

func (s Expression) Accept(v StatementVisitor) {
	v.VisitExpressionStmt(s)
}

type Function struct {
    Name scanner.Token
    Params []scanner.Token
    Body []Statement
}

func (s Function) Accept(v StatementVisitor) {
    v.VisitFunctionStmt(s)
}

type If struct {
	Conditional expression.Expr
	If_stmt     Statement
	Else_stmt   Statement
}

func NewIfStatement(conitional expression.Expr, if_stmt, else_stmt Statement) If {
	return If{Conditional: conitional, If_stmt: if_stmt, Else_stmt: else_stmt}
}

func (s If) Accept(v StatementVisitor) {
	v.VisitIfStmt(s)
}

type Print struct {
	Val expression.Expr
}

func NewPrintStmt(val expression.Expr) Print {
	return Print{Val: val}
}

func (s Print) Accept(v StatementVisitor) {
	v.VisitPrintStmt(s)
}

type Return struct {
    Return_expr expression.Expr
}

func (s Return) Accept(v StatementVisitor) {
    v.VisitReturnStmt(s)
}

type Var struct {
	Initializer expression.Expr
	Name        scanner.Token
}

func NewVarStmt(name scanner.Token, initializer expression.Expr) Var {
	return Var{Initializer: initializer, Name: name}
}

func (s Var) Accept(v StatementVisitor) {
	v.VisitVarStmt(s)
}

type While struct {
	Conditional expression.Expr
	Stmt        Statement
}

func NewWhileStmt(conditional expression.Expr, stmt Statement) While {
	return While{Conditional: conditional, Stmt: stmt}
}

func (s While) Accept(v StatementVisitor) {
	v.VisitWhileStmt(s)
}
