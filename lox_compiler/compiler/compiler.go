package compiler

import (
	"fmt"
	"lox-compiler/bytecode"
	"lox-compiler/debug"
	"lox-compiler/parser"
)

type CompilationError struct {
	err string
}

func (e CompilationError) Error() string {
	return fmt.Sprintf("a compilation error ocurred: %s", e.err)
}

type Compiler struct {
	rootChunk       *bytecode.Chunk
	curChunk        *bytecode.Chunk
	InteractiveMode bool
}

func (c *Compiler) Compile(source string) (*bytecode.Chunk, *CompilationError) {
	s := parser.NewScanner(source)
	tokens, err := s.ScanTokens()
	if err != nil {
		return nil, &CompilationError{err: err.Error()}
	}
	p := parser.NewParser(tokens)
	ast := p.Parse()
	debug.Printf("%v", tokens)
	debug.Printf("%s", ast)
	compilationErr := c.compileFromAST(ast)
	debug.Printf("%s", *c.rootChunk)
	// c.rootChunk.AddInst(bytecode.NewReturnInst(1))
	return c.rootChunk, compilationErr
}

func (c *Compiler) compileFromAST(nodes []parser.ASTNode) *CompilationError {
	chunk := bytecode.NewChunk()
	c.curChunk = &chunk
	c.rootChunk = c.curChunk
	for _, stmt := range nodes {
		err := c.compileStmt(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) compileStmt(stmt parser.Statement) *CompilationError {
	switch v := stmt.(type) {
	case parser.Block:
		return c.compileBlock(v)
	case parser.Class:
		return c.compileClass(v)
	case parser.ExpressionStmt:
		return c.compileExpressionStmt(v)
	case parser.Function:
		return c.compileFunction(v)
	case parser.If:
		return c.compileIf(v)
	case parser.Print:
		return c.compilePrint(v)
	case parser.Return:
		return c.compileReturn(v)
	case parser.Var:
		return c.compileVar(v)
	case parser.While:
		return c.compileWhile(v)
	}

	return &CompilationError{err: "expected a statement"}
}

func (c *Compiler) compileExpr(e parser.Expr) *CompilationError {
	switch v := e.(type) {
	case parser.Assign:
		return c.compileAssign(v)
	case parser.Binary:
		return c.compileBinary(v)
	case parser.Call:
		return c.compileCall(v)
	case parser.Get:
		return c.compileGet(v)
	case parser.Grouping:
		return c.compileGrouping(v)
	case parser.Literal:
		return c.compileLiteral(v)
	case parser.Logical:
		return c.compileLogical(v)
	case parser.Unary:
		return c.compileUnary(v)
	case parser.Set:
		return c.compileSet(v)
	case parser.Super:
		return c.compileSuper(v)
	case parser.This:
		return c.compileThis(v)
	case parser.Variable:
		return c.compileVariable(v)
	}

	return &CompilationError{err: "expected an expression"}
}

func (c *Compiler) compileBlock(stmt parser.Block) *CompilationError {
	return &CompilationError{err: "compiling `Block` statements is not implemented"}
}

func (c *Compiler) compileClass(stmt parser.Class) *CompilationError {
	return &CompilationError{err: "compiling `Class` statements is not implemented"}
}

func (c *Compiler) compileExpressionStmt(stmt parser.ExpressionStmt) *CompilationError {
	err := c.compileExpr(stmt.Val)
	if err != nil || !c.InteractiveMode {
		return err
	}
	switch stmt.Val.(type) {
	case parser.Assign:
		return nil
	default:
		printInst := bytecode.NewPrintInst(0)
		c.curChunk.AddInst(printInst)
	}
	return nil
}

func (c *Compiler) compileFunction(stmt parser.Function) *CompilationError {
	return &CompilationError{err: "compiling `Function` statements is not implemented"}
}

func (c *Compiler) compileIf(stmt parser.If) *CompilationError {
	return &CompilationError{err: "compiling `If` statements is not implemented"}
}

func (c *Compiler) compilePrint(stmt parser.Print) *CompilationError {
	c.compileExpr(stmt.Val)
	c.curChunk.AddInst(bytecode.NewPrintInst(0))

	return nil
}

func (c *Compiler) compileReturn(stmt parser.Return) *CompilationError {
	return &CompilationError{err: "compiling `Return` statements is not implemented"}
}

func (c *Compiler) compileVar(stmt parser.Var) *CompilationError {
	// declare Variable
	constIndex := c.curChunk.AddConstant(
		bytecode.LoxString(stmt.Name.Lexeme),
	)
	c.curChunk.AddInst(
		bytecode.NewConstantInst(
			bytecode.Operand(constIndex),
			stmt.Name.Line,
		),
	)
	c.curChunk.AddInst(
		bytecode.Instruction{Code: bytecode.OpDeclare, SourceLineNumer: stmt.Name.Line},
	)

	// if there's an Initializer
	// evaluate Initializer
	if stmt.Initializer != nil {
		err := c.compileExpr(stmt.Initializer)
		if err != nil {
			return err
		}
		// assign var to Initializer
		c.curChunk.AddInst(
			bytecode.NewConstantInst(
				bytecode.Operand(constIndex),
				stmt.Name.Line,
			),
		)
		c.curChunk.AddInst(
			bytecode.Instruction{Code: bytecode.OpAssign, SourceLineNumer: stmt.Name.Line},
		)
	}
	return nil
}

func (c *Compiler) compileWhile(stmt parser.While) *CompilationError {
	return &CompilationError{err: "compiling `While` statements is not implemented"}
}

func (c *Compiler) compileAssign(e parser.Assign) *CompilationError {
	// return &CompilationError{err: "compiling Assign expression is not implemented"}
	// expression
	err := c.compileExpr(e.Value)
	if err != nil {
		return err
	}
	// store var Name
	c.curChunk.AddInst(
		bytecode.NewConstantInst(
			bytecode.Operand(c.curChunk.AddConstant(
				bytecode.LoxString(e.Name.Lexeme),
			)),
			e.Name.Line,
		),
	)
	// write assign
	c.curChunk.AddInst(bytecode.Instruction{Code: bytecode.OpAssign, SourceLineNumer: e.Name.Line})

	return nil
}

func (c *Compiler) compileBinary(e parser.Binary) *CompilationError {
	token_op_mapping := map[parser.TokenType]bytecode.OpCode{
		parser.OR:            bytecode.OpOr,
		parser.AND:           bytecode.OpAnd,
		parser.LESS:          bytecode.OpLess,
		parser.LESS_EQUAL:    bytecode.OpLessEqual,
		parser.GREATER:       bytecode.OpGreater,
		parser.GREATER_EQUAL: bytecode.OpGreaterEqual,
		parser.EQUAL_EQUAL:   bytecode.OpEqualEqual,
		parser.BANG_EQUAL:    bytecode.OpNotEqual,
		parser.PLUS:          bytecode.OpAdd,
		parser.MINUS:         bytecode.OpSubtract,
		parser.STAR:          bytecode.OpMultiply,
		parser.SLASH:         bytecode.OpDivide,
	}
	err := c.compileExpr(e.Left)
	if err != nil {
		return err
	}

	err = c.compileExpr(e.Right)
	if err != nil {
		return err
	}

	// switch e.Operator.Token_type {
	// case parser.PLUS:
	// 	c.curChunk.AddInst(bytecode.NewAddInst(e.Operator.Line))
	// case parser.MINUS:
	// 	c.curChunk.AddInst(bytecode.NewSubtractInst(e.Operator.Line))
	// case parser.STAR:
	// 	c.curChunk.AddInst(bytecode.NewMultiplyInst(e.Operator.Line))
	// case parser.SLASH:
	// 	c.curChunk.AddInst(bytecode.NewDivideInst(e.Operator.Line))
	// default:
	// 	return &CompilationError{err: fmt.Sprintf("unknown token %s", e.Operator.Token_type)}
	// }
	c.curChunk.AddInst(
		bytecode.Instruction{
			Code:            token_op_mapping[e.Operator.Token_type],
			SourceLineNumer: e.Operator.Line,
		},
	)

	return nil
}

func (c *Compiler) compileCall(e parser.Call) *CompilationError {
	return &CompilationError{err: "compiling Call expression is not implemented"}
}

func (c *Compiler) compileGet(e parser.Get) *CompilationError {
	return &CompilationError{err: "compiling Get expression is not implemented"}
}

func (c *Compiler) compileGrouping(e parser.Grouping) *CompilationError {
	return c.compileExpr(e.Expr)
}

func (c *Compiler) compileLiteral(e parser.Literal) *CompilationError {
	v, err := bytecode.NewValue(e.Value)
	if err != nil {
		return &CompilationError{err: err.Error()}
	}

	i := c.curChunk.AddConstant(v)
	c.curChunk.AddInst(bytecode.NewConstantInst(bytecode.Operand(i), 0))
	return nil
}

func (c *Compiler) compileLogical(e parser.Logical) *CompilationError {
	token_op_mapping := map[parser.TokenType]bytecode.OpCode{
		parser.OR:            bytecode.OpOr,
		parser.AND:           bytecode.OpAnd,
		parser.LESS:          bytecode.OpLess,
		parser.LESS_EQUAL:    bytecode.OpLessEqual,
		parser.GREATER:       bytecode.OpGreater,
		parser.GREATER_EQUAL: bytecode.OpGreaterEqual,
		parser.EQUAL_EQUAL:   bytecode.OpEqualEqual,
		parser.BANG_EQUAL:    bytecode.OpNotEqual,
	}
	c.compileExpr(e.Left)
	c.compileExpr(e.Right)
	c.curChunk.AddInst(
		bytecode.Instruction{
			Code:            token_op_mapping[e.Operator.Token_type],
			SourceLineNumer: e.Operator.Line,
		},
	)

	return nil
}

func (c *Compiler) compileUnary(e parser.Unary) *CompilationError {
	err := c.compileExpr(e.Right)
	if err != nil {
		return err
	}
	switch e.Operator.Token_type {
	case parser.MINUS, parser.BANG:
		c.curChunk.AddInst(bytecode.NewNegateInst(e.Operator.Line))

	default:
		return &CompilationError{err: fmt.Sprintf("expected a unary operator but got %s", e.Operator.Lexeme)}
	}

	return nil
}

func (c *Compiler) compileSet(e parser.Set) *CompilationError {
	return &CompilationError{err: "compiling Set expression is not implemented"}
}

func (c *Compiler) compileSuper(e parser.Super) *CompilationError {
	return &CompilationError{err: "compiling Super expression is not implemented"}
}

func (c *Compiler) compileThis(e parser.This) *CompilationError {
	return &CompilationError{err: "compiling This expression is not implemented"}
}

func (c *Compiler) compileVariable(e parser.Variable) *CompilationError {
	// return &CompilationError{err: "compiling Variable expression is not implemented"}
	c.curChunk.AddInst(
		bytecode.NewConstantInst(
			bytecode.Operand(c.curChunk.AddConstant(
				bytecode.LoxString(e.Name.Lexeme),
			)),
			e.Name.Line,
		),
	)
	c.curChunk.AddInst(
		bytecode.Instruction{Code: bytecode.OpLookup, SourceLineNumer: e.Name.Line},
	)
	return nil
}
