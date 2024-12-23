package rgcore

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
	CAction struct {
		actionType C.int
		x          C.int
		y          C.int
	}
	Location struct {
		X int
		Y int
	}
	LocationType int
	ActionType   int
	Action       struct {
		ActionType ActionType
		X          int
		Y          int
	}
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
	RGError struct {
		Err  error
		Code int
	}
)

const (
	PORT_PLAYER_1 = 1111
	PORT_PLAYER_2 = 2222

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

	RG_CORE_SYSTEM_CORRUPTED_ERROR_MESSAGE = "RGCoreSystemCorruptedError"
	UNDEFINED_ACT_FUNCTION_ERROR_MESSAGE   = "UndefinedActFunctionError"
	INVALID_ACTION_FORMAT_ERROR_MESSAGE    = "InvalidActionFormatError"
	INVALID_ACTION_TYPE_ERROR_MESSAGE      = "InvalidActionTypeError"
	INVALID_MOVE_ERROR_MESSAGE             = "InvalidMoveError"
	DISQUALIFIED_ERROR_MESSAGE             = "Disqualification"

	TIMEOUT_ERROR_MESSAGE = "TimeoutError"
)

var (
	GRID            = C.GRID
	SPAWN_LOCATIONS = C.SPAWN_LOCATIONS

	RG_CORE_SYSTEM_CORRUPTED_ERROR = wrapRGError(101)
	UNDEFINED_ACT_FUNCTION_ERROR   = wrapRGError(102)
	INVALID_ACTION_FORMAT_ERROR    = wrapRGError(103)
	INVALID_ACTION_TYPE_ERROR      = wrapRGError(104)
	INVALID_MOVE_ERROR             = wrapRGError(105)
	DISQUALIFIED_ERROR             = wrapRGError(106)

	TIMEOUT_ERROR = wrapRGError(int(C.CUSTOM_TIMEOUT_ERROR))
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

func GetActionWithTimeout(pl unsafe.Pointer, bot Bot) (Action, *RGError) {
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
	err := GetRGError(errCode)
	if errors.Unwrap(err) != nil {
		return Action{ActionType: SUICIDE, X: -1, Y: -1}, err
	}
	if action.ActionType == MOVE {
		if WalkDist(bot.X, bot.Y, action.X, action.Y) != 1 {
			action.ActionType = GUARD
			err = INVALID_MOVE_ERROR
		}
	} else if action.ActionType == ATTACK {
		attackRange := WalkDist(bot.X, bot.Y, action.X, action.Y)
		if attackRange < 1 || attackRange > ATTACK_RANGE {
			action.ActionType = GUARD
			err = INVALID_MOVE_ERROR
		}
	}
	if action.ActionType == GUARD || action.ActionType == SUICIDE {
		action.X = -1
		action.Y = -1
	}
	return action, err
}

func GetInitialisationScript() string {
	return `MOVE = ` + fmt.Sprintf("%d", MOVE) + `
ATTACK = ` + fmt.Sprintf("%d", ATTACK) + `
GUARD = ` + fmt.Sprintf("%d", GUARD) + `
SUICIDE = ` + fmt.Sprintf("%d", SUICIDE) + `
NORMAL = ` + fmt.Sprintf("%d", NORMAL) + `
SPAWN = ` + fmt.Sprintf("%d", SPAWN) + `
OBSTACLE = ` + fmt.Sprintf("%d", OBSTACLE) + `
rg.SETTINGS = {
	spawn_delay = ` + fmt.Sprintf("%d", SPAWN_DELAY) + `,
	spawn_count = ` + fmt.Sprintf("%d", SPAWN_COUNT) + `,
	robot_hp = ` + fmt.Sprintf("%d", MAX_HP) + `,
	attack_range = ` + fmt.Sprintf("%d", ATTACK_RANGE) + `,
	attack_damage = { ` +
		`min=` + fmt.Sprintf("%d", ATTACK_DAMAGE_MIN) + `, ` +
		`max=` + fmt.Sprintf("%d", ATTACK_DAMAGE_MAX) +
		` },
	suicide_damage = ` + fmt.Sprintf("%d", SUICIDE_DAMAGE) + `,
	collision_damage = ` + fmt.Sprintf("%d", COLLISION_DAMAGE) + `,
	max_turn = ` + fmt.Sprintf("%d", MAX_TURN) + `,
}
__RG_CORE_SYSTEM = { 
	GRID_SIZE = rg.GRID_SIZE,
	self = {},
	game = {},
}
`
}

func GetLoadActScript() string {
	return `__RG_CORE_SYSTEM.act = act`
}

func InitRG(pl unsafe.Pointer, script string, fileName string) (int, error) {
	var err error
	warningCount := 0
	C.LoadRg(pl)
	err = RunScript(pl, GetInitialisationScript(), "[Initialisation Script]")
	if err != nil {
		fmt.Printf("%v\n", err)
		return WARNING_TOLERANCE + 1, err
	}
	initTimeout := BOT_INIT_TIME_BUDGET
	actionTimeout := BOT_ACTION_TIME_BUDGET
	timeoutBuffer := actionTimeout * WARNING_TOLERANCE
	var timeLeft int
	timeLeft, err = RunScriptWithTimeout(pl, script, fileName, initTimeout, timeoutBuffer)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return WARNING_TOLERANCE + 1, err
	}
	for timeLeft < 0 && warningCount <= WARNING_TOLERANCE {
		warningCount++
		timeLeft += actionTimeout
	}
	err = RunScript(pl, GetLoadActScript(), "[load act]")
	if err != nil {
		fmt.Printf("%v\n", err)
		return WARNING_TOLERANCE + 1, err
	}
	return warningCount, err
}
