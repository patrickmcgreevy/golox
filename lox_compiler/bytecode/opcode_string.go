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
	_ = x[OpGlobalLookup-7]
	_ = x[OpGreater-8]
	_ = x[OpGreaterEqual-9]
	_ = x[OpLess-10]
	_ = x[OpLessEqual-11]
	_ = x[OpLocalAssign-12]
	_ = x[OpLocalLookup-13]
	_ = x[OpMultiply-14]
	_ = x[OpNegate-15]
	_ = x[OpNotEqual-16]
	_ = x[OpOr-17]
	_ = x[OpPop-18]
	_ = x[OpPrint-19]
	_ = x[OpReturn-20]
	_ = x[OpSubtract-21]
}

const _OpCode_name = "OpAddOpAndOpAssignOpConstantOpDeclareGlobalOpDivideOpEqualEqualOpGlobalLookupOpGreaterOpGreaterEqualOpLessOpLessEqualOpLocalAssignOpLocalLookupOpMultiplyOpNegateOpNotEqualOpOrOpPopOpPrintOpReturnOpSubtract"

var _OpCode_index = [...]uint8{0, 5, 10, 18, 28, 43, 51, 63, 77, 86, 100, 106, 117, 130, 143, 153, 161, 171, 175, 180, 187, 195, 205}

func (i OpCode) String() string {
	if i >= OpCode(len(_OpCode_index)-1) {
		return "OpCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OpCode_name[_OpCode_index[i]:_OpCode_index[i+1]]
}
