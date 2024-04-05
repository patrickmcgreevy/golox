package bytecode

import "fmt"

type Value interface {
	Truthy() bool
	private()
}
type ValueSlice []Value
type ValueStack []Value

func NewValue(v any) (Value, error) {
	switch val := v.(type) {
	case float64:
		return LoxInt(val), nil
	case string:
		return LoxString(val), nil
	case bool:
		return LoxBool(val), nil
	case nil:
		return LoxNil(0), nil
	case LinearProbingHashMap:
		return LoxMap(val), nil
	}

	return nil, fmt.Errorf("%T is not a valid LoxValue type", v)
}

type LoxInt float64

func (v LoxInt) private() {}
func (v LoxInt) Truthy() bool {
	return v != 0
}

func (v LoxInt) Val() any {
	return v
}

type LoxString string

func (v LoxString) private() {}
func (v LoxString) Truthy() bool {
	return true
}

type LoxBool bool

func (v LoxBool) private() {}
func (v LoxBool) Truthy() bool {
	return bool(v)
}

type LoxNil int

func (v LoxNil) private() {}
func (v LoxNil) Truthy() bool {
	return false
}

func (v LoxNil) String() string {
	return "nil"
}

type LoxMap LinearProbingHashMap

func (v LoxMap) private() {}
func (v LoxMap) Truthy() bool {
	return true
}

type LoxFunc struct {
	Args  []LoxString
	Body  Chunk
	Name  LoxString
}

func NewLoxFunc(name string) LoxFunc {
	return LoxFunc{
		Args: make([]LoxString, 5),
		Body: NewChunk(),
        Name: LoxString(name),
	}
}

func (LoxFunc) private() {}
func (LoxFunc) Truthy() bool {
	return true
}

func (f LoxFunc) String() string {
	return fmt.Sprintf("fun %s(%v)",f.Name, f.Args)
}

func (f LoxFunc) Arity() int {
    return len(f.Args)
}

func (v *LoxMap) Insert(s LoxString, val Value) {
	(*LinearProbingHashMap)(v).Insert(s, val)
}

func (v *LoxMap) Get(s LoxString) (Value, error) {
	return (*LinearProbingHashMap)(v).Get(s)
}

func (v *LoxMap) Delete(s LoxString) {
	(*LinearProbingHashMap)(v).Delete(s)
}

func (vs ValueSlice) String() string {
	return fmt.Sprintf("Constants: %v", []Value(vs))
}

func (vs *ValueSlice) addConstant(v Value) int {
	*vs = append(*vs, v)

	return len(*vs) - 1
}

func (vs *ValueStack) Push(v Value) {
	*vs = append(*vs, v)
}

func (vs *ValueStack) Pop() Value {
	ret := (*vs)[len(*vs)-1]
	*vs = (*vs)[0 : len(*vs)-1]

	return ret
}

func (vs *ValueStack) Reset() {
	*vs = make(ValueStack, 0)
}
