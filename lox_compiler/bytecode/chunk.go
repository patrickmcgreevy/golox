package bytecode

import (
	"fmt"
	"strings"
)

type InstructionSlice []Instruction
type Chunk struct {
	InstructionSlice
	Constants ValueSlice
	Values    ValueStack
}

func NewChunk() Chunk {
	return Chunk{
		InstructionSlice: make(InstructionSlice, 0),
		Constants:        make(ValueSlice, 0),
	}
}

func (c Chunk) String() string {
    str := strings.Builder{}
    str.WriteString(fmt.Sprintf("Constants: %s\n", c.Constants))
	for i, v := range c.InstructionSlice {
		str.WriteString(fmt.Sprintf("%04d %-4d  %s\n", i, v.SourceLineNumer, v.String()))
    }

    return str.String()
}

// Add a Value to the Chunk's pool of constants.
func (c *Chunk) AddConstant(v Value) int {
	return c.Constants.addConstant(v)
}

// Append an Instruction to the Chunk
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
