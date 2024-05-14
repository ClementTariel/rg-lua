package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

/*
#cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu/ -lm -llua5.3
#cgo CFLAGS: -I/usr/include/lua5.3/
#include "luaBridge.h"
#include "luaBridge.c"
#include <stdio.h>
*/
import "C"

const (
	LUA_OK            = 0
	LUA_YIELD         = 1
	LUA_ERRRUN        = 2
	LUA_ERRSYNTAX     = 3
	LUA_ERRMEM        = 4
	LUA_ERRGCMM       = 5
	LUA_ERRERR        = 6
	LUA_MULTRET       = -1
	LUA_GLOBALSINDEX  = -10002
	LUA_REGISTRYINDEX = -10000
	LUA_ENVIRONINDEX  = -10001
)

func GetStack(pl unsafe.Pointer, depth int, ptDebug unsafe.Pointer) C.int {
	return C.getStackBridge(pl, C.int(depth), ptDebug)
}

func GetInfo(pl unsafe.Pointer, ptDebug unsafe.Pointer) {
	C.getInfoBridge(pl, ptDebug)
}

func ToString(pl unsafe.Pointer, n int) string {
	return C.GoString(C.toStringBridge(pl, C.int(n)))
}

func PopState(pl unsafe.Pointer, n int) {
	C.popStateBridge(pl, C.int(n))
}

func ErrorName(code int) string {
	switch code {
	case LUA_OK:
		return ""
	case LUA_YIELD:
		return "YieldError"
	case LUA_ERRRUN:
		return "RunTimeError"
	case LUA_ERRSYNTAX:
		return "SyntaxError"
	case LUA_ERRMEM:
		return "ERRMEM"
	case LUA_ERRGCMM:
		return "ERRMEM"
	case LUA_ERRERR:
		return "ERRERR"
	default:
		return "DefaultError"
	}
}

func lineErrorMessage(pl unsafe.Pointer, script string, errorName string, fileName string) error {
	msg := ToString(pl, -1)
	parts := strings.Split(msg, ":")
	PopState(pl, 1)
	currentLine, err := strconv.ParseInt(parts[1], 10, 64)
	lineContent := "Line not found"
	if err != nil {
		currentLine = 0
	} else {
		// added "\n" as a quick fix to allow lazy error handling,
		// causing a shift of the line count
		currentLine -= 1
		lineContent = strings.Split(script, "\n")[currentLine-1]
	}
	return errors.New("\n" +
		fmt.Sprintf("%s\n", errorName) +
		fmt.Sprintf("%s: At line %d:\n", fileName, currentLine) +
		fmt.Sprintf("%d| %s\n", currentLine, lineContent) +
		fmt.Sprintf("\t%s\n", strings.Join(parts[2:], ":")))
}

func NewState() unsafe.Pointer {
	return unsafe.Pointer(C.newStateBridgeBridge())
}

func CloseState(pl unsafe.Pointer) {
	C.closeBridge(pl)
}

func PushFunction(pl unsafe.Pointer, pfn unsafe.Pointer, name string) {
	C.pushCFunctionBridge(pl, pfn)
	C.setGlobalBridge(pl, C.CString(name))
}

func RunScript(pl unsafe.Pointer, script string, fileName string) error {
	var res int
	fmt.Printf("Running %s\n", fileName)
	// Add "\n" as a quick fix to allow lazy error handling
	res = int(C.loadStringBridge(pl, C.CString("\n"+script)))
	if res != LUA_OK {
		fmt.Println("Exit with error")
		return lineErrorMessage(pl, script, ErrorName(int(res)), fileName)
	}
	res = int(C.pcallBridge(pl, 0, -1, 0))
	if res != LUA_OK {
		fmt.Println("Exit with error")
		// var ptDebug unsafe.Pointer
		// currentLine := 0
		// depth := 1
		// stack := int(GetStack(pl, depth, ptDebug))
		// if stack == 1 {
		// 	GetInfo(pl, ptDebug)
		// 	currentLine = int(C.debugCurrentLine(ptDebug))
		// }
		// fmt.Printf("stack %v\n", stack)
		// var errormsg = "At line " + fmt.Sprintf("%v", currentLine) + "\n"
		// return errors.New(ErrorName(res) + ":" + errormsg)
		return lineErrorMessage(pl, script, ErrorName(int(res)), fileName)
	}
	return nil
}
