package bytecode_test

import (
	"fmt"
	"lox-compiler/bytecode"
	"testing"
)

func BenchmarkInserts(b *testing.B) {
	for i := 0; i < b.N; i++ {
        m := bytecode.LoxMap(bytecode.NewLinearProbingHashMap())
		for j := 0; j < 1000000; j++ {
			m.Insert(bytecode.LoxString(fmt.Sprint(j)), bytecode.LoxInt(j))
		}
	}
}

func BenchmarkDeletions(b *testing.B) {
	for i := 0; i < b.N; i++ {
        m := bytecode.LoxMap(bytecode.NewLinearProbingHashMap())
		for j := 0; j < 1000000; j++ {
			m.Insert(bytecode.LoxString(fmt.Sprint(j)), bytecode.LoxInt(j))
            if j % 2 == 0 {
                m.Delete(bytecode.LoxString(fmt.Sprint(j)))
            } else {
                m.Get(bytecode.LoxString(fmt.Sprint(j)))
            }
		}
	}
}

