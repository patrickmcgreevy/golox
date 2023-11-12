package errorhandling

import "fmt"
import "golox/token"

func Report(line int, where, message string) {
	fmt.Printf("[line %d] Error '%s': '%s'\n", line, where, message)
}

func Error(cur_token token.Token, message string) {
	if cur_token.Token_type == token.EOF {
		Report(cur_token.Line, " at end ", message)
	} else {
		Report(cur_token.Line, " at '"+cur_token.Lexeme+"'", message)
	}
}
