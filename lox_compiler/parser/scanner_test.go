package parser_test

import (
	"lox-compiler/parser"
	"testing"
)

func assertTokenTypesMatch(t *testing.T, expected []parser.TokenType, tokens []parser.Token) {
    if len(expected) != len(tokens) {
        t.Fatalf("Got more tokens than expected. Expected %v, got %v", expected, tokens)
    }
	for i, v := range expected {
		if tokens[i].Token_type != v {
			t.Fatalf("Expected '%s' but got '%s'", v.String(), tokens[i].Token_type.String())
		}
	}
}

// Calls the scanner with a simple binary statement
func TestScanBinary(t *testing.T) {
	expectedTokens := []parser.TokenType{parser.NUMBER, parser.PLUS, parser.NUMBER, parser.SEMICOLON, parser.EOF}
	toks, err := parser.Scan("2+2;")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

    assertTokenTypesMatch(t, expectedTokens, toks)
}

func TestInvalidStr(t *testing.T) {
	expectedTokens := []parser.TokenType{parser.ERROR, parser.EOF}
	toks, err := parser.Scan("\"asdfasdfasdf;")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

    assertTokenTypesMatch(t, expectedTokens, toks)
}
