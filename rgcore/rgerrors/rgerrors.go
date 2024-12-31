package rgerrors

import (
	"errors"
	"fmt"
)

/*
#cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu/ -lm -llua5.3 -Wl,--allow-multiple-definition
#cgo CFLAGS: -I/usr/include/lua5.3/
#include "../lua/luaBridge.c"
*/
import "C"

type RGError struct {
	Err  error
	Code int
}

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

	RG_CORE_SYSTEM_CORRUPTED_ERROR_MESSAGE = "RGCoreSystemCorruptedError"
	UNDEFINED_ACT_FUNCTION_ERROR_MESSAGE   = "UndefinedActFunctionError"
	INVALID_ACTION_FORMAT_ERROR_MESSAGE    = "InvalidActionFormatError"
	INVALID_ACTION_TYPE_ERROR_MESSAGE      = "InvalidActionTypeError"
	INVALID_MOVE_ERROR_MESSAGE             = "InvalidMoveError"
	DISQUALIFIED_ERROR_MESSAGE             = "Disqualification"

	TIMEOUT_ERROR_MESSAGE = "TimeoutError"
)

var (
	ErrLuaYield     = errors.New("YieldError")
	ErrLuaErrRun    = errors.New("RunTimeError")
	ErrLuaErrSyntax = errors.New("SyntaxError")
	ErrLuaErrMem    = errors.New("ERRMEM")
	ErrLuaErrGcmm   = errors.New("ERRGCMM")
	ErrLuaErrErr    = errors.New("ERRERR")

	RG_CORE_SYSTEM_CORRUPTED_ERROR = wrapRGError(101)
	UNDEFINED_ACT_FUNCTION_ERROR   = wrapRGError(102)
	INVALID_ACTION_FORMAT_ERROR    = wrapRGError(103)
	INVALID_ACTION_TYPE_ERROR      = wrapRGError(104)
	INVALID_MOVE_ERROR             = wrapRGError(105)
	DISQUALIFIED_ERROR             = wrapRGError(106)

	TIMEOUT_ERROR = wrapRGError(int(C.CUSTOM_TIMEOUT_ERROR))
)

func GetLuaError(code int) error {
	switch code {
	case LUA_OK:
		return nil
	case LUA_YIELD:
		return ErrLuaYield
	case LUA_ERRRUN:
		return ErrLuaErrRun
	case LUA_ERRSYNTAX:
		return ErrLuaErrSyntax
	case LUA_ERRMEM:
		return ErrLuaErrMem
	case LUA_ERRGCMM:
		return ErrLuaErrGcmm
	case LUA_ERRERR:
		return ErrLuaErrErr
	default:
		return errors.New("GenericError")
	}
}

func (e *RGError) Error() string { return fmt.Sprintf("%v", e.Err) }
func (e *RGError) Unwrap() error { return e.Err }
func wrapRGError(code int) *RGError {
	var err error
	switch code {
	case 101:
		err = errors.New(RG_CORE_SYSTEM_CORRUPTED_ERROR_MESSAGE)
	case 102:
		err = errors.New(UNDEFINED_ACT_FUNCTION_ERROR_MESSAGE)
	case 103:
		err = errors.New(INVALID_ACTION_FORMAT_ERROR_MESSAGE)
	case 104:
		err = errors.New(INVALID_ACTION_TYPE_ERROR_MESSAGE)
	case 105:
		err = errors.New(INVALID_MOVE_ERROR_MESSAGE)
	case 106:
		err = errors.New(DISQUALIFIED_ERROR_MESSAGE)
	case int(C.CUSTOM_TIMEOUT_ERROR):
		err = errors.New(TIMEOUT_ERROR_MESSAGE)
	default:
		err = GetLuaError(code)
	}
	return &RGError{
		Err:  err,
		Code: code,
	}
}

func GetRGError(code int) *RGError {
	switch code {
	case 101:
		return RG_CORE_SYSTEM_CORRUPTED_ERROR
	case 102:
		return UNDEFINED_ACT_FUNCTION_ERROR
	case 103:
		return INVALID_ACTION_FORMAT_ERROR
	case 104:
		return INVALID_ACTION_TYPE_ERROR
	case 105:
		return INVALID_MOVE_ERROR
	case 106:
		return DISQUALIFIED_ERROR
	case int(C.CUSTOM_TIMEOUT_ERROR):
		return TIMEOUT_ERROR
	default:
		return &RGError{
			Err:  GetLuaError(code),
			Code: code,
		}
	}
}
