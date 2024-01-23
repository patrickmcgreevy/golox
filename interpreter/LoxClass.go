package interpreter

import (
	"fmt"
	"golox/scanner"
)

type LoxClass struct {
	Name    string
	Methods map[string]UserCallable
}

func (c LoxClass) String() string {
	return c.Name
}

func (c LoxClass) Call(interp Interpreter, args []any) (any, *RuntimeError) {
	return NewLoxInstance(c), nil
}
func (c LoxClass) Arity() int {
	return 0
}

type LoxInstance struct {
	Class  LoxClass
	Fields map[string]any
}

func NewLoxInstance(class LoxClass) LoxInstance {
	fields := make(map[string]any)
	// for k, f := range class.Methods {
	// 	fields[k] = f
	// }
	return LoxInstance{Class: class, Fields: fields}
}

func (inst LoxInstance) String() string {
	return inst.Class.String() + " instance"
}

func (inst LoxInstance) Get(name scanner.Token) (any, *RuntimeError) {
	val, ok := inst.Fields[name.Lexeme]
	if ok {
        return val, nil
	}
    
    method, ok := inst.Class.Methods[name.Lexeme]
    if ok {
        method.Bind(inst)
        return method, nil
    }

    return nil, &RuntimeError{
        error: fmt.Sprintf(
            "\"%s\" is not a property of \"%s\"", name.Lexeme, inst.Class),
        tok: name,
    }
}

func (inst *LoxInstance) Set(name scanner.Token, val any) {
	inst.Fields[name.Lexeme] = val
}
