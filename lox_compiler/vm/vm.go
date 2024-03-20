package vm

import (
	"fmt"
	"lox-compiler/bytecode"
	"lox-compiler/compiler"
	"lox-compiler/debug"
	"strings"
)

type VirtualMachine struct {
	chunk           bytecode.Chunk
	pc              int
	InteractiveMode bool
	vars            map[bytecode.LoxString]bytecode.Value
}

const (
	outOfBoundsPC string = "out of bounds program counter"
	popEmptyStack        = "pop on an empty stack"
	wrongType            = "incorrect type"
	invalidOpCode        = "invalid OpCode"
	expectedInts         = "expected two ints"
	expectedStr          = "expected a string"
)

type InterpreterError struct {
	interpreterErr string
	line           int
}

func (e InterpreterError) Error() string {
	str := strings.Builder{}
	if e.line >= 0 {
		str.WriteString(fmt.Sprintf("[line %d]: ", e.line))
	}
    str.WriteString(fmt.Sprintf("encountered an error: %s", e.interpreterErr))

	return str.String()
}

func (vm *VirtualMachine) Interpret(s string) *InterpreterError {
	if vm.vars == nil {
		vm.vars = make(map[bytecode.LoxString]bytecode.Value)

	}
	vm.pc = 0
	c := compiler.Compiler{}
	c.InteractiveMode = vm.InteractiveMode
	chunk, err := c.Compile(s)
	if err != nil {
		return &InterpreterError{interpreterErr: err.Error(), line: -1}
	}

	return vm.run_bytecode(chunk)
}

func (vm *VirtualMachine) run_bytecode(c *bytecode.Chunk) *InterpreterError {
	vm.chunk = *c

	return vm.run()
}

