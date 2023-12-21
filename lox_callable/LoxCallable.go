package loxcallable

import "golox/interpreter"



type LoxCallable interface {
    Call(interp interpreter.Interpreter, args []any) any
    Arity() int
}
