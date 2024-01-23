package interpreter

import (
	"fmt"
	"golox/scanner"
)

var constructor_name string = "init"

type LoxClass struct {
	Name    string
	Methods map[string]UserCallable
}

func (c LoxClass) String() string {
	return c.Name
}

func (c LoxClass) Call(interp Interpreter, args []any) (any, *RuntimeError) {
    instance := NewLoxInstance(c)
    init, ok := c.Methods[constructor_name]
    if ok {
        init.Bind(instance)
        init.Call(interp, args)
    }
	return instance, nil
}
func (c LoxClass) Arity() int {
    val, ok := c.Methods[constructor_name]
    if ok {
        return val.Arity()
    }

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
