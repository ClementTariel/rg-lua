package rgcore

/*
#cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu/ -lm -llua5.3 -Wl,--allow-multiple-definition
#cgo CFLAGS: -I/usr/include/lua5.3/
#include "rg.c"
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/ClementTariel/rg-lua/rgcore/lua"
	"github.com/ClementTariel/rg-lua/rgcore/rgconst"
	"github.com/ClementTariel/rg-lua/rgcore/rgentities"
	"github.com/ClementTariel/rg-lua/rgcore/rgerrors"
	"github.com/ClementTariel/rg-lua/rgcore/rgutils"
)

type (
	CAction struct {
		actionType C.int
		x          C.int
		y          C.int
	}
)

func GetActionWithTimeout(pl unsafe.Pointer, bot rgentities.Bot) (rgentities.Action, *rgerrors.RGError) {
	cAction := CAction{
		actionType: (C.int)(rgconst.SUICIDE),
		x:          -1,
		y:          -1,
	}
	action := rgentities.Action{
		ActionType: rgconst.SUICIDE,
		X:          -1,
		Y:          -1,
	}
	errCode := int(C.GetActionWithTimeoutBridge(pl, unsafe.Pointer(&cAction), C.int(bot.Id), rgconst.BOT_ACTION_TIME_BUDGET))
	err := rgerrors.GetRGError(errCode)
	if errCode != 0 {
		fmt.Printf("Error after timed exection: %v\n", err)
		return action, err
	}
	if rgentities.ActionType(cAction.actionType) == rgconst.MOVE ||
		rgentities.ActionType(cAction.actionType) == rgconst.ATTACK ||
		rgentities.ActionType(cAction.actionType) == rgconst.GUARD ||
		rgentities.ActionType(cAction.actionType) == rgconst.SUICIDE {
		action.ActionType = rgentities.ActionType(cAction.actionType)
	} else {
		action.ActionType = rgconst.GUARD
		fmt.Printf("Error after timed exection: %v\n", rgerrors.INVALID_ACTION_TYPE_ERROR)
		return action, rgerrors.INVALID_ACTION_TYPE_ERROR
	}
	if int(cAction.x) >= 0 &&
		int(cAction.x) < rgconst.GRID_SIZE &&
		int(cAction.y) >= 0 &&
		int(cAction.y) < rgconst.GRID_SIZE {
		action.X = int(cAction.x)
		action.Y = int(cAction.y)
	} else {
		action.ActionType = rgconst.GUARD
		err = rgerrors.INVALID_MOVE_ERROR
	}
	if action.ActionType == rgconst.MOVE {
		if rgutils.WalkDist(bot.X, bot.Y, action.X, action.Y) != 1 {
			action.ActionType = rgconst.GUARD
			err = rgerrors.INVALID_MOVE_ERROR
		}
	} else if action.ActionType == rgconst.ATTACK {
		attackRange := rgutils.WalkDist(bot.X, bot.Y, action.X, action.Y)
		if attackRange < 1 || attackRange > rgconst.ATTACK_RANGE {
			action.ActionType = rgconst.GUARD
			err = rgerrors.INVALID_MOVE_ERROR
		}
	}
	if action.ActionType == rgconst.GUARD || action.ActionType == rgconst.SUICIDE {
		action.X = -1
		action.Y = -1
	}
	return action, err
}

func NewRGState() (unsafe.Pointer, error) {
	pl := lua.NewState()
	C.LoadRg(pl)
	err := lua.RunScript(pl, GetInitialisationScript(), "[Initialisation Script]")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	return pl, nil
}

func InitRG(pl unsafe.Pointer, script string, fileName string) (int, error) {
	var err error
	warningCount := 0
	timeoutBuffer := rgconst.BOT_ACTION_TIME_BUDGET * rgconst.WARNING_TOLERANCE
	var timeLeft int
	timeLeft, err = lua.RunScriptWithTimeout(pl, script, fileName, rgconst.BOT_INIT_TIME_BUDGET, timeoutBuffer)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return rgconst.WARNING_TOLERANCE + 1, err
	}
	for timeLeft < 0 && warningCount <= rgconst.WARNING_TOLERANCE {
		warningCount++
		timeLeft += rgconst.BOT_ACTION_TIME_BUDGET
	}
	err = lua.RunScript(pl, GetLoadActScript(), "[load act]")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return rgconst.WARNING_TOLERANCE + 1, err
	}
	return warningCount, err
}
