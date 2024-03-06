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
    vm.InteractiveMode = true

    for ;; {
        fmt.Print("> ")
        line, err := reader.ReadString(byte('\n'))
        if err != nil {
            return
        }

        if err := vm.Interpret(line[:len(line)-1]); err != nil {
            fmt.Println(err.Error())
        }
    }
}

func runFile(path string) {
    vm := vm.VirtualMachine{}
    code, err := os.ReadFile(path)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    if err := vm.Interpret(string(code)); err != nil {
        fmt.Println(err.Error())
    }
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
