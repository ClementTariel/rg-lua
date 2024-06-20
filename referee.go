package main

import (
	"math/rand"
	"unsafe"
)

func Allies(playerId int, bots []Bot) []Bot {
	allies := []Bot{}
	for _, bot := range bots {
		if bot.playerId == playerId {
			ally := bot
			allies = append(allies, ally)
		}
	}
	return allies
}
func Enemies(playerId int, bots []Bot) []Bot {
	enemies := []Bot{}
	for _, bot := range bots {
		if bot.playerId != playerId {
			enemy := bot
			enemy.id = 0
			enemies = append(enemies, enemy)
		}
	}
	return enemies
}

func RemoveBotsOnSpawn(bots []Bot, grid map[Location]LocationType) []Bot {
	filteredBots := []Bot{}
	for _, bot := range bots {
		if grid[Location{x: bot.x, y: bot.y}] != SPAWN {
			filteredBots = append(filteredBots, bot)
		}
	}
	return filteredBots
}

func RemoveDeadBots(bots map[int]Bot) []Bot {
	filteredBots := []Bot{}
	for _, bot := range bots {
		if bot.hp > 0 {
			filteredBots = append(filteredBots, bot)
		}
	}
	return filteredBots
}

func GenerateSpawnLocations(spawnLocations []Location) []Location {
	selectedSpawnLocations := []Location{}
	for i := 0; i < SPAWN_COUNT; i++ {
		validSpawn := false
		var newSpawn Location
		for !validSpawn {
			newSpawn = spawnLocations[rand.Intn(len(spawnLocations))]
			validSpawn = true
			for _, spawn := range selectedSpawnLocations {
				if (spawn == newSpawn) ||
					(Location{x: (GRID_SIZE - 1 - spawn.x), y: (GRID_SIZE - 1 - spawn.y)} == newSpawn) {
					validSpawn = false
					break
				}
			}
		}
		selectedSpawnLocations = append(selectedSpawnLocations, newSpawn)
	}
	return selectedSpawnLocations
}

func ClaimLocation(loc Location, bot Bot, claimedMoves map[Location][]Bot, grid map[Location]LocationType) {
	botLoc := Location{x: bot.x, y: bot.y}
	otherBots, ok := claimedMoves[loc]
	if !ok {
		otherBots = []Bot{}
	}
	for _, otherBot := range otherBots {
		if otherBot.id == bot.id {
			return
		}
	}
	claimedMoves[loc] = append(otherBots, bot)
	conflict := len(otherBots) >= 1 || (grid[loc] != NORMAL && grid[loc] != SPAWN)
	if !conflict {
		// Make sure that there aren't 2 bots swapping their places
		// (this conflict cannot be detected in the previous check)
		potentialSwapBots, potentialSwapBotsFound := claimedMoves[botLoc]
		if !potentialSwapBotsFound {
			return
		}
		for _, potentialSwapBot := range potentialSwapBots {
			potentialSwapBotLoc := Location{x: potentialSwapBot.x, y: potentialSwapBot.y}
			if bot.id != potentialSwapBot.id && potentialSwapBotLoc == loc {
				// - bot wants to move to potentialSwapBotLoc
				// - potentialSwapBot != bot
				// - potentialSwapBot wants to move to loc
				// => potentialSwapBot and bot are trying to swap places
				conflict = true
				break
			}
		}
		if !conflict {
			return
		}
	}
	otherBots, ok = claimedMoves[botLoc]
	if !ok {
		otherBots = []Bot{}
	}
	for _, otherBot := range otherBots {
		if otherBot.id == bot.id {
			return
		}
	}
	claimedMoves[botLoc] = append(otherBots, bot)
	for _, otherBot := range otherBots {
		if otherBot.id != bot.id {
			otherBotLoc := Location{x: otherBot.x, y: otherBot.y}
			ClaimLocation(otherBotLoc, otherBot, claimedMoves, grid)
		}
	}

}

