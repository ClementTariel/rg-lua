package services

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"

	"github.com/ClementTariel/rg-lua/referee/internal/domain/external"
	"github.com/ClementTariel/rg-lua/referee/internal/domain/repositories"
	"github.com/ClementTariel/rg-lua/referee/internal/infrastructure/rest"
	"github.com/ClementTariel/rg-lua/rgcore/rgconst"
	"github.com/ClementTariel/rg-lua/rgcore/rgentities"
	"github.com/ClementTariel/rg-lua/rgcore/rgutils"
)

type RefereeService struct {
	matchId      uuid.UUID
	botRepo      repositories.BotRepository
	playerMS     external.PlayerMS
	matchmakerMS external.MatchmakerMS
	mu           sync.Mutex
}

func NewRefereeService(botRepo repositories.BotRepository) RefereeService {
	return RefereeService{
		matchId:      uuid.Nil,
		botRepo:      botRepo,
		playerMS:     rest.NewPlayerMS(),
		matchmakerMS: rest.NewMatchmakerMS(),
	}
}

func (s *RefereeService) StopMatch() (uuid.UUID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	matchId := uuid.Nil
	var err error
	if s.matchId != uuid.Nil {
		matchId = s.matchId
		blue := true

		var blueErr error
		var redErr error

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			_, blueErr = s.playerMS.Kill(blue)
			wg.Done()
		}()
		go func() {
			_, redErr = s.playerMS.Kill(!blue)
			wg.Done()
		}()
		wg.Wait()

		if blueErr != nil {
			err = blueErr
			fmt.Printf("An Error Occured : %v\n", err)
		}
		if redErr != nil {
			err = redErr
			fmt.Printf("An Error Occured : %v\n", err)
		}
		s.matchId = uuid.Nil
	}
	return matchId, err
}

func (s *RefereeService) StartMatch(matchId uuid.UUID, blueId uuid.UUID, redId uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.matchId != uuid.Nil {
		err := errors.New("Previous match not killed")
		fmt.Printf("Error: %v\n", err)
		return err
	}
	blueBot, err := s.botRepo.GetById(blueId)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	redBot, err := s.botRepo.GetById(redId)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	blueWarningCount, redWarningCount, err := s.initMatch(blueBot.Name, blueBot.Script, redBot.Name, redBot.Script)
	if err != nil {
		fmt.Printf("An Error Occured : %v\n", err)
		return err
	}
	s.matchId = matchId
	go func(matchId uuid.UUID) {
		match, matchErr := s.playMatch(blueWarningCount, redWarningCount)
		if matchErr != nil {
			fmt.Printf("Error: %v\n", matchErr)
			s.matchmakerMS.CancelMatch(s.matchId, matchErr)
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.matchId != matchId {
			matchErr = errors.New("Match stopped before the end")
			fmt.Printf("Error: %v\n", matchErr)
			return
		}
		// TODO: error handling system to mae sure the match has been saved
		matchErr = s.matchmakerMS.SaveMatch(s.matchId, match)
		fmt.Printf("Error: %v\n", matchErr)
		s.matchId = uuid.Nil
	}(matchId)
	return nil
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
		fmt.Printf("blue error: %v\n", blueErr)
		err = blueErr
	}
	if redErr != nil {
		fmt.Printf("red error: %v\n", redErr)
		err = redErr
	}
	return blueWarningCount, redWarningCount, err

}

func (s *RefereeService) initPlayer(isBlue bool, name string, script string) (int, error) {
	warningCount, err := s.playerMS.Init(isBlue, name, script)
	if err != nil {
		warningCount = rgconst.WARNING_TOLERANCE + 1
	}
	return warningCount, err
}

func (s *RefereeService) playTurn(playerId int, turn int, allBots []rgentities.Bot, previousWarningCount int) (map[int]rgentities.Action, int) {
	actions, warningCount, err := s.playerMS.PlayTurn(playerId == rgconst.BLUE_ID, turn, rgutils.Allies(playerId, allBots), rgutils.Enemies(playerId, allBots), previousWarningCount)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		warningCount = rgconst.WARNING_TOLERANCE + 1
		for _, bot := range rgutils.Allies(playerId, allBots) {
			actions[bot.Id] = rgentities.Action{
				ActionType: rgconst.SUICIDE,
				X:          -1,
				Y:          -1,
			}
		}
	}
	return actions, warningCount
}

