package main

/*
#include "rg.c"
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"unsafe"
)

type (
	ActionType C.int
	Action     struct {
		actionType ActionType
		x          C.int
		y          C.int
	}
	Bot struct {
		x        int
		y        int
		hp       int
		id       int
		playerId int
	}
	RGError struct {
		Err  error
		Code int
	}
)

const (
	MOVE ActionType = iota
	ATTACK
	GUARD
	SUICIDE

	GRID_SIZE = 19

	SPAWN_DELAY       = 10
	SPAWN_COUNT       = 5
	MAX_HP            = 50
	ATTACK_RANGE      = 1
	SUICIDE_RANGE     = 1
	ATTACK_DAMAGE_MIN = 8
	ATTACK_DAMAGE_MAX = 10
	COLLISION_DAMAGE  = 5
	SUICIDE_DAMAGE    = 15
	MAX_TURN          = 100

	WARNING_TOLERANCE = 3
	BOT_TIME_BUDGET   = 10

	RG_CORE_SYSTEM_CORRUPTED_ERROR_MESSAGE = "RGCoreSystemCorruptedError"
	UNDEFINED_ACT_FUNCTION_ERROR_MESSAGE   = "UndefinedActFunctionError"
	INVALID_ACTION_FORMAT_ERROR_MESSAGE    = "InvalidActionFormatError"
	INVALID_ACTION_TYPE_ERROR_MESSAGE      = "InvalidActionTypeError"
	INVALID_MOVE_ERROR_MESSAGE             = "InvalidMoveError"
	TIMEOUT_ERROR_MESSAGE                  = "TimeoutError"
)

var (
	RG_CORE_SYSTEM_CORRUPTED_ERROR = NewRGError(101)
	UNDEFINED_ACT_FUNCTION_ERROR   = NewRGError(102)
	INVALID_ACTION_FORMAT_ERROR    = NewRGError(103)
	INVALID_ACTION_TYPE_ERROR      = NewRGError(104)
	INVALID_MOVE_ERROR             = NewRGError(105)
	TIMEOUT_ERROR                  = NewRGError(106)
)

func (e *RGError) Error() string { return e.Err.Error() }
func (e *RGError) Unwrap() error { return e.Err }
func NewRGError(code int) *RGError {
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
		err = errors.New(TIMEOUT_ERROR_MESSAGE)
	default:
		err = errors.New(ErrorName(code))
	}
	return &RGError{
		Err:  err,
		Code: code,
	}
}

func Abs(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}

func WalkDist(x1, y1, x2, y2 int) int {
	return int(Abs(x1-x2) + Abs(y1-y2))
}

func Dist(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)))
}

func GetActionWithTimeout(pl unsafe.Pointer, bot Bot) (Action, error) {
	action := Action{
		actionType: SUICIDE,
		x:          -1,
		y:          -1,
	}
	var err error = nil
	errCode := int(C.getActionWithTimeoutBridge(pl, unsafe.Pointer(&action), BOT_TIME_BUDGET))
	if errCode == 0 &&
		action.actionType != MOVE &&
		action.actionType != ATTACK &&
		action.actionType != GUARD &&
		action.actionType != SUICIDE {
		errCode = 104
	}
	if errCode == 0 &&
		action.x < 0 ||
		action.x >= GRID_SIZE ||
		action.y < 0 ||
		action.y >= GRID_SIZE {
		action.x = -1
		action.y = -1
	}
	switch errCode {
	case 0:
		break
	case 101:
		err = RG_CORE_SYSTEM_CORRUPTED_ERROR
	case 102:
		err = UNDEFINED_ACT_FUNCTION_ERROR
	case 103:
		err = INVALID_ACTION_FORMAT_ERROR
	case 104:
		err = INVALID_ACTION_TYPE_ERROR
	case 105:
		err = INVALID_MOVE_ERROR
	case 106:
		err = TIMEOUT_ERROR
	default:
		err = errors.New(ErrorName(errCode))
	}
	if err != nil {
		return Action{actionType: SUICIDE, x: -1, y: -1}, err
	}
	if action.actionType == MOVE {
		// TODO: replace out of bound moves with GUARD too
		if WalkDist(bot.x, bot.y, int(action.x), int(action.y)) != 1 {
			action.actionType = GUARD
			err = INVALID_MOVE_ERROR
		}
	}
	if action.actionType == ATTACK {
		attackRange := WalkDist(bot.x, bot.y, int(action.x), int(action.y))
		if attackRange < 1 || attackRange > ATTACK_RANGE {
			action.actionType = GUARD
			err = INVALID_MOVE_ERROR
		}
	}
	if action.actionType == GUARD || action.actionType == SUICIDE {
		action.x = -1
		action.y = -1
	}
	return action, err
}

func GetInitialisationScript() string {
	c := fmt.Sprintf("%d", (GRID_SIZE-1)/2)
	return `MOVE = 0
ATTACK = 1
GUARD = 2
SUICIDE = 3
rg = {
	CENTER_POINT = { x=` + c + `, y=` + c + `},
	GRID_SIZE = ` + fmt.Sprintf("%d", GRID_SIZE) + `,
	settings = {
		spawn_delay = ` + fmt.Sprintf("%d", SPAWN_DELAY) + `,
		spawn_count = ` + fmt.Sprintf("%d", SPAWN_COUNT) + `,
		robot_hp = ` + fmt.Sprintf("%d", MAX_HP) + `,
		attack_range = ` + fmt.Sprintf("%d", ATTACK_RANGE) + `,
		attack_damage = {
			min=` + fmt.Sprintf("%d", ATTACK_DAMAGE_MIN) + `,
			max=` + fmt.Sprintf("%d", ATTACK_DAMAGE_MAX) + `
		},
		suicide_range = ` + fmt.Sprintf("%d", SUICIDE_RANGE) + `,
		suicide_damage = ` + fmt.Sprintf("%d", SUICIDE_DAMAGE) + `,
		collision_damage = ` + fmt.Sprintf("%d", COLLISION_DAMAGE) + `,
		max_turn = ` + fmt.Sprintf("%d", MAX_TURN) + `
	}
}
`
}

func GetLoadActScript() string {
	return `__RG_CORE_SYSTEM = {act = act}`
}

func InitRG(pl unsafe.Pointer, script string, fileName string) error {
	var err error
	err = RunScript(pl, GetInitialisationScript(), "[Initialisation Script]")
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	err = RunScript(pl, script, fileName)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	err = RunScript(pl, GetLoadActScript(), "[load act]")
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}
	return err
}