// This is a performance critical path. There are techniques to speed it up.
// If you want to learn some of these techniques, look up “direct threaded code”, “jump table”, and “computed goto”.
func (vm *VirtualMachine) run() *InterpreterError {
	var err *InterpreterError
	var inst bytecode.Instruction

	for inst, err = vm.read_inst(); err == nil; inst, err = vm.read_inst() {
		debug.Printf("%s", inst.String())
		switch inst.Code {
		case bytecode.OpReturn:
			fmt.Println(vm.chunk.Values.Pop())
			return nil

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
				return &InterpreterError{interpreterErr: expectedInts, line: inst.SourceLineNumer}
			}
			lInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return &InterpreterError{interpreterErr: expectedInts, line: inst.SourceLineNumer}
			}
			vm.chunk.Values.Push(bytecode.LoxBool(lInt < rInt))

		case bytecode.OpLessEqual:
			rInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return &InterpreterError{interpreterErr: expectedInts, line: inst.SourceLineNumer}
			}
			lInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return &InterpreterError{interpreterErr: expectedInts, line: inst.SourceLineNumer}
			}
			vm.chunk.Values.Push(bytecode.LoxBool(lInt <= rInt))

		case bytecode.OpGreater:
			rInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return &InterpreterError{interpreterErr: expectedInts, line: inst.SourceLineNumer}
			}
			lInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return &InterpreterError{interpreterErr: expectedInts, line: inst.SourceLineNumer}
			}
			vm.chunk.Values.Push(bytecode.LoxBool(lInt > rInt))

		case bytecode.OpGreaterEqual:
			rInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return &InterpreterError{interpreterErr: expectedInts, line: inst.SourceLineNumer}
			}
			lInt, ok := vm.chunk.Values.Pop().(bytecode.LoxInt)
			if !ok {
				return &InterpreterError{interpreterErr: expectedInts, line: inst.SourceLineNumer}
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
				return err
			}

		case bytecode.OpPrint:
			fmt.Println(vm.chunk.Values.Pop())

		case bytecode.OpDeclareGlobal:
			//
			val := vm.chunk.Values.Pop()
			name, ok := val.(bytecode.LoxString)
			if !ok {
				return &InterpreterError{interpreterErr: expectedStr, line: inst.SourceLineNumer}
			}
			vm.vars[name] = nil

		case bytecode.OpAssign:
			// pop name
			name, ok := vm.chunk.Values.Pop().(bytecode.LoxString)
			if !ok {
				return &InterpreterError{interpreterErr: expectedStr, line: inst.SourceLineNumer}
			}
            // Don't pop the value, because an expression needs a result
			vm.vars[name] = vm.chunk.Values[len(vm.chunk.Values)-1]

		case bytecode.OpGlobalLookup:
			name, ok := vm.chunk.Values.Pop().(bytecode.LoxString)
			if !ok {
				return &InterpreterError{interpreterErr: expectedStr, line: inst.SourceLineNumer}
			}
			val, ok := vm.vars[name]
			if !ok {
				return &InterpreterError{interpreterErr: fmt.Sprintf("variable %s is not defined in this scope", name), line: inst.SourceLineNumer}
			}
			vm.chunk.Values.Push(val)

		case bytecode.OpLocalLookup:
			vm.chunk.Values.Push(vm.chunk.Values[inst.Operands[0]])

		case bytecode.OpLocalAssign:
            // Don't pop the value, that's the result of the assignment expression
			vm.chunk.Values[inst.Operands[0]] = vm.chunk.Values[len(vm.chunk.Values)-1]

		case bytecode.OpPop:
			vm.chunk.Values.Pop()

        case bytecode.OpConditionalJump:
            cond := vm.chunk.Values.Pop().Truthy()
            if !cond {
                vm.pc += int(vm.chunk.Constants[inst.Operands[1]].(bytecode.LoxInt))
            } 

        case bytecode.OpJump:
            vm.pc += int(vm.chunk.Constants[inst.Operands[0]].(bytecode.LoxInt))

        case bytecode.OpAnd, bytecode.OpOr:
            err := vm.run_logical_op(inst)
            if err != nil {
                return err
            }

		default:
			fmt.Println("unknown instruction ", inst.String())
			return &InterpreterError{interpreterErr: "unkown instruction", line: inst.SourceLineNumer}
		}
		debug.Printf("%v", vm.chunk.Values)

	}

	if err != nil {
		if err.interpreterErr == outOfBoundsPC {
			return nil
		}
		return &InterpreterError{interpreterErr: outOfBoundsPC, line: inst.SourceLineNumer}
	}

	return nil
}

func (vm *VirtualMachine) read_inst() (bytecode.Instruction, *InterpreterError) {
	if vm.pc >= len(vm.chunk.InstructionSlice) {
		return bytecode.Instruction{}, &InterpreterError{interpreterErr: outOfBoundsPC}
	}
	i := vm.chunk.InstructionSlice[vm.pc]
	vm.pc += 1
	return i, nil
}

func (vm VirtualMachine) read_const(i bytecode.Instruction) bytecode.Value {
	return vm.chunk.Constants[i.Operands[0]]
}

func (vm *VirtualMachine) run_logical_op(i bytecode.Instruction) *InterpreterError {
	var ret bytecode.Value
	rVal, lVal := vm.chunk.Values.Pop(), vm.chunk.Values.Pop()

	switch i.Code {
	case bytecode.OpOr:
		ret = bytecode.LoxBool(rVal.Truthy() || lVal.Truthy())
	case bytecode.OpAnd:
		ret = bytecode.LoxBool(rVal.Truthy() && lVal.Truthy())
    default:
        return &InterpreterError{interpreterErr: "invalid opcode"}
	}

	vm.chunk.Values.Push(ret)
	return nil
}

func (vm *VirtualMachine) run_binary_op(i bytecode.Instruction) *InterpreterError {
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
			return &InterpreterError{interpreterErr: wrongType}
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
			return &InterpreterError{interpreterErr: invalidOpCode}
		}
	}

	vm.chunk.Values.Push(ret)
	return nil
}
