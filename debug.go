package main

import (
	"fmt"
)

var VERBOSE = false

func VPrintf(format string, a ...any) (int, error) {
	if VERBOSE {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

func VPrintln(format string, a ...any) (int, error) {
	if VERBOSE {
		return fmt.Println(a...)
	}
	return 0, nil
}
