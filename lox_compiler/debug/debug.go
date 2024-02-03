// +build debug

package debug

import (
	"fmt"
)

func Printf(format string, args ...any) {
	fmt.Println(fmt.Sprintf(format, args...))
}
