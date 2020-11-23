package debug

import (
	"fmt"
	"os"
)

func Log(format string, a ...interface{}) {
	if _, ok := os.LookupEnv("VERMIN_DEBUG"); ok {
		fmt.Printf(format+"\n", a...)
	}
}
