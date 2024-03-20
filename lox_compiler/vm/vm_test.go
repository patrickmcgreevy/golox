package vm_test

import (
	"bufio"
	"io"
	"lox-compiler/vm"
	"os"
	"testing"
)

func test_interp(t *testing.T, s string) {
	vm := vm.VirtualMachine{}
	ret := vm.Interpret(s)
	if ret != nil {
		t.Fatalf("%s", ret.Error())
	}
}

func test_interp_output(t *testing.T, input, output string) {
	vm := vm.VirtualMachine{}
	s := os.Stdout
	r, w, _ := os.Pipe()
	reader := bufio.NewReader(r)

	os.Stdout = w
	defer func() { os.Stdout = s }()

	ret := vm.Interpret(input)
	if ret != nil {
		t.Fatalf("%s", ret.Error())
	}

	str, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return
		}
		t.Fatalf("fail: %s", err.Error())
	}
    if str != output {
        t.Fatalf("expected: '%s'\ngot: '%s'", output, str)
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

func TestIfOutput(t *testing.T) {
    test_interp_output(t, "if (true) {print \"yes\";} else { print \"no\";}", "yes\n")
    test_interp_output(t, "if (false) {print \"yes\";} else { print \"no\";}", "no\n")
}
