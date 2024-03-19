package compiler_test

import (
	"fmt"
	"lox-compiler/compiler"
	"strings"
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

func TestLocalRedefine(t *testing.T) {
    c := compiler.Compiler{}
    chunk, err := c.Compile("{var a; var a;}")
    if err == nil {
        chunk.Disassemble("main")
        t.Fatalf("expected compilation to fail")
    }
}

func TestExceedMaxLocalVars(t *testing.T) {
    str := strings.Builder{}
    c := compiler.Compiler{}
    str.WriteString("{")
    for i := 0; i < 1000; i++ {
        str.WriteString(fmt.Sprintf("var a%d = 1;\n", i))
    }
    str.WriteString("}")
    chunk, err := c.Compile(str.String())
    if err == nil {
        chunk.Disassemble("main")
        t.Fatalf("expected compilation to fail due to too many local variables")
    }
}

func TestIf(t *testing.T) {
    test_compilation(t, "if (true) {print 1;}")
    test_compilation(t, "if (true) {print 1;} else {print 2;}")
}
