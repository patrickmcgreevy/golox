// Code generated by "stringer -type=OpCode"; DO NOT EDIT.

package bytecode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpAdd-0]
	_ = x[OpAnd-1]
	_ = x[OpAssign-2]
	_ = x[OpConstant-3]
	_ = x[OpDeclareGlobal-4]
	_ = x[OpDivide-5]
	_ = x[OpEqualEqual-6]
	_ = x[OpGreater-7]
	_ = x[OpGreaterEqual-8]
	_ = x[OpLess-9]
	_ = x[OpLessEqual-10]
	_ = x[OpLookup-11]
	_ = x[OpMultiply-12]
	_ = x[OpNegate-13]
	_ = x[OpNotEqual-14]
	_ = x[OpOr-15]
	_ = x[OpPrint-16]
	_ = x[OpReturn-17]
	_ = x[OpSubtract-18]
	_ = x[OpPop-19]
}

const _OpCode_name = "OpAddOpAndOpAssignOpConstantOpDeclareGlobalOpDivideOpEqualEqualOpGreaterOpGreaterEqualOpLessOpLessEqualOpLookupOpMultiplyOpNegateOpNotEqualOpOrOpPrintOpReturnOpSubtractOpPop"

var _OpCode_index = [...]uint8{0, 5, 10, 18, 28, 43, 51, 63, 72, 86, 92, 103, 111, 121, 129, 139, 143, 150, 158, 168, 173}

func (i OpCode) String() string {
	if i >= OpCode(len(_OpCode_index)-1) {
		return "OpCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OpCode_name[_OpCode_index[i]:_OpCode_index[i+1]]
}
