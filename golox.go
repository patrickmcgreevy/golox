package main

import (
	"fmt"
	"os"
    "golox/token"
    "golox/scanner"
    "golox/errorhandling"
)

var hadError bool // Improvement idea: Implement an ErrorHandling interface so we can pass different strategies

func main() {
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

}

func raise_error(line int, message string) {
	errorhandling.Report(line, "", message)
}



type Scanner struct {
}

func (s *Scanner) scanTokens() []token.Token {
    return []token.Token{}
}
