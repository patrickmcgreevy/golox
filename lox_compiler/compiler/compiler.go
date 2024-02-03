package compiler

import (
	"fmt"
	"lox-compiler/bytecode"
	"lox-compiler/scanner"
    "lox-compiler/debug"
)

type CompilationError struct {
    err string
}

func (e CompilationError) Error() string {
    return fmt.Sprintf("a compilation error ocurred: %s", e.err)
}

func Compile(source string) (*bytecode.Chunk, *CompilationError) {
    s := scanner.NewScanner(source)
    tokens, err := s.ScanTokens()
    if err != nil {
        return nil, &CompilationError{err: err.Error()}
    }
    debug.Printf("%v", tokens)
    c := bytecode.NewChunk()
    c.AddInst(bytecode.NewConstantInst(bytecode.Operand(c.AddConstant(32)), 0))
    c.AddInst(bytecode.NewReturnInst(0))
    return &c, nil
}
