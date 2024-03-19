package vm_test

import (
	"lox-compiler/vm"
	"testing"
)

func test_interp(t *testing.T, s string) {
    vm := vm.VirtualMachine{}
    ret := vm.Interpret(s)
    if ret != nil {
        t.Fatalf("%s", ret.Error())
    }
}

func TestVMBinary(t *testing.T) {
    v := vm.VirtualMachine{}
    ret := v.Interpret("2+2;")

    if ret != nil {
        t.Fatalf("%s", ret.Error())
    }
    ret = v.Interpret("2+2;")

    if ret != nil {
        t.Fatalf("%s", ret.Error())
    }
}

func TestDeclareGlobal(t *testing.T) {
    vm := vm.VirtualMachine{}
    ret := vm.Interpret("var a = 2; print a; a = 3; print a;")
    if ret != nil {
        t.Fatalf("%s", ret.Error())
    }
}

func TestLocalVars(t *testing.T) {
    test_interp(t, "{var a = 0; var b = 1; a = 2; print a; print b;}")
}

func TestIfStmt(t *testing.T) {
    test_interp(t, "if (true)  {print \"true\"; print \"block\";} else {print \"false\"; print \"block\";}")
    test_interp(t, "if (false) {print \"true\"; print \"block\";} else {print \"false\"; print \"block\";}")
}
