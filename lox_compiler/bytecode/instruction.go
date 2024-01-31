package bytecode

import (
	"fmt"
	"strings"
)

type OpCode uint8
type Operand int
type OperandArray [2]Operand

//go:generate stringer -type=OpCode
const (
	OpReturn OpCode = iota
	OpConstant
)

type Instruction struct {
	Code     OpCode
	Operands OperandArray
    SourceLineNumer int
}

func NewReturnInst(line int) Instruction {
	return Instruction{Code: OpReturn, SourceLineNumer: line}
}

func NewConstantInst(val Operand, line int) Instruction {
	ret := Instruction{Code: OpConstant, SourceLineNumer: line}
	ret.Operands[0] = val

	return ret
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
