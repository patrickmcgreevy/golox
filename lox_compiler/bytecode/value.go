package bytecode

import "fmt"

type Value float64
type ValueSlice []Value

func (vs ValueSlice) String() string {
    return fmt.Sprintf("Constants: %v", []Value(vs))
}

func (vs *ValueSlice) addConstant(v Value) int {
    *vs = append(*vs, v)

    return len(*vs)
}