func (s *RefereeService) generateSpawnLocations() ([]rgentities.Location, error) {
	var err error
	selectedSpawnLocations := []rgentities.Location{}
	for i := 0; i < rgconst.SPAWN_COUNT; i++ {
		validSpawn := false
		var newSpawn rgentities.Location
		for !validSpawn {
			spawnIndex := rand.Intn(rgconst.SPAWN_LEN)
			newSpawn, err = rgutils.GetSpawnLocation(spawnIndex)
			if err != nil {
				return []rgentities.Location{}, err
			}
			validSpawn = true
			for _, spawn := range selectedSpawnLocations {
				if (spawn == newSpawn) ||
					(rgentities.Location{X: (rgconst.GRID_SIZE - 1 - spawn.X), Y: (rgconst.GRID_SIZE - 1 - spawn.Y)} == newSpawn) {
					validSpawn = false
					break
				}
			}
		}
		selectedSpawnLocations = append(selectedSpawnLocations, newSpawn)
	}
	return selectedSpawnLocations, err
}

func (s *RefereeService) claimLocation(loc rgentities.Location, bot rgentities.Bot, claimedMoves map[rgentities.Location][]rgentities.Bot) {
	botLoc := rgentities.Location{X: bot.X, Y: bot.Y}
	otherBots, ok := claimedMoves[loc]
	if !ok {
		otherBots = []rgentities.Bot{}
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
			potentialSwapBotLoc := rgentities.Location{X: potentialSwapBot.X, Y: potentialSwapBot.Y}
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
	botsInConflict = append([]rgentities.Bot{}, botsInConflict...)
	for _, otherBot := range botsInConflict {
		// Bots involved in the conflict cannot move
		// => they stay in place and claim their current location
		otherBotLoc := rgentities.Location{X: otherBot.X, Y: otherBot.Y}
		s.claimLocation(otherBotLoc, otherBot, claimedMoves)
	}
}

func (s *RefereeService) playMatch(blueWarningCount int, redWarningCount int) ([]map[int]rgentities.BotState, error) {
	game := []map[int]rgentities.BotState{}
	allBots := []rgentities.Bot{}
	lastId := 0
	for turn := 0; turn < rgconst.MAX_TURN; turn++ {
		// Spawn
		if turn%rgconst.SPAWN_DELAY == 0 {
			// Kill bots on spawn tiles
			allBots = rgutils.FilterOutBotsOnSpawn(allBots)
			// Generate random spawns
			newSpawnLocations, err := s.generateSpawnLocations()
			if err != nil {
				return game, err
			}
			for _, loc := range newSpawnLocations {
				lastId += 1
				allBots = append(allBots, rgentities.Bot{
					X:        loc.X,
					Y:        loc.Y,
					Hp:       rgconst.MAX_HP,
					Id:       lastId,
					PlayerId: rgconst.BLUE_ID,
				})
				lastId += 1
				allBots = append(allBots, rgentities.Bot{
					X:        rgconst.GRID_SIZE - 1 - loc.X,
					Y:        rgconst.GRID_SIZE - 1 - loc.Y,
					Hp:       rgconst.MAX_HP,
					Id:       lastId,
					PlayerId: rgconst.RED_ID,
				})
			}
		}

		// Get actions
		var blueActions map[int]rgentities.Action
		var redActions map[int]rgentities.Action
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			blueActions, blueWarningCount = s.playTurn(rgconst.BLUE_ID, turn, allBots, blueWarningCount)
			wg.Done()
		}()
		go func() {
			redActions, redWarningCount = s.playTurn(rgconst.RED_ID, turn, allBots, redWarningCount)
			wg.Done()
		}()
		wg.Wait()

		// Add actions to game state
		currentGameState := map[int]rgentities.BotState{}
		for _, bot := range allBots {
			botState := rgentities.BotState{Bot: bot}
			if bot.PlayerId == rgconst.BLUE_ID {
				botState.Action = blueActions[bot.Id]
			} else {
				botState.Action = redActions[bot.Id]
			}
			currentGameState[bot.Id] = botState
		}
		game = append(game, currentGameState)

		// Apply actions
		for id, botState := range game[len(game)-1] {
			nextLoc := rgentities.Location{X: botState.Action.X, Y: botState.Action.Y}
			if (botState.Action.ActionType == rgconst.MOVE || botState.Action.ActionType == rgconst.ATTACK) &&
				rgutils.GetLocationType(nextLoc.X, nextLoc.Y) == rgconst.OBSTACLE {
				guardAction := rgentities.Action{ActionType: rgconst.GUARD, X: botState.Bot.X, Y: botState.Bot.Y}
				botState.Action = guardAction
				game[len(game)-1][id] = botState
			}
		}
		claimedMoves := map[rgentities.Location][]rgentities.Bot{}
		for _, botState := range game[len(game)-1] {
			var loc rgentities.Location
			if botState.Action.ActionType == rgconst.MOVE {
				loc = rgentities.Location{X: botState.Action.X, Y: botState.Action.Y}
			} else {
				loc = rgentities.Location{X: botState.Bot.X, Y: botState.Bot.Y}
			}
			s.claimLocation(loc, botState.Bot, claimedMoves)
		}
		updatedBots := map[int]rgentities.BotState{}
		// Move bots
		for loc, bots := range claimedMoves {
			if len(bots) == 1 {
				updatedBot := game[len(game)-1][bots[0].Id]
				updatedBot.Bot.X = loc.X
				updatedBot.Bot.Y = loc.Y
				updatedBots[updatedBot.Bot.Id] = updatedBot
				continue
			}
			for _, bot := range bots {
				updatedBots[bot.Id] = game[len(game)-1][bot.Id]
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
			if game[len(game)-1][updatedBot.Bot.Id].Action.ActionType == rgconst.GUARD {
				continue
			}
			updatedBot.Bot.Hp -= rgconst.COLLISION_DAMAGE * len(collisions)
			updatedBots[updatedBot.Bot.Id] = updatedBot
		}
		// Attack & Suicide damage
		damage := map[int]map[rgentities.Location]int{}
		damage[rgconst.BLUE_ID] = map[rgentities.Location]int{}
		damage[rgconst.RED_ID] = map[rgentities.Location]int{}
		for _, botState := range game[len(game)-1] {
			if botState.Action.ActionType == rgconst.ATTACK {
				loc := rgentities.Location{X: botState.Action.X, Y: botState.Action.Y}
				currentDamage, ok := damage[botState.Bot.PlayerId][loc]
				if !ok {
					currentDamage = 0
				}
				damage[botState.Bot.PlayerId][loc] = currentDamage +
					rgconst.ATTACK_DAMAGE_MIN +
					rand.Intn(rgconst.ATTACK_DAMAGE_MAX+1-rgconst.ATTACK_DAMAGE_MIN)
			} else if botState.Action.ActionType == rgconst.SUICIDE {
				updatedBot := updatedBots[botState.Bot.Id]
				updatedBot.Bot.Hp = 0
				updatedBots[updatedBot.Bot.Id] = updatedBot
				for i := 0; i < 4; i++ {
					x := botState.Bot.X + (i%2)*(i-2)
					y := botState.Bot.Y + ((i+1)%2)*(i-1)
					loc := rgentities.Location{X: x, Y: y}
					currentDamage, ok := damage[botState.Bot.PlayerId][loc]
					if !ok {
						currentDamage = 0
					}
					damage[botState.Bot.PlayerId][loc] = currentDamage + rgconst.SUICIDE_DAMAGE
				}

			}
		}
		for _, botState := range updatedBots {
			loc := rgentities.Location{X: botState.Bot.X, Y: botState.Bot.Y}
			totalDamage := 0
			ok := false
			if botState.Bot.PlayerId == rgconst.BLUE_ID {
				totalDamage, ok = damage[rgconst.RED_ID][loc]
			} else {
				totalDamage, ok = damage[rgconst.BLUE_ID][loc]
			}
			if !ok {
				continue
			}
			updatedBot := updatedBots[botState.Bot.Id]
			if botState.Action.ActionType == rgconst.GUARD {
				updatedBot.Bot.Hp -= totalDamage / 2
			} else {
				updatedBot.Bot.Hp -= totalDamage
			}
			updatedBots[updatedBot.Bot.Id] = updatedBot
		}
		// Remove dead bots
		allBots = rgutils.FilterOutDeadBots(updatedBots)
	}
	// Add final state with all robot guarding
	currentGameState := map[int]rgentities.BotState{}
	for _, bot := range allBots {
		botState := rgentities.BotState{Bot: bot}
		botState.Action = rgentities.Action{ActionType: rgconst.GUARD, X: -1, Y: -1}
		currentGameState[bot.Id] = botState
	}
	// TODO: kill players states (call "/kill")
	game = append(game, currentGameState)
	return game, nil
}
