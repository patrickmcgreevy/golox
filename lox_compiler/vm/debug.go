// +build debug

package vm

import (
	"fmt"
)

func debug(format string, args ...any) {
    
	fmt.Println(fmt.Sprintf(format, args...))
}
