package main

import (
	"fmt"
	"lox-compiler/bytecode"
	"lox-compiler/vm"
)

func main() {
    vm := vm.VirtualMachine{}
	c := bytecode.NewChunk()
	c.AddInst(bytecode.NewConstantInst(bytecode.Operand(c.AddConstant(32)), 0))
	c.AddInst(bytecode.NewConstantInst(bytecode.Operand(c.AddConstant(64)), 0))
	c.AddInst(bytecode.NewReturnInst(0))
	// fmt.Println(c)
	// c.Disassemble("main bytecode")
    fmt.Println(vm.Interpret(&c))
}
