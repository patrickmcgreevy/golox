package main

import (
	"fmt"
	"lox-compiler/bytecode"
	"lox-compiler/vm"
)

func main() {
    vm := vm.VirtualMachine{}
	c := bytecode.NewChunk()
    c.AddInst(bytecode.NewConstantInst(bytecode.Operand(c.AddConstant(1.2)), 0))
    for i := 0; i < 100000; i += 1 {
        c.AddInst(bytecode.NewNegateInst(0))
    }
    // c.AddInst(bytecode.NewConstantInst(bytecode.Operand(c.AddConstant(3.4)), 0))
    // c.AddInst(bytecode.NewAddInst(0))
    // c.AddInst(bytecode.NewConstantInst(bytecode.Operand(c.AddConstant(5.6)), 0))
    // c.AddInst(bytecode.NewDivideInst(0))
    // c.AddInst(bytecode.NewNegateInst(0))
	c.AddInst(bytecode.NewReturnInst(0))
    fmt.Println(vm.Interpret(&c))
}
