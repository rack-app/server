package debug

import (
	"fmt"
)

func Println(v ...interface{}) (n int, err error) {
	if inDebug {
		return fmt.Println(v...)
	}

	return 0, nil
}

func Printf(format string, a ...interface{}) (n int, err error) {
	if inDebug {
		return fmt.Printf(format, a...)
	}

	return 0, nil
}
