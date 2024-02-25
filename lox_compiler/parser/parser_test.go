package parser_test

import (
	"lox-compiler/parser"
	"testing"
)

func TestBinaryStmt(t *testing.T) {
	toks, err := parser.Scan("2+2;")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	p := parser.NewParser(toks)
	stmts := p.Parse()
	if len(stmts) > 1 {
		t.Fatal(stmts)
	}
	s, _ := stmts[0].(parser.ExpressionStmt)
	b, ok := s.Val.(parser.Binary)
	if !ok {
		t.Fatalf("Expected %s but got %T", "parser.Binary", stmts[0])
	}
	if b.String() != "PLUS 2 2" {
		t.Fatal(b)
	}
}

func TestBinaryPrecedence(t *testing.T) {
	toks, _ := parser.Scan("2+2*3-8;")
	p := parser.NewParser(toks)
	stmts := p.Parse()
	if stmts[0].(parser.ExpressionStmt).Val.(parser.Binary).String() != "PLUS 2 MINUS STAR 2 3 8" {
		t.Fatalf("%v", stmts)
	}
}

func TestBinaryPrecedence2(t *testing.T) {
	toks, _ := parser.Scan("8*2-2/4+10;")
	p := parser.NewParser(toks)
	stmts := p.Parse()
	if stmts[0].(parser.ExpressionStmt).Val.(parser.Binary).String() != "MINUS STAR 8 2 PLUS SLASH 2 4 10" {
		t.Fatalf("%v", stmts)
	}
}

func TestCall(t *testing.T) {
	toks, _ := parser.Scan("foo.bar();")
	p := parser.NewParser(toks)
	stmts := p.Parse()
	if stmts[0].(parser.ExpressionStmt).Val.(parser.Call).String() != "CALL GET foo.bar ([])" {
		t.Fatalf("%v", stmts)
	}
}
