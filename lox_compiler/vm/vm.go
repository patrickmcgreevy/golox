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

type runtimeErrorCode int

//go:generate stringer -type=runtimeErrorCode -linecomment
const (
	outOfBoundsPC runtimeErrorCode = iota // out of bounds program counter
	popEmptyStack                         // pop on an empty stack
)

type RuntimeError struct {
	errCode runtimeErrorCode
}

func (e RuntimeError) Error() string {
    return fmt.Sprintf("encountered a(n) %v error", e.errCode)
}

func (vm *VirtualMachine) Interpret(c *bytecode.Chunk) InterpreterResult {
	vm.chunk = *c
	// vm.chunk.Disassemble("main")

	return vm.run()
}

// This is a performance critical path. There are techniques to speed it up.
// If you want to learn some of these techniques, look up “direct threaded code”, “jump table”, and “computed goto”.
func (vm *VirtualMachine) run() InterpreterResult {
    var err *RuntimeError
    var inst bytecode.Instruction

    for inst, err = vm.read_inst(); err == nil; inst, err = vm.read_inst() {
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
		case bytecode.OpNegate:
			// vm.chunk.Values.Push(-vm.chunk.Values.Pop())
            vm.chunk.Values[len(vm.chunk.Values)-1] = -vm.chunk.Values[len(vm.chunk.Values)-1]
		case bytecode.OpAdd, bytecode.OpSubtract, bytecode.OpMultiply, bytecode.OpDivide:
			vm.run_binary_op(inst)
		}

	}

    if err != nil {
        if err.errCode == outOfBoundsPC {
            return Interpret_OK
        }
        return Interpret_Runtime_Error
    }

    return Interpret_OK
}

func (vm *VirtualMachine) read_inst() (bytecode.Instruction, *RuntimeError) {
    if vm.pc >= len(vm.chunk.InstructionSlice) {
        return bytecode.Instruction{}, &RuntimeError{errCode: outOfBoundsPC}
    }
	i := vm.chunk.InstructionSlice[vm.pc]
	vm.pc += 1
	return i, nil
}

func (vm VirtualMachine) read_const(i bytecode.Instruction) bytecode.Value {
	return vm.chunk.Constants[i.Operands[0]]
}

func (vm *VirtualMachine) run_binary_op(i bytecode.Instruction) {
	var ret bytecode.Value
	r, l := vm.chunk.Values.Pop(), vm.chunk.Values.Pop()

	switch i.Code {
	case bytecode.OpAdd:
		ret = l + r
	case bytecode.OpSubtract:
		ret = l - r
	case bytecode.OpMultiply:
		ret = l * r
	case bytecode.OpDivide:
		ret = l / r
	}

	vm.chunk.Values.Push(ret)
}
