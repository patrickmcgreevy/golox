package compiler_test

import (
	"lox-compiler/compiler"
	"testing"
)

func test_compilation(t *testing.T, s string) {
    c := compiler.Compiler{}
    chunk, err := c.Compile(s)
    if err != nil {
        t.Fatalf("%s", err.Error())
    }

    chunk.Disassemble("main")
}

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

func TestCompileVarNoInitializer( t *testing.T) {
    test_compilation(t, "var b; var a =1;")
}

func TestCompileBlock(t *testing.T) {
    test_compilation(t, "{print \"123\";}")
}

func TestBlockVars(t *testing.T) {
    test_compilation(t, "{var a; var b = 1; a = 2;}")
}
