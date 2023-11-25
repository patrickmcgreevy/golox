package main

import (
	"bufio"
	"fmt"
	"golox/interpreter"
	"golox/parser"
	"golox/scanner"
	"os"
)

var hadError bool // Improvement idea: Implement an ErrorHandling interface so we can pass different strategies
var interp interpreter.Interpreter

func main() {
    interp = interpreter.NewInterpreter()

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
    reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
        line, err := reader.ReadString('\n')
		if err != nil {
            if err.Error() == "EOF" {
                return nil
            }
			panic(err)
		}
		run(line)
		hadError = false
	}
}

func run(source string) {
	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()
    parser := parser.NewParser(tokens)
    // interp := interpreter.NewInterpreter()

    statements := parser.Parse()
    if statements == nil {
        return
    }

    interp.Interpret(statements)
}
