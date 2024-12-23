package rgcore

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

/*
#cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu/ -lm -llua5.3 -Wl,--allow-multiple-definition
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

var (
	LUA_YIELD_ERROR     = errors.New("YieldError")
	LUA_ERRRUN_ERROR    = errors.New("RunTimeError")
	LUA_ERRSYNTAX_ERROR = errors.New("SyntaxError")
	LUA_ERRMEM_ERROR    = errors.New("ERRMEM")
	LUA_ERRGCMM_ERROR   = errors.New("ERRGCMM")
	LUA_ERRERR_ERROR    = errors.New("ERRERR")
)

func GetStack(pl unsafe.Pointer, depth int, ptDebug unsafe.Pointer) C.int {
	return C.GetStackBridge(pl, C.int(depth), ptDebug)
}

func GetInfo(pl unsafe.Pointer, ptDebug unsafe.Pointer) {
	C.GetInfoBridge(pl, ptDebug)
}

func ToString(pl unsafe.Pointer, n int) string {
	return C.GoString(C.ToStringBridge(pl, C.int(n)))
}

func PopState(pl unsafe.Pointer, n int) {
	C.PopStateBridge(pl, C.int(n))
}

func GetLuaError(code int) error {
	switch code {
	case LUA_OK:
		return nil
	case LUA_YIELD:
		return LUA_YIELD_ERROR
	case LUA_ERRRUN:
		return LUA_ERRRUN_ERROR
	case LUA_ERRSYNTAX:
		return LUA_ERRSYNTAX_ERROR
	case LUA_ERRMEM:
		return LUA_ERRMEM_ERROR
	case LUA_ERRGCMM:
		return LUA_ERRGCMM_ERROR
	case LUA_ERRERR:
		return LUA_ERRERR_ERROR
	default:
		return errors.New("GenericError")
	}
}

func lineErrorMessage(pl unsafe.Pointer, script string, errorName string, fileName string) error {
	msg := ToString(pl, -1)
	parts := strings.Split(msg, ":")
	PopState(pl, 1)
	// TODO handle empty error messages
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
	return unsafe.Pointer(C.NewStateBridge())
}

func CloseState(pl unsafe.Pointer) {
	C.CloseBridge(pl)
}

func PushFunction(pl unsafe.Pointer, pfn unsafe.Pointer, name string) {
	C.PushCFunctionBridge(pl, pfn)
	C.SetGlobalBridge(pl, C.CString(name))
}

func RunScript(pl unsafe.Pointer, script string, fileName string) error {
	var res int
	VPrintf("Running %s\n", fileName)
	// Add "\n" as a quick fix to allow lazy error handling
	res = int(C.LoadStringBridge(pl, C.CString("\n"+script)))
	if res != LUA_OK {
		fmt.Println("Exit with error")
		return lineErrorMessage(pl, script, GetLuaError(int(res)).Error(), fileName)
	}
	res = int(C.PcallBridge(pl, 0, -1, 0))
	if res != LUA_OK {
		fmt.Println("Exit with error")
		// var ptDebug unsafe.Pointer
		// currentLine := 0
		// depth := 1
		// stack := int(GetStack(pl, depth, ptDebug))
		// if stack == 1 {
		// 	GetInfo(pl, ptDebug)
		// 	currentLine = int(C.DebugCurrentLine(ptDebug))
		// }
		// fmt.Printf("stack %v\n", stack)
		// var errormsg = "At line " + fmt.Sprintf("%v", currentLine) + "\n"
		// return errors.New(GetLuaError(int(res)).Error() + ":" + errormsg)
		return lineErrorMessage(pl, script, GetLuaError(int(res)).Error(), fileName)
	}
	return nil
}

func RunScriptWithTimeout(pl unsafe.Pointer, script string, fileName string, timeout int, timeoutBuffer int) (int, error) {
	startTime := time.Now()
	timeLeft := timeout
	var res int
	VPrintf("Running %s:\n%s\n", fileName, "\n"+script+"\n")
	// Add "\n" as a quick fix to allow lazy error handling
	res = int(C.LoadStringBridge(pl, C.CString("\n"+script+"\n")))
	if res != LUA_OK {
		timeLeft = timeout - int(time.Since(startTime).Milliseconds())
		return timeLeft, lineErrorMessage(pl, script, GetLuaError(int(res)).Error(), fileName)
	}
	timeLeft = timeout - int(time.Since(startTime).Milliseconds())
	if timeoutBuffer+timeLeft < 0 {
		return timeLeft, GetRGError(TIMEOUT_ERROR.Code)
	}
	res = int(C.PcallWithTimeoutBridge(pl, 0, -1, 0, C.int(timeLeft)))
	if errors.Is(GetRGError(int(res)), TIMEOUT_ERROR) {
		timeLeft = timeout - int(time.Since(startTime).Milliseconds())
		return timeLeft, GetRGError(int(res))
	}
	if res != LUA_OK {
		// var ptDebug unsafe.Pointer
		// currentLine := 0
		// depth := 1
		// stack := int(GetStack(pl, depth, ptDebug))
		// if stack == 1 {
		// 	GetInfo(pl, ptDebug)
		// 	currentLine = int(C.DebugCurrentLine(ptDebug))
		// }
		// fmt.Printf("stack %v\n", stack)
		// var errormsg = "At line " + fmt.Sprintf("%v", currentLine) + "\n"
		// return errors.New(GetLuaError(int(res)).Error() + ":" + errormsg)
		timeLeft = timeout - int(time.Since(startTime).Milliseconds())
		return timeLeft, lineErrorMessage(pl, script, GetLuaError(int(res)).Error(), fileName)
	}
	timeLeft = timeout - int(time.Since(startTime).Milliseconds())
	if timeoutBuffer+timeLeft < 0 {
		return timeLeft, GetRGError(TIMEOUT_ERROR.Code)
	}
	return timeLeft, nil
}
