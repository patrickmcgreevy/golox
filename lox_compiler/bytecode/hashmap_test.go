package bytecode_test

import (
	"fmt"
    "lox-compiler/bytecode"
	"testing"
)

func TestInsertElem(t *testing.T) {
    m := bytecode.NewLinearProbingHashMap()
    fmt.Println(m)
    m.Insert("a", bytecode.LoxInt(1))
    fmt.Println(m)
    m.Insert("b", bytecode.LoxInt(2))
    fmt.Println(m)
}

func TestGetInvalid(t *testing.T) {
    m := bytecode.NewLinearProbingHashMap()
    _, err := m.Get("a")
    if err == nil {
        t.FailNow()
    }

    fmt.Println(err)
}

func TestGetValid(t *testing.T) {
    m := bytecode.NewLinearProbingHashMap()
    m.Insert("a", bytecode.LoxString("asdf"))
    m.Insert("b", bytecode.LoxString("1234"))
    val, err := m.Get("a")
    if err != nil {
        t.Fatalf(err.Error())
    }

    fmt.Println(val)
    val, err = m.Get("b")
    if err != nil {
        fmt.Println(m)
        fmt.Println(val)
        t.Fatalf(err.Error())
    }

    fmt.Println(val)
}

func TestRehash(t *testing.T) {
    m := bytecode.NewLinearProbingHashMap()
    for i := 0; i < 1000000; i++ {
        m.Insert(bytecode.LoxString(fmt.Sprint(i)), bytecode.LoxInt(i))
        v, err := m.Get(bytecode.LoxString(fmt.Sprint(i)))
        if err != nil {
            t.Fatalf("%s", err)
        }
        if v != bytecode.LoxInt(i) {
            t.Fatalf("%v != %v", v, i)
        }
    }
}

func TestDelElem(t *testing.T) {
    m := bytecode.NewLinearProbingHashMap()
    m.Insert("a", bytecode.LoxInt(1))
    fmt.Println(m)
    // m.Insert("b", bytecode.LoxInt(2))
    m.Delete("a")
    fmt.Println(m)
    _, err := m.Get("a")
    if err == nil {
        t.Fatalf("%s", m)
    }
}