func PlayGame(pl1 unsafe.Pointer, pl2 unsafe.Pointer) ([]map[int]BotState, error) {
	game := []map[int]BotState{}

	// TODO: precompute
	spawnLocations := []Location{}
	grid := map[Location]LocationType{}
	radius := (GRID_SIZE - 3) * 0.5
	center := (GRID_SIZE - 1) / 2
	min2 := (radius - 0.5) * (radius - 0.5)
	max2 := (radius + 0.5) * (radius + 0.5)
	for i := 0; i < GRID_SIZE; i++ {
		for j := 0; j < GRID_SIZE; j++ {
			loc := Location{x: i, y: j}
			d2 := (float64(i-center))*(float64(i-center)) +
				(float64(j-center))*(float64(j-center))
			if d2 < min2 {
				grid[loc] = NORMAL
			} else if d2 < max2 {
				spawnLocations = append(spawnLocations, Location{x: i, y: j})
				grid[loc] = SPAWN
			} else {
				grid[loc] = OBSTACLE
			}
		}
	}

	allBots := []Bot{}
	warningCount1 := 0
	warningCount2 := 0

	lastId := 0
	for turn := 0; turn < MAX_TURN; turn++ {
		// Spawn
		if turn%SPAWN_DELAY == 0 {
			// Kill bots on spawn tiles
			allBots = RemoveBotsOnSpawn(allBots, grid)
			// Generate random spawns
			newSpawnLocations := GenerateSpawnLocations(spawnLocations)
			for _, loc := range newSpawnLocations {
				lastId += 1
				allBots = append(allBots, Bot{
					x:        loc.x,
					y:        loc.y,
					hp:       MAX_HP,
					id:       lastId,
					playerId: 1,
				})
				lastId += 1
				allBots = append(allBots, Bot{
					x:        GRID_SIZE - 1 - loc.x,
					y:        GRID_SIZE - 1 - loc.y,
					hp:       MAX_HP,
					id:       lastId,
					playerId: 2,
				})
			}
		}

		// Get actions
		var actions1 map[int]Action
		if warningCount1 > WARNING_TOLERANCE {
			actions1 = map[int]Action{}
			for _, bot := range Allies(1, allBots) {
				actions1[bot.id] = Action{
					actionType: SUICIDE,
					x:          -1,
					y:          -1,
				}
			}
		} else {
			actions1, warningCount1 = PlayTurn(
				pl1,
				turn,
				Allies(1, allBots),
				Enemies(1, allBots),
				warningCount1)
		}

		var actions2 map[int]Action
		if warningCount2 > WARNING_TOLERANCE {
			actions2 = map[int]Action{}
			for _, bot := range Allies(2, allBots) {
				actions2[bot.id] = Action{
					actionType: SUICIDE,
					x:          -1,
					y:          -1,
				}
			}
		} else {
			actions2, warningCount2 = PlayTurn(
				pl2,
				turn,
				Allies(2, allBots),
				Enemies(2, allBots),
				warningCount2)
		}

		// Add actions to game state
		currentGameState := map[int]BotState{}
		for _, bot := range allBots {
			botState := BotState{bot: bot}
			if bot.playerId == 1 {
				botState.action = actions1[bot.id]
			} else {
				botState.action = actions2[bot.id]
			}
			currentGameState[bot.id] = botState
		}
		game = append(game, currentGameState)

		// Apply actions
		claimedMoves := map[Location][]Bot{}
		for _, botState := range game[len(game)-1] {
			var loc Location
			if botState.action.actionType == MOVE {
				loc = Location{x: botState.action.x, y: botState.action.y}
			} else {
				loc = Location{x: botState.bot.x, y: botState.bot.y}
			}
			ClaimLocation(loc, botState.bot, claimedMoves, grid)
		}
		updatedBots := map[int]Bot{}
		// Move bots
		for loc, bots := range claimedMoves {
			if len(bots) == 1 {
				updatedBot := bots[0]
				updatedBot.x = loc.x
				updatedBot.y = loc.y
				updatedBots[updatedBot.id] = updatedBot
				continue
			}
			for _, bot := range bots {
				updatedBot := bot
				updatedBots[updatedBot.id] = updatedBot
			}
		}
		// Collision damages
		for _, bots := range claimedMoves {
			if len(bots) == 1 {
				continue
			}
			for _, bot := range bots {
				updatedBot := updatedBots[bot.id]
				if game[len(game)-1][updatedBot.id].action.actionType == GUARD {
					continue
				}
				for _, otherBot := range bots {
					if bot.playerId != otherBot.playerId {
						updatedBot.hp -= COLLISION_DAMAGE
					}
				}
				updatedBots[updatedBot.id] = updatedBot
			}
		}
		// Attack & Suicide damages
		damages := map[Location]int{}
		for _, botState := range game[len(game)-1] {
			if botState.action.actionType == ATTACK {
				loc := Location{x: botState.action.x, y: botState.action.y}
				currentDamages, ok := damages[loc]
				if !ok {
					currentDamages = 0
				}
				damages[loc] = currentDamages +
					ATTACK_DAMAGE_MIN +
					rand.Intn(ATTACK_DAMAGE_MAX+1-ATTACK_DAMAGE_MIN)
			} else if botState.action.actionType == SUICIDE {
				updatedBot := updatedBots[botState.bot.id]
				updatedBot.hp = 0
				updatedBots[updatedBot.id] = updatedBot
				for i := 0; i < 4; i++ {
					x := botState.bot.x + (i%2)*(i-2)
					y := botState.bot.y + ((i+1)%2)*(i-1)
					loc := Location{x: x, y: y}
					currentDamages, ok := damages[loc]
					if !ok {
						currentDamages = 0
					}
					damages[loc] = currentDamages + SUICIDE_DAMAGE
				}

			}
		}
		for _, botState := range game[len(game)-1] {
			loc := Location{x: botState.bot.x, y: botState.bot.y}
			totalDamages, ok := damages[loc]
			if !ok {
				continue
			}
			updatedBot := updatedBots[botState.bot.id]
			if botState.action.actionType == GUARD {
				updatedBot.hp -= totalDamages / 2
			} else {
				updatedBot.hp -= totalDamages
			}
			updatedBots[updatedBot.id] = updatedBot
		}
		// Remove dead bots
		allBots = RemoveDeadBots(updatedBots)
	}
	// Add final state with all robot guarding
	currentGameState := map[int]BotState{}
	for _, bot := range allBots {
		botState := BotState{bot: bot}
		botState.action = Action{actionType: GUARD, x: -1, y: -1}
		currentGameState[bot.id] = botState
	}
	game = append(game, currentGameState)
	return game, nil
}
