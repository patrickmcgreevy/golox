package vm

import (
	"fmt"
	"lox-compiler/bytecode"
	"lox-compiler/compiler"
	"lox-compiler/debug"
)

type InterpreterResult int

//go:generate stringer -type=InterpreterResult
const (
	Interpret_OK InterpreterResult = iota
	Interpret_Compile_Error
	Interpret_Runtime_Error
)

type VirtualMachine struct {
	chunk           bytecode.Chunk
	pc              int
	InteractiveMode bool
	vars            map[bytecode.LoxString]bytecode.Value
}

type runtimeErrorCode int

//go:generate stringer -type=runtimeErrorCode -linecomment
const (
	outOfBoundsPC runtimeErrorCode = iota // out of bounds program counter
	popEmptyStack                         // pop on an empty stack
	wrongType                             // incorrect type
	invalidOpCode                         // invalid OpCode
)

type RuntimeError struct {
	errCode runtimeErrorCode
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("encountered a(n) %v error", e.errCode)
}

func (vm *VirtualMachine) Interpret(s string) InterpreterResult {
    if vm.vars == nil {
        vm.vars = make(map[bytecode.LoxString]bytecode.Value)

    }
	vm.pc = 0
	c := compiler.Compiler{}
	c.InteractiveMode = vm.InteractiveMode
	chunk, err := c.Compile(s)
	if err != nil {
		fmt.Println(err.Error())
		return Interpret_Compile_Error
	}

	return vm.run_bytecode(chunk)
}

func (vm *VirtualMachine) run_bytecode(c *bytecode.Chunk) InterpreterResult {
	vm.chunk = *c

	return vm.run()
}

// This is a performance critical path. There are techniques to speed it up.
// If you want to learn some of these techniques, look up “direct threaded code”, “jump table”, and “computed goto”.
func (vm *VirtualMachine) run() InterpreterResult {
	var err *RuntimeError
	var inst bytecode.Instruction

	for inst, err = vm.read_inst(); err == nil; inst, err = vm.read_inst() {
		debug.Printf("%v", vm.chunk.Values)
		debug.Printf("%s", inst.String())
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
			loxInt, ok := vm.chunk.Values[len(vm.chunk.Values)-1].(bytecode.LoxInt)
			if ok {
				vm.chunk.Values[len(vm.chunk.Values)-1] = -loxInt
			} else {
				vm.chunk.Values[len(vm.chunk.Values)-1] = bytecode.LoxBool(!vm.chunk.Values[len(vm.chunk.Values)-1].Truthy())
			}
		case bytecode.OpLess:
			rInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return Interpret_Runtime_Error
			}
			lInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return Interpret_Runtime_Error
			}
			vm.chunk.Values.Push(bytecode.LoxBool(lInt < rInt))

		case bytecode.OpLessEqual:
			rInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return Interpret_Runtime_Error
			}
			lInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return Interpret_Runtime_Error
			}
			vm.chunk.Values.Push(bytecode.LoxBool(lInt <= rInt))
		case bytecode.OpGreater:
			rInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return Interpret_Runtime_Error
			}
			lInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return Interpret_Runtime_Error
			}
			vm.chunk.Values.Push(bytecode.LoxBool(lInt > rInt))
		case bytecode.OpGreaterEqual:
			rInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return Interpret_Runtime_Error
			}
			lInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return Interpret_Runtime_Error
			}
			vm.chunk.Values.Push(bytecode.LoxBool(lInt >= rInt))
		case bytecode.OpEqualEqual:
			r := vm.chunk.Values.Pop()
			l := vm.chunk.Values.Pop()
			vm.chunk.Values.Push(bytecode.LoxBool(l == r))
		case bytecode.OpNotEqual:
			r := vm.chunk.Values.Pop()
			l := vm.chunk.Values.Pop()
			vm.chunk.Values.Push(bytecode.LoxBool(l != r))
		case bytecode.OpAdd, bytecode.OpSubtract, bytecode.OpMultiply, bytecode.OpDivide:
			err = vm.run_binary_op(inst)
			if err != nil {
				return Interpret_Runtime_Error
			}
		case bytecode.OpPrint:
			fmt.Println(vm.chunk.Values.Pop())
		case bytecode.OpDeclare:
			//
            val := vm.chunk.Values.Pop()
            name, ok := val.(bytecode.LoxString)
            if !ok {
                return Interpret_Runtime_Error
            }
            vm.vars[name] = nil
        case bytecode.OpAssign:
            // pop name
            name, ok := vm.chunk.Values.Pop().(bytecode.LoxString)
            if !ok {
                return Interpret_Runtime_Error
            }
            // pop val
            val := vm.chunk.Values.Pop()
            vm.vars[name] = val
        case bytecode.OpLookup:
            name, ok := vm.chunk.Values.Pop().(bytecode.LoxString)
            if !ok {
                return Interpret_Runtime_Error
            }
            val, ok := vm.vars[name]
            if !ok {
                return Interpret_Runtime_Error
            }

            vm.chunk.Values.Push(val)

		default:
			fmt.Println("unknown instruction ", inst.String())
			return Interpret_Runtime_Error
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

func (vm *VirtualMachine) run_logical_op(i bytecode.Instruction) *RuntimeError {
	var ret bytecode.Value
	rVal, lVal := vm.chunk.Values.Pop(), vm.chunk.Values.Pop()

	switch i.Code {
	case bytecode.OpOr:
		ret = bytecode.LoxBool(rVal.Truthy() || lVal.Truthy())
	case bytecode.OpAdd:
		ret = bytecode.LoxBool(rVal.Truthy() || lVal.Truthy())
	}

	vm.chunk.Values.Push(ret)
	return nil
}

func (vm *VirtualMachine) run_binary_op(i bytecode.Instruction) *RuntimeError {
	var ret bytecode.Value
	rVal, lVal := vm.chunk.Values.Pop(), vm.chunk.Values.Pop()
	lInt, lOK := lVal.(bytecode.LoxInt)
	rInt, rOK := rVal.(bytecode.LoxInt)
	if !lOK || !rOK {
		lStr, lOK := lVal.(bytecode.LoxString)
		rStr, rOK := rVal.(bytecode.LoxString)
		if (!lOK || !rOK) || i.Code != bytecode.OpAdd {
			// error!!
			// Only + supports str and int other sneed int
			debug.Printf("line[]: expected integers but got (%T, %T)", rVal, lVal)
			// return fmt.Errorf()
			return &RuntimeError{errCode: wrongType}
		} else {
			ret = lStr + rStr
		}
	} else {
		switch i.Code {
		case bytecode.OpAdd:
			ret = lInt + rInt
		case bytecode.OpSubtract:
			ret = lInt - rInt
		case bytecode.OpMultiply:
			ret = lInt * rInt
		case bytecode.OpDivide:
			ret = lInt / rInt
		default:
			return &RuntimeError{errCode: invalidOpCode}
		}
	}

	vm.chunk.Values.Push(ret)
	return nil
}
