package statement

import (
    "golox/token"
	"golox/expression"
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
}

type Class struct {
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
}

type If struct {
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
}

type Var struct {
    Initializer expression.Expr
    Name token.Token
}

func NewVarStmt( name token.Token, initializer expression.Expr) Var {
    return Var{Initializer: nil, Name: name}
}

func (s Var) Accept(v StatementVisitor) {
    v.VisitVarStmt(s)
}



type While struct {
}

