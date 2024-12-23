package main

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/ClementTariel/rg-lua/rgcore"
)

var pl unsafe.Pointer
var running = false
var err error

func KillCurrentMatch() bool {
	rgcore.VPrintln("kill")
	killed := false
	if running {
		killed = true
		running = false
		// TODO: may need additional step to avoid some issue on the long term,
		// depending on how memory is cleared by lua, with potential memory leak or fragmentation
		rgcore.CloseState(pl)
	}
	rgcore.VPrintf("kill status: %v\n", killed)
	return killed
}

func InitNewMatch(name string, script string) (int, error) {
	KillCurrentMatch()
	pl = rgcore.NewState()
	rgcore.PushFunction(pl, rgcore.GetPrintInLuaFunctionPointer(), "print")
	var warningCount int
	warningCount, err = rgcore.InitRG(pl, script, name)
	if err != nil {
		fmt.Printf("%v\n", err)
		return warningCount, err
	}
	rgcore.VPrintln("[Successfully initialized]")
	running = true
	return warningCount, nil
}

func ResetGame(turn int) error {
	const resetScript = `__RG_CORE_SYSTEM.game.robots = rg.Robots()
for i = 0,%[1]d-1,1 do
	__RG_CORE_SYSTEM.game.robots[i] = {}
end
__RG_CORE_SYSTEM.game.turn = %[2]d
`
	return rgcore.RunScript(pl, fmt.Sprintf(resetScript, rgcore.GRID_SIZE, turn), "[reset game data]")
}

func LoadGameBot(bot rgcore.Bot) error {
	botId := "nil"
	if (bot.Id) > 0 {
		botId = fmt.Sprintf("%d", bot.Id)
	}
	const loadScript = `__RG_CORE_SYSTEM.game.robots[%[1]d][%[2]d] = {
	location = rg.Loc(%[1]d, %[2]d),
	hp = %[3]d,
	player_id = %[4]d,
	id = %[5]s,
}
`
	botDescription := fmt.Sprintf("bot %s", botId)
	if botId == "nil" {
		botDescription = "enemy bot"
	}
	return rgcore.RunScript(pl,
		fmt.Sprintf(loadScript, bot.X, bot.Y, bot.Hp, bot.PlayerId, botId),
		fmt.Sprintf("[loading game data - %s]", botDescription))
}

func LoadSelf(bot rgcore.Bot) error {
	const loadScript = `if __RG_CORE_SYSTEM.self[%[5]d] == nil then
	__RG_CORE_SYSTEM.self[%[5]d] = {}
end
__RG_CORE_SYSTEM.self[%[5]d].id = %[5]d
__RG_CORE_SYSTEM.self[%[5]d].location = rg.Loc(%[1]d, %[2]d)
__RG_CORE_SYSTEM.self[%[5]d].hp = %[3]d
__RG_CORE_SYSTEM.self[%[5]d].player_id = %[4]d
__RG_CORE_SYSTEM.self[%[5]d].id = %[5]d
`
	return rgcore.RunScript(pl,
		fmt.Sprintf(loadScript, bot.X, bot.Y, bot.Hp, bot.PlayerId, bot.Id),
		fmt.Sprintf("[loading self data - bot %d]", bot.Id))
}

func PlayTurn(turn int, allies []rgcore.Bot, enemies []rgcore.Bot, warningCount int) (map[int]rgcore.Action, int) {
	err := ResetGame(turn)
	if err != nil {
		fmt.Printf("error when reseting game: %v\n", err)
		warningCount = rgcore.WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
	}
	actions := map[int]rgcore.Action{}
	if !(warningCount > rgcore.WARNING_TOLERANCE) {
		for _, bot := range allies {
			err = LoadGameBot(bot)
			if err != nil {
				fmt.Printf("error when loading game bot %v: %v\n", bot, err)
				warningCount = rgcore.WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
				break
			}
		}
		if warningCount <= rgcore.WARNING_TOLERANCE {
			for _, bot := range enemies {
				err = LoadGameBot(bot)
				if err != nil {
					fmt.Printf("error when loading game bot %v: %v\n", bot, err)
					warningCount = rgcore.WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
					break
				}
			}
		}
	}
	for _, bot := range allies {
		actions[bot.Id] = rgcore.Action{
			ActionType: rgcore.SUICIDE,
			X:          -1,
			Y:          -1,
		}
		if !(warningCount > rgcore.WARNING_TOLERANCE) {
			err = LoadSelf(bot)
			if err != nil {
				fmt.Printf("error when loading self (bot %v) %v\n", bot, err)
				warningCount = rgcore.WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
				continue
			}
			action, err := rgcore.GetActionWithTimeout(pl, bot)

			rgcore.VPrintf("bot %d (%v) act (%d,%d,%d), %v\n", bot.Id, bot, action.ActionType, action.X, action.Y, err)
			if errors.Is(err, rgcore.INVALID_MOVE_ERROR) {
				action.ActionType = rgcore.GUARD
				action.X = -1
				action.Y = -1
			} else if errors.Is(err, rgcore.TIMEOUT_ERROR) {
				warningCount++
				action.ActionType = rgcore.GUARD
				action.X = -1
				action.Y = -1
			} else if errors.Unwrap(err) != nil {
				fmt.Printf("disqualified because of %v\n", err)
				warningCount = rgcore.WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
			}
			if warningCount > rgcore.WARNING_TOLERANCE {
				continue
			}
			actions[bot.Id] = action
		}
	}
	return actions, warningCount
}
