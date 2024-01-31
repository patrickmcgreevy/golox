package vm

import (
	"lox-compiler/bytecode"
)

type InterpreterResult int

//go:generate stringer -type=InterpreterResult
const (
	Interpret_OK InterpreterResult = iota
	Interpret_Compile_Error
	Interpret_Runtime_Error
)

type VirtualMachine struct {
	chunk bytecode.Chunk
}

func (vm *VirtualMachine) Interpret(c *bytecode.Chunk) InterpreterResult {
    return Interpret_OK
}

func (vm *VirtualMachine) run()
