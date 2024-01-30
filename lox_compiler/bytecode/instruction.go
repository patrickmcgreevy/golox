package bytecode

import (
	"fmt"
)

type OpCode uint8

//go:generate stringer -type=OpCode
const (
	OpReturn OpCode = iota
    OpConstant
)

type Instruction struct {
	Code OpCode
}

func (c Instruction) String() string {
	return fmt.Sprintf("%s", c.Code.String())
}

func (i Instruction) DisassembleInst() {
	fmt.Println(i)
}
