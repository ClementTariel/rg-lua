package rgconst

/*
#include "grid/grid.c"
#include "rgconst.c"
*/
import "C"

import "github.com/ClementTariel/rg-lua/rgcore/rgentities"

const (
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

	MOVE    = rgentities.ActionType(C.MOVE)
	ATTACK  = rgentities.ActionType(C.ATTACK)
	GUARD   = rgentities.ActionType(C.GUARD)
	SUICIDE = rgentities.ActionType(C.SUICIDE)

	NORMAL   = rgentities.LocationType(C.NORMAL)
	SPAWN    = rgentities.LocationType(C.SPAWN)
	OBSTACLE = rgentities.LocationType(C.OBSTACLE)
)

var (
	SPAWN_LEN    = int(C.SPAWN_LEN)
	ARENA_RADIUS = float64(int(C.ARENA_RADIUS))
	GRID_SIZE    = int(C.GRID_SIZE)

	GRID            = C.GRID
	SPAWN_LOCATIONS = C.SPAWN_LOCATIONS
)
