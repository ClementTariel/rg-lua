package services

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/ClementTariel/rg-lua/referee/internal/domain/external"
	"github.com/ClementTariel/rg-lua/referee/internal/domain/repositories"
	"github.com/ClementTariel/rg-lua/referee/internal/infrastructure/rest"
	"github.com/ClementTariel/rg-lua/rgcore"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
)

type RefereeService struct {
	matchId      string
	botRepo      repositories.BotRepository
	playerMS     external.PlayerMS
	matchmakerMS external.MatchmakerMS
}

func NewRefereeService(botRepo repositories.BotRepository) RefereeService {
	return RefereeService{
		matchId:      "",
		botRepo:      botRepo,
		playerMS:     rest.NewPlayerMS(),
		matchmakerMS: rest.NewMatchmakerMS(),
	}
}

func (s *RefereeService) StopMatch() (string, error) {
	rgdebug.VPrintln("Stop match")
	matchId := ""
	var err error
	if s.matchId != "" {
		matchId = s.matchId
		blue := true
		_, err1 := s.playerMS.Kill(blue)
		_, err2 := s.playerMS.Kill(!blue)
		if err1 != nil {
			err = err1
			fmt.Printf("An Error Occured : %v\n", err)
		}
		if err2 != nil {
			err = err2
			fmt.Printf("An Error Occured : %v\n", err)
		}
		s.matchId = ""
	}
	rgdebug.VPrintf("match stopped : %v\n", matchId)
	return matchId, err
}

func (s *RefereeService) StartMatch(matchId string, blueName string, redName string) bool {
	blueBot, err := s.botRepo.GetByName(blueName)
	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}
	redBot, err := s.botRepo.GetByName(redName)
	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}
	blueWarningCount, redWarningCount, err := s.initMatch(blueBot.Name, blueBot.Script, redBot.Name, redBot.Script)
	if err != nil {
		fmt.Printf("An Error Occured : %v\n", err)
		return false
	}
	s.matchId = matchId
	go func() {
		match, err := s.playMatch(blueWarningCount, redWarningCount)
		if err != nil {
			fmt.Printf("%v\n", err)
			s.matchmakerMS.CancelMatch(s.matchId, err)
			return
		}
		s.matchmakerMS.SaveMatch(s.matchId, match)
	}()
	return true
}

func (s *RefereeService) initMatch(blueName string, blueScript string, redName string, redScript string) (int, int, error) {
	var blueErr error
	blueWarningCount := 0
	var redErr error
	redWarningCount := 0
	blue := true
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		blueWarningCount, blueErr = s.initPlayer(blue, blueName, blueScript)
		wg.Done()
	}()
	go func() {
		redWarningCount, redErr = s.initPlayer(!blue, redName, redScript)
		wg.Done()
	}()
	wg.Wait()
	var err error
	if blueErr != nil {
		err = blueErr
	}
	if redErr != nil {
		err = redErr
	}
	return blueWarningCount, redWarningCount, err

}

func (s *RefereeService) initPlayer(isBlue bool, name string, script string) (int, error) {
	warningCount, err := s.playerMS.Init(isBlue, name, script)
	if err != nil {
		warningCount = rgcore.WARNING_TOLERANCE + 1
	}
	return warningCount, err
}

func (s *RefereeService) playTurn(playerId int, turn int, allBots []rgcore.Bot, previousWarningCount int) (map[int]rgcore.Action, int) {
	actions, warningCount, err := s.playerMS.PlayTurn(playerId == 1, turn, rgcore.Allies(playerId, allBots), rgcore.Enemies(playerId, allBots), previousWarningCount)
	if err != nil {
		warningCount = rgcore.WARNING_TOLERANCE + 1
		for _, bot := range rgcore.Allies(playerId, allBots) {
			actions[bot.Id] = rgcore.Action{
				ActionType: rgcore.SUICIDE,
				X:          -1,
				Y:          -1,
			}
		}
	}
	return actions, warningCount
}

