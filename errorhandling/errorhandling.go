package errorhandling

import (
	"fmt"
)

func Report(line int, where, message string) {
	fmt.Printf("[line %d] Error %s: '%s'\n", line, where, message)
}

func RuntimeError(err error) {
	fmt.Println(err.Error())
}
