package bytecode

import (
	"fmt"
)

type InstructionSlice []Instruction
type Chunk struct {
	InstructionSlice
	Constants ValueSlice
}

func NewChunk() Chunk {
	return Chunk{
		InstructionSlice: make(InstructionSlice, 0),
		Constants:        make(ValueSlice, 0),
	}
}

func (c *Chunk) AddConstant(v Value) int {
	return c.Constants.addConstant(v)
}

func (c *Chunk) AddInst(i Instruction) {
    c.InstructionSlice = append(c.InstructionSlice, i)
}

func (c Chunk) Disassemble(name string) error {
	fmt.Println(fmt.Sprintf("== %s ==", name))
	fmt.Println(c.Constants)
	for i, v := range c.InstructionSlice {
		_, err := fmt.Printf("%04d %-4d  %s\n", i, v.SourceLineNumer, v.String())
		if err != nil {
			return err
		}
	}

	return nil
}
