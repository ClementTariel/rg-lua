package main

import (
	"errors"
	"fmt"
	"unsafe"
)

func ResetGame(pl unsafe.Pointer, turn int) error {
	const resetScript = `__RG_CORE_SYSTEM.game = {}
for i = 0,%[1]d-1,1 do
	__RG_CORE_SYSTEM.game[i] = {}
end
__RG_CORE_SYSTEM.game.turn = %[2]d
`
	return RunScript(pl, fmt.Sprintf(resetScript, GRID_SIZE, turn), "[reset game data]")
}

func LoadGameBot(pl unsafe.Pointer, bot Bot) error {
	botId := "nil"
	if (bot.id) > 0 {
		botId = fmt.Sprintf("%d", bot.id)
	}
	const loadScript = `__RG_CORE_SYSTEM.game[%[1]d][%[2]d] = {}
__RG_CORE_SYSTEM.game[%[1]d][%[2]d].x = %[1]d
__RG_CORE_SYSTEM.game[%[1]d][%[2]d].y = %[2]d
__RG_CORE_SYSTEM.game[%[1]d][%[2]d].hp = %[3]d
__RG_CORE_SYSTEM.game[%[1]d][%[2]d].player_id = %[4]d
__RG_CORE_SYSTEM.game[%[1]d][%[2]d].id = %[5]s
`
	botDescription := fmt.Sprintf("bot %s", botId)
	if botId == "nil" {
		botDescription = "enemy bot"
	}
	return RunScript(pl,
		fmt.Sprintf(loadScript, bot.x, bot.y, bot.hp, bot.playerId, botId),
		fmt.Sprintf("[loading game data - %s]", botDescription))
}

func LoadSelf(pl unsafe.Pointer, bot Bot) error {
	const loadScript = `__RG_CORE_SYSTEM.self = {}
__RG_CORE_SYSTEM.self.id = %[1]d
__RG_CORE_SYSTEM.self.hp = %[2]d
__RG_CORE_SYSTEM.self.x = %[3]d
__RG_CORE_SYSTEM.self.y = %[4]d
`
	return RunScript(pl,
		fmt.Sprintf(loadScript, bot.id, bot.hp, bot.x, bot.y),
		fmt.Sprintf("[loading self data - bot %d]", bot.id))
}

func PlayTurn(pl unsafe.Pointer, turn int, allies map[int]Bot, enemies []Bot, warningCount int) (map[int]Action, int) {
	err := ResetGame(pl, turn)
	if err != nil {
		fmt.Printf("error when reseting game: %v\n", err)
		warningCount = WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
	}
	actions := map[int]Action{}
	if !(warningCount > WARNING_TOLERANCE) {
		for _, bot := range allies {
			err = LoadGameBot(pl, bot)
			if err != nil {
				fmt.Printf("error when loading game bot %v: %v\n", bot, err)
				warningCount = WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
				break
			}
		}
		if warningCount <= WARNING_TOLERANCE {
			for _, bot := range enemies {
				err = LoadGameBot(pl, bot)
				if err != nil {
					fmt.Printf("error when loading game bot %v: %v\n", bot, err)
					warningCount = WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
					break
				}
			}
		}
	}
	// TODO: randomize order
	for id, bot := range allies {
		actions[id] = Action{
			actionType: SUICIDE,
			x:          -1,
			y:          -1,
		}
		if !(warningCount > WARNING_TOLERANCE) {
			err = LoadSelf(pl, bot)
			if err != nil {
				fmt.Printf("error when loading self (bot %v) %v\n", bot, err)
				warningCount = WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
				continue
			}
			action, err := GetActionWithTimeout(pl, bot)

			fmt.Printf("bot %d (%v) act (%d,%d,%d), %v\n", id, bot, action.actionType, action.x, action.y, err)
			if errors.Is(err, INVALID_MOVE_ERROR) {
				action.actionType = GUARD
				action.x = -1
				action.y = -1
			} else if errors.Is(err, TIMEOUT_ERROR) {
				warningCount++
				action.actionType = GUARD
				action.x = -1
				action.y = -1
			} else if err != nil {
				fmt.Printf("disqualifed because of %v\n", err)
				warningCount = WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
			}
			if warningCount > WARNING_TOLERANCE {
				continue
			}
			actions[id] = action
		}
	}
	return actions, warningCount
}
