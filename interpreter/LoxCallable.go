package interpreter

import (
	"golox/statement"
)

type BuiltinCallable struct {
    params []string
	foo   func(interp Interpreter, args []any) any
}

func (c BuiltinCallable) Arity() int {
	return len(c.params)
}

func (c BuiltinCallable) Call(interp Interpreter, args []any) any {
	return c.foo(interp, args)
}

func (c BuiltinCallable) String() string {
    return "<native fn>"
}

type UserCallable struct {
    declaration statement.Function
    closure Environment
}

func (c UserCallable) Arity() int {
    return len(c.declaration.Params)
}

func (c UserCallable) Call(interp Interpreter, args []any) (any, *RuntimeError) {
   // Create env
   env := NewEnvironment()
   env.SetEnclosing(&c.closure)
   // Map param name to arg values
   for i, v := range args {
       env.Define(c.declaration.Params[i].Lexeme, v)
   }
   // Evaluate block
   interp.executeBlock(c.declaration.Body, env)

   return interp.val, interp.err
}

func (c UserCallable) String() string {
    return "<fn " + c.declaration.Name.Lexeme + ">"
}
