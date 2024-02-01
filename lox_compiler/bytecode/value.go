package bytecode

import "fmt"

type Value float64
type ValueSlice []Value
type ValueStack []Value

func (vs ValueSlice) String() string {
	return fmt.Sprintf("Constants: %v", []Value(vs))
}

func (vs *ValueSlice) addConstant(v Value) int {
	*vs = append(*vs, v)

	return len(*vs)-1
}

func (vs *ValueStack) Push(v Value) {
	*vs = append(*vs, v)
}

func (vs *ValueStack) Pop() Value {
    ret := (*vs)[len(*vs)-1]
    *vs = (*vs)[0:len(*vs)-1]

    return ret
}

func (vs *ValueStack) Reset() {
    *vs = make(ValueStack, 0)
}