func (s *RefereeService) generateSpawnLocations() ([]rgcore.Location, error) {
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

func (s *RefereeService) claimLocation(loc rgcore.Location, bot rgcore.Bot, claimedMoves map[rgcore.Location][]rgcore.Bot) {
	botLoc := rgcore.Location{X: bot.X, Y: bot.Y}
	otherBots, ok := claimedMoves[loc]
	if !ok {
		otherBots = []rgcore.Bot{}
	}
	// Check if the bot has already claimed this location (base case)
	for _, otherBot := range otherBots {
		if otherBot.Id == bot.Id {
			return
		}
	}
	conflict := len(otherBots) >= 1
	claimedMoves[loc] = append(otherBots, bot)
	if !conflict {
		// Make sure that there aren't 2 bots swapping their places
		potentialSwapBots, botLocIsClaimed := claimedMoves[botLoc]
		if !botLocIsClaimed {
			return
		}
		for _, potentialSwapBot := range potentialSwapBots {
			potentialSwapBotLoc := rgcore.Location{X: potentialSwapBot.X, Y: potentialSwapBot.Y}
			if bot.Id != potentialSwapBot.Id && potentialSwapBotLoc == loc {
				// - bot wants to move from botLoc to potentialSwapBotLoc
				// - potentialSwapBot != bot
				// - potentialSwapBot wants to move from potentialSwapBotLoc to botLoc
				// => potentialSwapBot and bot are trying to swap places
				conflict = true
				break
			}
		}
		if !conflict {
			return
		}
	}
	botsInConflict := claimedMoves[loc]
	// Deep copy to avoid unexpected behaviours
	botsInConflict = append([]rgcore.Bot{}, botsInConflict...)
	for _, otherBot := range botsInConflict {
		// Bots involved in the conflict cannot move
		// => they stay in place and claim their current location
		otherBotLoc := rgcore.Location{X: otherBot.X, Y: otherBot.Y}
		s.claimLocation(otherBotLoc, otherBot, claimedMoves)
	}
}

func (s *RefereeService) playMatch(blueWarningCount int, redWarningCount int) ([]map[int]rgcore.BotState, error) {
	game := []map[int]rgcore.BotState{}
	allBots := []rgcore.Bot{}
	lastId := 0
	for turn := 0; turn < rgcore.MAX_TURN; turn++ {
		// Spawn
		if turn%rgcore.SPAWN_DELAY == 0 {
			// Kill bots on spawn tiles
			allBots = rgcore.FilterOutBotsOnSpawn(allBots)
			// Generate random spawns
			newSpawnLocations, err := s.generateSpawnLocations()
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
		var blueActions map[int]rgcore.Action
		var redActions map[int]rgcore.Action
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			blueActions, blueWarningCount = s.playTurn(1, turn, allBots, blueWarningCount)
			wg.Done()
		}()
		go func() {
			redActions, redWarningCount = s.playTurn(2, turn, allBots, redWarningCount)
			wg.Done()
		}()
		wg.Wait()

		// Add actions to game state
		currentGameState := map[int]rgcore.BotState{}
		for _, bot := range allBots {
			botState := rgcore.BotState{Bot: bot}
			if bot.PlayerId == 1 {
				botState.Action = blueActions[bot.Id]
			} else {
				botState.Action = redActions[bot.Id]
			}
			currentGameState[bot.Id] = botState
		}
		game = append(game, currentGameState)

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
			s.claimLocation(loc, botState.Bot, claimedMoves)
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
		// Collision damage
		collisions := map[int][]int{}
		// Register all collisions
		for _, bots := range claimedMoves {
			if len(bots) == 1 {
				continue
			}
			for _, bot := range bots {
				botCollisions, ok := collisions[bot.Id]
				if !ok {
					botCollisions = []int{}
				}
				for _, otherBot := range bots {
					if otherBot.PlayerId != bot.PlayerId {
						// otherBot deals damage to bot
						collisionAlreadyRegistered := false
						for _, colliderId := range botCollisions {
							if colliderId == otherBot.Id {
								collisionAlreadyRegistered = true
								break
							}
						}
						if !collisionAlreadyRegistered {
							botCollisions = append(botCollisions, otherBot.Id)
						}
					}
				}
				collisions[bot.Id] = botCollisions
			}
		}
		// Deal collision damage
		for botId, collisions := range collisions {
			updatedBot := updatedBots[botId]
			if game[len(game)-1][updatedBot.Id].Action.ActionType == rgcore.GUARD {
				continue
			}
			updatedBot.Hp -= rgcore.COLLISION_DAMAGE * len(collisions)
			updatedBots[updatedBot.Id] = updatedBot
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
		allBots = rgcore.FilterOutDeadBots(updatedBots)
	}
	// Add final state with all robot guarding
	currentGameState := map[int]rgcore.BotState{}
	for _, bot := range allBots {
		botState := rgcore.BotState{Bot: bot}
		botState.Action = rgcore.Action{ActionType: rgcore.GUARD, X: -1, Y: -1}
		currentGameState[bot.Id] = botState
	}
	// TODO: kill players states (call "/kill")
	game = append(game, currentGameState)
	return game, nil
}
