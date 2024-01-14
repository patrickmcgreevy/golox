package interpreter

import (
	"fmt"
	"golox/scanner"
)

type undefinedVariableError struct {
	name string
}

func (err undefinedVariableError) Error() string {
	return fmt.Sprintf("'%s' is not defined.", err.name)
}

func newUndefinedVariableError(name string) undefinedVariableError {
	return undefinedVariableError{name: name}
}

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment() Environment {
	return Environment{values: make(map[string]any), enclosing: nil}
}

func (e *Environment) SetEnclosing(enclosing *Environment) {
	e.enclosing = enclosing
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name string, value any) error {
	_, ok := e.values[name]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.Assign(name, value)
		}
		return newUndefinedVariableError(name)
	}
	e.values[name] = value

	return nil
}

func (e *Environment) AssignAt(dist int, name string, value any) error {
    var env *Environment
    var i int 

    for i, env = 0, e; i<=dist && env != nil; i, env = i+1, env.GetEnclosing() {
        if i != dist {
            continue
        }

        env.Assign(name, value)
        return nil
    }

    return newUndefinedVariableError(name)
}

func (e Environment) Get(name scanner.Token) (any, error) {
	val, ok := e.values[name.Lexeme]
	if ok {
		return val, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, newUndefinedVariableError(name.Lexeme) 
}

func (e Environment) GetAt(dist int, name scanner.Token) (any, error) {
    var env *Environment
    var i int

    for i, env = 0, &e; i <= dist && env != nil; i, env = i+1, env.GetEnclosing() {
        if i != dist {
            continue
        }

        return env.Get(name)
    }

    return nil, newUndefinedVariableError(name.Lexeme)
}

func (e Environment) GetEnclosing() *Environment {
	return e.enclosing
}
