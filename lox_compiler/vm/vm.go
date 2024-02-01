package vm

import (
	"fmt"
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
	pc    int
}

func (vm *VirtualMachine) Interpret(c *bytecode.Chunk) InterpreterResult {
	vm.chunk = *c
	vm.chunk.Disassemble("main")

	return vm.run()
}

func (vm *VirtualMachine) run() InterpreterResult {
	// This is a performance critical path. There are techniques to speed it up.
	// If you want to learn some of these techniques, look up “direct threaded code”, “jump table”, and “computed goto”.
	var inst bytecode.Instruction

	for inst = vm.read_inst(); ; inst = vm.read_inst() {
        debug("%v", vm.chunk.Values)
		debug("%s", inst.String())
		switch inst.Code {
		case bytecode.OpReturn:
            fmt.Println(vm.chunk.Values.Pop())
			return Interpret_OK
		case bytecode.OpConstant:
			// We could define some type aliases and methods on those aliases for each
			// Instruction type?? Would this be slow as balls? Any good?
			// fmt.Println(vm.chunk.Constants[inst.Operands[0]])
            vm.chunk.Values.Push(vm.read_const(inst))
		}

	}
}

func (vm *VirtualMachine) read_inst() bytecode.Instruction {
	i := vm.chunk.InstructionSlice[vm.pc]
	vm.pc += 1
	return i
}

func (vm VirtualMachine) read_const(i bytecode.Instruction) bytecode.Value {
    return vm.chunk.Constants[i.Operands[0]]
}
