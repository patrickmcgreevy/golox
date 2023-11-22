package statement

type statementVistitor interface {
	visitBlockStmt(stmt Block)
	visitClassStmt(stmt Class)
	visitExpressionStmt(stmt Expression)
	visitFunctionStmt(stmt Function)
	visitIfStmt(stmt If)
	visitPrintStmt(stmt Print)
	visitReturnStmt(stmt Return)
	visitVarStmt(stmt Var)
	visitWhileStmt(stmt While)
}

type Statement interface {
    accept(statementVistitor)
}

type Block struct {
}

type Class struct {
}

type Expression struct {
}

type Function struct {
}

type If struct {
}

type Print struct {
}

type Return struct {
}

type Var struct {
}

type While struct {
}

