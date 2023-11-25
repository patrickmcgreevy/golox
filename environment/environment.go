package environment

import (
	"fmt"
	"golox/token"
)

type undefinedVariableError struct {
    name token.Token
}

// func (err undefinedVariableError) Error() string {
//     return fmt.Sprintf("'%s' is not defined in this environment.", err.name)
// }

func (err undefinedVariableError) Error() string {
    return fmt.Sprintf("'%s' is not defined in this environment.", err.name.Lexeme)
}

func newUndefinedVariableError(name token.Token) undefinedVariableError {
    return undefinedVariableError{name: name}
}


type environment struct {
    values map[string]any
}

func NewEnvironment() environment {
    return environment{values: make(map[string]any)}
}

func (e *environment) Define(name string, value any) {
    e.values[name] = value
}

func (e environment) Get(name token.Token) (any, error) {
    var err undefinedVariableError
    val, ok := e.values[name.Lexeme]
    if ok {
        return val, nil
    }
    err = newUndefinedVariableError(name)
    return nil, err
}
