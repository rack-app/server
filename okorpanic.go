package main

import "fmt"

func OkOrPanic(fn func() []error) {
	if errs := fn(); len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}
}
