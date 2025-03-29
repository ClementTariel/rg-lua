package rgcore

/*
#cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu/ -lm -llua5.3 -Wl,--allow-multiple-definition
#cgo CFLAGS: -I/usr/include/lua5.3/
#include "rg.c"
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/ClementTariel/rg-lua/rgcore/lua"
	"github.com/ClementTariel/rg-lua/rgcore/rgerrors"
)

type (
	CAction struct {
		actionType C.int
		x          C.int
		y          C.int
	}
	ActionType int
	Action     struct {
		ActionType ActionType
		X          int
		Y          int
	}

	Location struct {
		X int
		Y int
	}
	LocationType int

	Bot struct {
		X        int
		Y        int
		Hp       int
		Id       int
		PlayerId int
	}
	BotState struct {
		Bot    Bot
		Action Action
	}
)

const (
	MOVE    = ActionType(C.MOVE)
	ATTACK  = ActionType(C.ATTACK)
	GUARD   = ActionType(C.GUARD)
	SUICIDE = ActionType(C.SUICIDE)

	NORMAL    = LocationType(C.NORMAL)
	SPAWN     = LocationType(C.SPAWN)
	OBSTACLE  = LocationType(C.OBSTACLE)
	SPAWN_LEN = int(C.SPAWN_LEN)

	ARENA_RADIUS = float64(C.ARENA_RADIUS)
	GRID_SIZE    = int(C.GRID_SIZE)

	BLUE_ID = 1
	RED_ID  = 2

	SPAWN_DELAY       = 10
	SPAWN_COUNT       = 5
	MAX_HP            = 50
	ATTACK_RANGE      = 1
	ATTACK_DAMAGE_MIN = 8
	ATTACK_DAMAGE_MAX = 10
	COLLISION_DAMAGE  = 5
	SUICIDE_DAMAGE    = 15
	MAX_TURN          = 100

	WARNING_TOLERANCE      = 3
	BOT_INIT_TIME_BUDGET   = 1000
	BOT_ACTION_TIME_BUDGET = 10
)

var (
	GRID            = C.GRID
	SPAWN_LOCATIONS = C.SPAWN_LOCATIONS
)

func GetLocationType(x int, y int) LocationType {
	return LocationType(GRID[C.int(x)][C.int(y)])
}

func GetSpawnLocation(i int) (Location, error) {
	if i < 0 || i >= len(SPAWN_LOCATIONS) {
		return Location{X: -1, Y: -1}, errors.New("spawn index out of range in spawn generation")
	}
	return Location{X: int(SPAWN_LOCATIONS[C.int(i)].X), Y: int(SPAWN_LOCATIONS[C.int(i)].Y)}, nil
}

func GetActionWithTimeout(pl unsafe.Pointer, bot Bot) (Action, *rgerrors.RGError) {
	cAction := CAction{
		actionType: (C.int)(SUICIDE),
		x:          -1,
		y:          -1,
	}
	action := Action{
		ActionType: SUICIDE,
		X:          -1,
		Y:          -1,
	}
	errCode := int(C.GetActionWithTimeoutBridge(pl, unsafe.Pointer(&cAction), C.int(bot.Id), BOT_ACTION_TIME_BUDGET))
	if errCode == 0 &&
		ActionType(cAction.actionType) == MOVE ||
		ActionType(cAction.actionType) == ATTACK ||
		ActionType(cAction.actionType) == GUARD ||
		ActionType(cAction.actionType) == SUICIDE {
		action.ActionType = ActionType(cAction.actionType)
	} else {
		errCode = 104
	}
	if errCode == 0 &&
		int(cAction.x) >= 0 &&
		int(cAction.x) < GRID_SIZE &&
		int(cAction.y) >= 0 &&
		int(cAction.y) < GRID_SIZE {
		action.X = int(cAction.x)
		action.Y = int(cAction.y)
	} else {
		action.X = -1
		action.Y = -1
	}
	err := rgerrors.GetRGError(errCode)
	if errors.Unwrap(err) != nil {
		return Action{ActionType: SUICIDE, X: -1, Y: -1}, err
	}
	if action.ActionType == MOVE {
		if WalkDist(bot.X, bot.Y, action.X, action.Y) != 1 {
			action.ActionType = GUARD
			err = rgerrors.INVALID_MOVE_ERROR
		}
	} else if action.ActionType == ATTACK {
		attackRange := WalkDist(bot.X, bot.Y, action.X, action.Y)
		if attackRange < 1 || attackRange > ATTACK_RANGE {
			action.ActionType = GUARD
			err = rgerrors.INVALID_MOVE_ERROR
		}
	}
	if action.ActionType == GUARD || action.ActionType == SUICIDE {
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
		fmt.Printf("%v\n", err)
		return nil, err
	}
	return pl, nil
}

func InitRG(pl unsafe.Pointer, script string, fileName string) (int, error) {
	var err error
	warningCount := 0
	timeoutBuffer := BOT_ACTION_TIME_BUDGET * WARNING_TOLERANCE
	var timeLeft int
	timeLeft, err = lua.RunScriptWithTimeout(pl, script, fileName, BOT_INIT_TIME_BUDGET, timeoutBuffer)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return WARNING_TOLERANCE + 1, err
	}
	for timeLeft < 0 && warningCount <= WARNING_TOLERANCE {
		warningCount++
		timeLeft += BOT_ACTION_TIME_BUDGET
	}
	err = lua.RunScript(pl, GetLoadActScript(), "[load act]")
	if err != nil {
		fmt.Printf("%v\n", err)
		return WARNING_TOLERANCE + 1, err
	}
	return warningCount, err
}
