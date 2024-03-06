// Code generated by "stringer -type=OpCode"; DO NOT EDIT.

package bytecode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpReturn-0]
	_ = x[OpConstant-1]
	_ = x[OpNegate-2]
	_ = x[OpAdd-3]
	_ = x[OpSubtract-4]
	_ = x[OpMultiply-5]
	_ = x[OpDivide-6]
	_ = x[OpPrint-7]
	_ = x[OpOr-8]
	_ = x[OpAnd-9]
	_ = x[OpLess-10]
	_ = x[OpGreater-11]
	_ = x[OpLessEqual-12]
	_ = x[OpGreaterEqual-13]
	_ = x[OpEqualEqual-14]
	_ = x[OpNotEqual-15]
	_ = x[OpDeclare-16]
	_ = x[OpAssign-17]
	_ = x[OpLookup-18]
}

const _OpCode_name = "OpReturnOpConstantOpNegateOpAddOpSubtractOpMultiplyOpDivideOpPrintOpOrOpAndOpLessOpGreaterOpLessEqualOpGreaterEqualOpEqualEqualOpNotEqualOpDeclareOpAssignOpLookup"

var _OpCode_index = [...]uint8{0, 8, 18, 26, 31, 41, 51, 59, 66, 70, 75, 81, 90, 101, 115, 127, 137, 146, 154, 162}

func (i OpCode) String() string {
	if i >= OpCode(len(_OpCode_index)-1) {
		return "OpCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OpCode_name[_OpCode_index[i]:_OpCode_index[i+1]]
}
