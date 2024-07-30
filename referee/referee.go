package main

import (
	"fmt"
	"math/rand"
	"strings"
	"unsafe"

	"github.com/ClementTariel/rg-lua/player"
	"github.com/ClementTariel/rg-lua/rgcore"
)

func Allies(playerId int, bots []rgcore.Bot) []rgcore.Bot {
	allies := []rgcore.Bot{}
	for _, bot := range bots {
		if bot.PlayerId == playerId {
			ally := bot
			allies = append(allies, ally)
		}
	}
	return allies
}
func Enemies(playerId int, bots []rgcore.Bot) []rgcore.Bot {
	enemies := []rgcore.Bot{}
	for _, bot := range bots {
		if bot.PlayerId != playerId {
			enemy := bot
			enemy.Id = 0
			enemies = append(enemies, enemy)
		}
	}
	return enemies
}

func RemoveBotsOnSpawn(bots []rgcore.Bot) []rgcore.Bot {
	filteredBots := []rgcore.Bot{}
	for _, bot := range bots {
		if rgcore.GetLocationType(bot.X, bot.Y) != rgcore.SPAWN {
			filteredBots = append(filteredBots, bot)
		}
	}
	return filteredBots
}

func RemoveDeadBots(bots map[int]rgcore.Bot) []rgcore.Bot {
	filteredBots := []rgcore.Bot{}
	for _, bot := range bots {
		if bot.Hp > 0 {
			filteredBots = append(filteredBots, bot)
		}
	}
	return filteredBots
}

func GenerateSpawnLocations() ([]rgcore.Location, error) {
	var err error
	selectedSpawnLocations := []rgcore.Location{}
	for i := 0; i < rgcore.SPAWN_COUNT; i++ {
		validSpawn := false
		var newSpawn rgcore.Location
		for !validSpawn {
			spawnIndex := rand.Intn(rgcore.SPAWN_LEN)
			newSpawn, err = rgcore.GetSpawnLocation(spawnIndex)
			if err != nil {
				return []rgcore.Location{}, err
			}
			validSpawn = true
			for _, spawn := range selectedSpawnLocations {
				if (spawn == newSpawn) ||
					(rgcore.Location{X: (rgcore.GRID_SIZE - 1 - spawn.X), Y: (rgcore.GRID_SIZE - 1 - spawn.Y)} == newSpawn) {
					validSpawn = false
					break
				}
			}
		}
		selectedSpawnLocations = append(selectedSpawnLocations, newSpawn)
	}
	return selectedSpawnLocations, err
}

