package compiler_test

import (
	"lox-compiler/compiler"
	"testing"
)

func TestExprStmt(t *testing.T) {
    c := compiler.Compiler{}
    chunk, err := c.Compile("2+2;")
    if err != nil {
        t.Fatalf("%s", err.Error())
    }

    chunk.Disassemble("main")
}

func TestCompileGrouping(t *testing.T) {
    c := compiler.Compiler{}
    chunk, err := c.Compile("(2+2);")
    if err != nil {
        t.Fatalf("%s", err.Error())
    }

    chunk.Disassemble("main")
}

func TestCompileVar(t *testing.T) {
    c := compiler.Compiler{}
    chunk, err := c.Compile("var a = 1+3;")
    if err != nil {
        t.Fatalf("%s", err.Error())
    }

    chunk.Disassemble("main")
}
