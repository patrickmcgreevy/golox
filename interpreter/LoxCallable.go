package interpreter

// import "golox/interpreter"



type LoxCallable interface {
    Call(interp Interpreter, args []any) any
    Arity() int
}
