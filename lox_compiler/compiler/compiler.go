package compiler

import (
	"fmt"
	"lox-compiler/bytecode"
	"lox-compiler/debug"
	"lox-compiler/parser"
	"math"
)

const maxLocals int = math.MaxUint8

type CompilationError struct {
	err string
}

func (e CompilationError) Error() string {
	return fmt.Sprintf("compilation error: %s", e.err)
}

type Compiler struct {
	rootChunk       *bytecode.Chunk
	curChunk        *bytecode.Chunk
	InteractiveMode bool
	localCount      int
	scopeDepth      int
	locals          [maxLocals]local
}

type local struct {
	name  parser.Token
	depth int
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
	c.beginScope()
	defer c.endScope()
	for _, s := range stmt.Statements {
		if err := c.compileStmt(s); err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) beginScope() {
	c.scopeDepth++
}

func (c *Compiler) endScope() {
	c.scopeDepth--
	for c.localCount > 0 && (c.locals[c.localCount-1].depth > c.scopeDepth) {
		// pop
		c.curChunk.AddInst(bytecode.NewInst(bytecode.OpPop, -1))
		c.localCount--
	}
}

func (c *Compiler) compileClass(stmt parser.Class) *CompilationError {
	return &CompilationError{err: "compiling `Class` statements is not implemented"}
}

func (c *Compiler) compileExpressionStmt(stmt parser.ExpressionStmt) *CompilationError {
	err := c.compileExpr(stmt.Val)
	if err != nil {
		return err
	}

	if c.InteractiveMode {
		printInst := bytecode.NewPrintInst(0)
		c.curChunk.AddInst(printInst)
	} else {
		c.curChunk.AddInst(bytecode.NewInst(bytecode.OpPop, 0))
	}

	return nil
}

func (c *Compiler) compileFunction(stmt parser.Function) *CompilationError {
	return &CompilationError{err: "compiling `Function` statements is not implemented"}
}

func (c *Compiler) compileIf(stmt parser.If) *CompilationError {
	// we need to add a two operand instruction here. The first holds the const
	// index of the "true" jump and the second the index of the "false" jump
	var curLen, falseJmpOffsetIndex, falseBlockSizeIndex int
	// falseBlockSizeIndex = c.curChunk.AddConstant(bytecode.LoxInt(0))
	if err := c.compileExpr(stmt.Conditional); err != nil {
		return err
	}
    _, falseJmpOffsetIndex = c.addConditionalJmp()
	curLen = len(c.curChunk.InstructionSlice)
	if err := c.compileStmt(stmt.If_stmt); err != nil {
		return err
	}
	// Skip the "else" statement rather than falling through into it
    falseBlockSizeIndex = c.addJmp()

	// backpatch the offsets
	// c.curChunk.Constants[trueJmpOffsetIndex] = bytecode.LoxInt(0) // The vm increments the pc all on its own.
	c.curChunk.Constants[falseJmpOffsetIndex] = bytecode.LoxInt(len(c.curChunk.InstructionSlice) - curLen)
	if stmt.Else_stmt == nil {
		return nil
	}
	curLen = len(c.curChunk.InstructionSlice)
	if err := c.compileStmt(stmt.Else_stmt); err != nil {
		return err
	}
	// backpatch the "else" jump
	c.curChunk.Constants[falseBlockSizeIndex] = bytecode.LoxInt(len(c.curChunk.InstructionSlice) - curLen)

	return nil
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
	var err *CompilationError
	if c.scopeDepth > 0 {
		err = c.compileLocalVar(stmt)
	} else {
		err = c.compileGlobalVar(stmt)
	}
	if err != nil {
		return err
	}
	if c.scopeDepth == 0 {
		c.curChunk.AddInst(bytecode.NewInst(bytecode.OpPop, stmt.Name.Line))
	}
	return nil
}

func (c *Compiler) compileLocalVar(stmt parser.Var) *CompilationError {
	var local *local
	for i := c.localCount - 1; i >= 0; i-- {
		local = &c.locals[i]
		if local.depth != -1 && local.depth < c.scopeDepth {
			break
		}
		if stmt.Name.Lexeme == local.name.Lexeme {
			return &CompilationError{err: "already a variable with this name in this scope"}
		}
	}
	// The result of this operation becomes the top of the stack, then
	// that index of the stack becomes a local variable
	err := c.compileExpr(stmt.Initializer)
	if err != nil {
		return err
	}
	return c.addLocal(stmt.Name)
}

func (c *Compiler) compileGlobalVar(stmt parser.Var) *CompilationError {
	// declare a global Variable
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
		bytecode.Instruction{Code: bytecode.OpDeclareGlobal, SourceLineNumer: stmt.Name.Line},
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

func (c *Compiler) addLocal(name parser.Token) *CompilationError {
	if c.localCount > maxLocals-1 {
		return &CompilationError{err: "too many local variables declared"}
	}
	var local *local = &c.locals[c.localCount]
	c.localCount++
	local.name = name
	local.depth = c.scopeDepth
	return nil
}

func (c *Compiler) compileWhile(stmt parser.While) *CompilationError {
	return &CompilationError{err: "compiling `While` statements is not implemented"}
}

func (c *Compiler) compileAssign(e parser.Assign) *CompilationError {
	err := c.compileExpr(e.Value)
	if err != nil {
		return err
	}
	if l, i := c.getLocalVar(e.Name); l != nil {
		c.curChunk.AddInst(
			bytecode.Instruction{
				Code:            bytecode.OpLocalAssign,
				Operands:        bytecode.OperandArray{bytecode.Operand(i)},
				SourceLineNumer: e.Name.Line,
			},
		)
		return nil
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
	l, i := c.getLocalVar(e.Name)
	if l != nil {
		return c.compileLocalLookup(i)
	}
	return c.compileGlobalLookup(e)
}

func (c *Compiler) getLocalVar(name parser.Token) (*local, int) {
	for i := c.localCount; i >= 0; i-- {
		if c.locals[i].name.Lexeme == name.Lexeme {
			return &c.locals[i], i
		}
	}

	return nil, -1
}

func (c *Compiler) compileLocalLookup(index int) *CompilationError {
	c.curChunk.AddInst(
		bytecode.Instruction{
			Code:     bytecode.OpLocalLookup,
			Operands: bytecode.OperandArray{bytecode.Operand(index)},
		},
	)

	return nil
}

func (c *Compiler) compileGlobalLookup(e parser.Variable) *CompilationError {
	c.curChunk.AddInst(
		bytecode.NewConstantInst(
			bytecode.Operand(c.curChunk.AddConstant(
				bytecode.LoxString(e.Name.Lexeme),
			)),
			e.Name.Line,
		),
	)
	c.curChunk.AddInst(
		bytecode.Instruction{Code: bytecode.OpGlobalLookup, SourceLineNumer: e.Name.Line},
	)
	return nil
}

// Define a new conditional jump instruction and return the
// indicies that we will store the two offsets.
func (c *Compiler) addConditionalJmp() (trueJmpIndex, falseJmpIndex int) {
	trueJmpIndex = c.curChunk.AddConstant(bytecode.LoxInt(0))
	falseJmpIndex = c.curChunk.AddConstant(bytecode.LoxInt(0))

	c.curChunk.AddInst(
		bytecode.Instruction{
			Code: bytecode.OpConditionalJump,
			Operands: bytecode.OperandArray{
				bytecode.Operand(trueJmpIndex),
				bytecode.Operand(falseJmpIndex),
			},
		},
	)

    return trueJmpIndex, falseJmpIndex
}

// Define a new jump instruction and return the
// index that we will store the offset.
func (c *Compiler) addJmp() (jmpIndex int) {
	jmpIndex = c.curChunk.AddConstant(bytecode.LoxInt(0))

	c.curChunk.AddInst(
		bytecode.Instruction{
			Code: bytecode.OpJump,
			Operands: bytecode.OperandArray{
				bytecode.Operand(jmpIndex),
			},
		},
	)

    return jmpIndex
}
