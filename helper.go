package main

import "fmt"

func print(args ...interface{}) {
	fmt.Println(args...)
}

func printf(format string, args ...interface{}) {
	fmt.Printf(format + "\n", args ...)
}