func ClaimLocation(loc rgcore.Location, bot rgcore.Bot, claimedMoves map[rgcore.Location][]rgcore.Bot) {
	botLoc := rgcore.Location{X: bot.X, Y: bot.Y}
	otherBots, ok := claimedMoves[loc]
	if !ok {
		otherBots = []rgcore.Bot{}
	}
	// Check if the bot has already claimed this location
	for _, otherBot := range otherBots {
		if otherBot.Id == bot.Id {
			return
		}
	}
	claimedMoves[loc] = append(otherBots, bot)
	conflict := len(otherBots) >= 1
	if !conflict {
		// Make sure that there aren't 2 bots swapping their places
		potentialSwapBots, botLocIsClaimed := claimedMoves[botLoc]
		if !botLocIsClaimed {
			return
		}
		for _, potentialSwapBot := range potentialSwapBots {
			potentialSwapBotLoc := rgcore.Location{X: potentialSwapBot.X, Y: potentialSwapBot.Y}
			if bot.Id != potentialSwapBot.Id && potentialSwapBotLoc == loc {
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
	if len(claimedMoves[loc]) == 2 {
		// New conflict with another bot
		// (otherwise the conflict has already been propagated)
		otherBot := claimedMoves[loc][0]
		otherBotLoc := rgcore.Location{X: otherBot.X, Y: otherBot.Y}
		ClaimLocation(otherBotLoc, otherBot, claimedMoves)
	}
	otherBots, ok = claimedMoves[botLoc]
	if !ok || len(otherBots) == 0 {
		// No conflict (yet) for the claim of botLoc
		claimedMoves[botLoc] = []rgcore.Bot{bot}
		return
	}
	if len(otherBots) > 1 {
		claimedMoves[botLoc] = append(otherBots, bot)
		// Conflicts already detected for the claim of botLoc
		return
	}
	// New conflict with another bot
	otherBot := otherBots[0]
	otherBotLoc := rgcore.Location{X: otherBot.X, Y: otherBot.Y}
	claimedMoves[botLoc] = append(otherBots, bot)
	ClaimLocation(otherBotLoc, otherBot, claimedMoves)
}

func printGrid(currentGameState map[int]rgcore.BotState) {
	gameStateAsStr := ""
	for i := 0; i < rgcore.GRID_SIZE; i++ {
		for j := 0; j < rgcore.GRID_SIZE; j++ {
			tile := " "
			if rgcore.GetLocationType(j, i) == rgcore.OBSTACLE {
				tile = "#"
			}
			gameStateAsStr += tile + " "
		}
		gameStateAsStr += "\n"
	}
	for _, botState := range currentGameState {
		tile := "O"
		if botState.Bot.PlayerId == 1 {
			tile = "X"
		}
		tileIndex := ((2*rgcore.GRID_SIZE+1)*botState.Bot.Y + (2 * botState.Bot.X))
		gameStateAsStr = gameStateAsStr[:tileIndex] + tile + gameStateAsStr[tileIndex+1:]
	}
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "# ", "\033[40m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "X ", "\033[41m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "O ", "\033[46m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "\n", "\033[0m\n\033[47m")
	fmt.Printf("\033[47m%s\033[0m\n", gameStateAsStr)
}

func PlayGame(pl1 unsafe.Pointer, pl2 unsafe.Pointer) ([]map[int]rgcore.BotState, error) {
	game := []map[int]rgcore.BotState{}

	allBots := []rgcore.Bot{}
	warningCount1 := 0
	warningCount2 := 0

	lastId := 0
	for turn := 0; turn < rgcore.MAX_TURN; turn++ {
		// Spawn
		if turn%rgcore.SPAWN_DELAY == 0 {
			// Kill bots on spawn tiles
			allBots = RemoveBotsOnSpawn(allBots)
			// Generate random spawns
			newSpawnLocations, err := GenerateSpawnLocations()
			if err != nil {
				return game, err
			}
			for _, loc := range newSpawnLocations {
				lastId += 1
				allBots = append(allBots, rgcore.Bot{
					X:        loc.X,
					Y:        loc.Y,
					Hp:       rgcore.MAX_HP,
					Id:       lastId,
					PlayerId: 1,
				})
				lastId += 1
				allBots = append(allBots, rgcore.Bot{
					X:        rgcore.GRID_SIZE - 1 - loc.X,
					Y:        rgcore.GRID_SIZE - 1 - loc.Y,
					Hp:       rgcore.MAX_HP,
					Id:       lastId,
					PlayerId: 2,
				})
			}
		}

		// Get actions
		var actions1 map[int]rgcore.Action
		if warningCount1 > rgcore.WARNING_TOLERANCE {
			actions1 = map[int]rgcore.Action{}
			for _, bot := range Allies(1, allBots) {
				actions1[bot.Id] = rgcore.Action{
					ActionType: rgcore.SUICIDE,
					X:          -1,
					Y:          -1,
				}
			}
		} else {
			actions1, warningCount1 = player.PlayTurn(
				pl1,
				turn,
				Allies(1, allBots),
				Enemies(1, allBots),
				warningCount1)
		}

		var actions2 map[int]rgcore.Action
		if warningCount2 > rgcore.WARNING_TOLERANCE {
			actions2 = map[int]rgcore.Action{}
			for _, bot := range Allies(2, allBots) {
				actions2[bot.Id] = rgcore.Action{
					ActionType: rgcore.SUICIDE,
					X:          -1,
					Y:          -1,
				}
			}
		} else {
			actions2, warningCount2 = player.PlayTurn(
				pl2,
				turn,
				Allies(2, allBots),
				Enemies(2, allBots),
				warningCount2)
		}

		// Add actions to game state
		currentGameState := map[int]rgcore.BotState{}
		for _, bot := range allBots {
			botState := rgcore.BotState{Bot: bot}
			if bot.PlayerId == 1 {
				botState.Action = actions1[bot.Id]
			} else {
				botState.Action = actions2[bot.Id]
			}
			currentGameState[bot.Id] = botState
		}
		game = append(game, currentGameState)
		fmt.Printf("turn %d\n", len(game))
		printGrid(currentGameState)

		// Apply actions
		for id, botState := range game[len(game)-1] {
			nextLoc := rgcore.Location{X: botState.Action.X, Y: botState.Action.Y}
			if (botState.Action.ActionType == rgcore.MOVE || botState.Action.ActionType == rgcore.ATTACK) &&
				rgcore.GetLocationType(nextLoc.X, nextLoc.Y) == rgcore.OBSTACLE {
				guardAction := rgcore.Action{ActionType: rgcore.GUARD, X: botState.Bot.X, Y: botState.Bot.Y}
				botState.Action = guardAction
				game[len(game)-1][id] = botState
			}
		}
		claimedMoves := map[rgcore.Location][]rgcore.Bot{}
		for _, botState := range game[len(game)-1] {
			var loc rgcore.Location
			if botState.Action.ActionType == rgcore.MOVE {
				loc = rgcore.Location{X: botState.Action.X, Y: botState.Action.Y}
			} else {
				loc = rgcore.Location{X: botState.Bot.X, Y: botState.Bot.Y}
			}
			ClaimLocation(loc, botState.Bot, claimedMoves)
		}
		updatedBots := map[int]rgcore.Bot{}
		// Move bots
		for loc, bots := range claimedMoves {
			if len(bots) == 1 {
				updatedBot := bots[0]
				updatedBot.X = loc.X
				updatedBot.Y = loc.Y
				updatedBots[updatedBot.Id] = updatedBot
				continue
			}
			for _, bot := range bots {
				updatedBot := bot
				updatedBots[updatedBot.Id] = updatedBot
			}
		}
		// Collision damages
		for _, bots := range claimedMoves {
			if len(bots) == 1 {
				continue
			}
			for _, bot := range bots {
				updatedBot := updatedBots[bot.Id]
				if game[len(game)-1][updatedBot.Id].Action.ActionType == rgcore.GUARD {
					continue
				}
				for _, otherBot := range bots {
					if bot.PlayerId != otherBot.PlayerId {
						updatedBot.Hp -= rgcore.COLLISION_DAMAGE
					}
				}
				updatedBots[updatedBot.Id] = updatedBot
			}
		}
		// Attack & Suicide damages
		damages := map[rgcore.Location]int{}
		for _, botState := range game[len(game)-1] {
			if botState.Action.ActionType == rgcore.ATTACK {
				loc := rgcore.Location{X: botState.Action.X, Y: botState.Action.Y}
				currentDamages, ok := damages[loc]
				if !ok {
					currentDamages = 0
				}
				damages[loc] = currentDamages +
					rgcore.ATTACK_DAMAGE_MIN +
					rand.Intn(rgcore.ATTACK_DAMAGE_MAX+1-rgcore.ATTACK_DAMAGE_MIN)
			} else if botState.Action.ActionType == rgcore.SUICIDE {
				updatedBot := updatedBots[botState.Bot.Id]
				updatedBot.Hp = 0
				updatedBots[updatedBot.Id] = updatedBot
				for i := 0; i < 4; i++ {
					x := botState.Bot.X + (i%2)*(i-2)
					y := botState.Bot.Y + ((i+1)%2)*(i-1)
					loc := rgcore.Location{X: x, Y: y}
					currentDamages, ok := damages[loc]
					if !ok {
						currentDamages = 0
					}
					damages[loc] = currentDamages + rgcore.SUICIDE_DAMAGE
				}

			}
		}
		for _, botState := range game[len(game)-1] {
			loc := rgcore.Location{X: botState.Bot.X, Y: botState.Bot.Y}
			totalDamages, ok := damages[loc]
			if !ok {
				continue
			}
			updatedBot := updatedBots[botState.Bot.Id]
			if botState.Action.ActionType == rgcore.GUARD {
				updatedBot.Hp -= totalDamages / 2
			} else {
				updatedBot.Hp -= totalDamages
			}
			updatedBots[updatedBot.Id] = updatedBot
		}
		// Remove dead bots
		allBots = RemoveDeadBots(updatedBots)
	}
	// Add final state with all robot guarding
	currentGameState := map[int]rgcore.BotState{}
	for _, bot := range allBots {
		botState := rgcore.BotState{Bot: bot}
		botState.Action = rgcore.Action{ActionType: rgcore.GUARD, X: -1, Y: -1}
		currentGameState[bot.Id] = botState
	}
	game = append(game, currentGameState)
	fmt.Println("RESULT")
	printGrid(currentGameState)
	return game, nil
}
