package rgcore

import (
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu/ -lm -llua5.3 -Wl,--allow-multiple-definition
#cgo CFLAGS: -I/usr/include/lua5.3/
#include <stdio.h>
#include "debug/custom_string.c"
#include "debug/string_array.c"
#include "debug/debug.c"
*/
import "C"

var VERBOSE = false

func VPrintf(format string, a ...any) (int, error) {
	if VERBOSE {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

func VPrintln(a ...any) (int, error) {
	if VERBOSE {
		return fmt.Println(a...)
	}
	return 0, nil
}

func SetPrintMemoryBudget(budget int) {
	C.PRINT_MEMORY_BUDGET = C.size_t(budget)
}

func GetPrintInLuaFunctionPointer() unsafe.Pointer {
	return unsafe.Pointer(C.printInLua)
}
