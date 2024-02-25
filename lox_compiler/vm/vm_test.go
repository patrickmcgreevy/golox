package vm_test

import (
	"lox-compiler/vm"
	"testing"
)

func TestVMBinary(t *testing.T) {
    v := vm.VirtualMachine{}
    ret := v.Interpret("2+2;")

    if ret != vm.Interpret_OK {
        t.Fatalf("%s", ret)
    }
    ret = v.Interpret("2+2;")

    if ret != vm.Interpret_OK {
        t.Fatalf("%s", ret)
    }
}
