package main

import (
	"lox-compiler/bytecode"
)

func main() {
	c := bytecode.NewChunk()
	c.InstructionSlice = append(c.InstructionSlice, bytecode.Instruction{Op: bytecode.OpReturn})
    c.AddConstant(32)
	// fmt.Println(c)
	c.Disassemble("main bytecode")
}
