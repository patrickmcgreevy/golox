package main

import (
	"fmt"
	"golox/errorhandling"
	"golox/expression"
	"golox/parser"
	"golox/scanner"
	"golox/token"
	"os"
	// "runtime"
)

var hadError bool // Improvement idea: Implement an ErrorHandling interface so we can pass different strategies

func main() {
	test_expr()
	if len(os.Args) > 2 {
		panic("Need two or more args")
	} else if len(os.Args) == 2 {
		fmt.Println(os.Args)
		file, err := os.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		fmt.Println(string(file))
		run(string(file))
	} else {
		runPrompt()
	}
}

func test_expr() {
	expr := expression.Binary{
		Left: expression.Binary{
			Left: expression.Unary{
				Operator: token.Token{Token_type: token.BANG, Lexeme: "!", Literal: nil, Line: 0},
				Right:    expression.Literal{Value: 420}},
			Operator: token.Token{Token_type: token.MINUS, Lexeme: "-", Literal: nil, Line: 0},
			Right:    expression.Grouping{Expr: expression.Literal{Value: "patrick"}}},
		Operator: token.Token{Token_type: token.PLUS, Lexeme: "+", Literal: nil, Line: 0},
		Right:    expression.Literal{Value: "def"},
	}

	fmt.Println(expr.Expand_to_string())

	// Expr expression = new Expr.Binary(
	//         new Expr.Unary(
	//             new Token(TokenType.MINUS, "-", null, 1),
	//             new Expr.Literal(123)),
	//         new Token(TokenType.STAR, "*", null, 1),
	//         new Expr.Grouping(
	//             new Expr.Literal(45.67)));
	// (* (- 123) (group 45.67))

	expr2 := expression.Binary{
		Left: expression.Unary{
			Operator: token.Token{Token_type: token.MINUS, Lexeme: "-", Literal: nil, Line: 0},
			Right:    expression.Literal{Value: 123},
		},
		Operator: token.Token{Token_type: token.STAR, Lexeme: "*", Literal: nil, Line: 0},
		Right: expression.Grouping{
			Expr: expression.Literal{Value: 45.67},
		},
	}

    fmt.Println(expr2.Expand_to_string())

}

func runFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fmt.Println(file)
	run(string(file))

	if hadError {
		os.Exit(65)
	}

	return nil
}

func runPrompt() error {
	var line string
	for {
		fmt.Print("> ")
		_, err := fmt.Scanln(&line)
		if err != nil {
			panic(err)
		}
		run(line)
		hadError = false
	}
}

func run(source string) {
	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()

	for _, v := range tokens {
		fmt.Println(v)
	}
    parser := parser.NewParser(tokens)

    // runtime.Breakpoint()
    expr, _ := parser.Expression()
    fmt.Println(expr.Expand_to_string())
}

func raise_error(line int, message string) {
	errorhandling.Report(line, "", message)
}

type Scanner struct {
}

func (s *Scanner) scanTokens() []token.Token {
	return []token.Token{}
}
