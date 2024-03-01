package bytecode

import (
	"fmt"
	"strings"
)

type OpCode uint8
type Operand int
type OperandArray [1]Operand

//go:generate stringer -type=OpCode
const (
	OpReturn OpCode = iota
	OpConstant
	OpNegate
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
    OpPrint
    OpOr
    OpAnd
    OpLess
    OpGreater
    OpLessEqual
    OpGreaterEqual
    OpEqualEqual
    OpNotEqual
)

type Instruction struct {
	Code            OpCode
	Operands        OperandArray
	SourceLineNumer int
}

func NewPrintInst(line int) Instruction {
    return Instruction{Code: OpPrint, SourceLineNumer: line}
}

func NewReturnInst(line int) Instruction {
	return Instruction{Code: OpReturn, SourceLineNumer: line}
}

func NewConstantInst(val Operand, line int) Instruction {
	ret := Instruction{Code: OpConstant, SourceLineNumer: line}
	ret.Operands[0] = val

	return ret
}

func NewNegateInst(line int) Instruction {
	return Instruction{Code: OpNegate, SourceLineNumer: line}
}

func NewAddInst(line int) Instruction {
	return Instruction{Code: OpAdd, SourceLineNumer: line}
}

func NewSubtractInst(line int) Instruction {
	return Instruction{Code: OpSubtract, SourceLineNumer: line}
}

func NewMultiplyInst(line int) Instruction {
	return Instruction{Code: OpMultiply, SourceLineNumer: line}
}

func NewDivideInst(line int) Instruction {
	return Instruction{Code: OpDivide, SourceLineNumer: line}
}

func (c Instruction) String() string {
	return fmt.Sprintf("%-16s %v", c.Code.String(), c.Operands)
}

func (op Operand) String() string {
	return fmt.Sprintf("%04d", op)
}

func (a OperandArray) String() string {
	var ret strings.Builder

	for i, v := range a {
		if i > 0 {
			ret.WriteString(" ")
		}
		ret.WriteString(v.String())
	}

	return ret.String()
}

func (i Instruction) DisassembleInst() {
	fmt.Println(i)
}
