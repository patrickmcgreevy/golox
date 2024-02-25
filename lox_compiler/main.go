package main

import (
	"flag"
	"fmt"
	"lox-compiler/vm"
    "os"
    "bufio"
)

func repl() {
	reader := bufio.NewReader(os.Stdin)
    vm := vm.VirtualMachine{}

    for ;; {
        fmt.Print("> ")
        line, err := reader.ReadString(byte('\n'))
        if err != nil {
            return
        }

        fmt.Println(vm.Interpret(line[:len(line)-1]))
    }
}

func runFile(path string) {
    vm := vm.VirtualMachine{}
    code, err := os.ReadFile(path)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    fmt.Println(vm.Interpret(string(code)))
}

func usage() {
    fmt.Fprintln(os.Stderr, "usage: lox [path]")
}

func main() {
    flag.Parse()
	args := flag.Args()
    if len(args) == 0 {
        repl()
    } else if len(args) == 1 {
        runFile(args[0])
    } else {
        usage()
        os.Exit(64)
    }
}

//
// func main() {
//     vm := vm.VirtualMachine{}
// 	c := bytecode.NewChunk()
//     c.AddInst(bytecode.NewConstantInst(bytecode.Operand(c.AddConstant(1.2)), 0))
//     for i := 0; i < 100000; i += 1 {
//         c.AddInst(bytecode.NewNegateInst(0))
//     }
//     // c.AddInst(bytecode.NewConstantInst(bytecode.Operand(c.AddConstant(3.4)), 0))
//     // c.AddInst(bytecode.NewAddInst(0))
//     // c.AddInst(bytecode.NewConstantInst(bytecode.Operand(c.AddConstant(5.6)), 0))
//     // c.AddInst(bytecode.NewDivideInst(0))
//     // c.AddInst(bytecode.NewNegateInst(0))
// 	c.AddInst(bytecode.NewReturnInst(0))
//     fmt.Println(vm.Interpret(&c))
// }
