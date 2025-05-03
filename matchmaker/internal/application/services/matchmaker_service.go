package services

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/entities"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/external"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/repositories"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/infrastructure/rest"
	"github.com/ClementTariel/rg-lua/rgcore"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
)

const MATCH_TIMEOUT = 2 *
	// Convert from nanoseconds to milliseconds
	1000000 *
	// Init time
	((rgcore.BOT_INIT_TIME_BUDGET *
		// number of bots per wave
		2 * rgcore.SPAWN_COUNT *
		// number of waves
		rgcore.MAX_TURN / rgcore.SPAWN_DELAY) +

		// Action time
		(rgcore.BOT_ACTION_TIME_BUDGET *
			// duration of a wave
			rgcore.SPAWN_DELAY *
			// sum of the max number of bots per wave
			((rgcore.MAX_TURN / rgcore.SPAWN_DELAY) *
				((rgcore.MAX_TURN / rgcore.SPAWN_DELAY) + 1) / 2) *
			2 * rgcore.SPAWN_COUNT))

type MatchmakerService struct {
	botRepo       repositories.BotRepository
	matchRepo     repositories.MatchRepository
	refereeMS     external.RefereeMS
	matchQueue    entities.MatchQueue
	isRunning     bool
	currentMatch  entities.PendingMatch
	debounceTimer *time.Timer
	mu            sync.Mutex
}

func NewMatchmakerService(botRepo repositories.BotRepository, matchRepo repositories.MatchRepository) MatchmakerService {
	return MatchmakerService{
		botRepo:    botRepo,
		matchRepo:  matchRepo,
		refereeMS:  rest.NewRefereeMS(),
		matchQueue: entities.NewMatchQueue(),
		isRunning:  false,
	}
}

func (s *MatchmakerService) printGrid(currentGameState map[int]rgcore.BotState) {
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
	blueCount := 0
	redCount := 0
	for _, botState := range currentGameState {
		tile := "O"
		blueCount += 1
		if botState.Bot.PlayerId == rgcore.RED_ID {
			tile = "X"
			blueCount -= 1
			redCount += 1
		}
		tileIndex := ((2*rgcore.GRID_SIZE+1)*botState.Bot.Y + (2 * botState.Bot.X))
		gameStateAsStr = gameStateAsStr[:tileIndex] + tile + gameStateAsStr[tileIndex+1:]
	}
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "# ", "\033[40m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "X ", "\033[41m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "O ", "\033[46m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "\n", "\033[0m\n\033[47m")
	fmt.Printf("%d - %d\n", blueCount, redCount)
	fmt.Printf("\033[47m%s\033[0m\n", gameStateAsStr)
}

func (s *MatchmakerService) SaveMatch(matchId uuid.UUID, game []map[int]rgcore.BotState) error {
	rgdebug.VPrintf("save %v\n", matchId)
	for i, state := range game {
		fmt.Printf("turn %d\n", i+1)
		s.printGrid(state)
	}
	score1 := 0
	score2 := 0
	for _, botState := range game[len(game)-1] {
		if botState.Bot.PlayerId == rgcore.BLUE_ID {
			score1 += 1
		} else {
			score2 += 1
		}
	}
	fmt.Printf("%v - %v\n", score1, score2)
	if matchId != s.currentMatch.Id {
		err := errors.New("Corrupted Match ID")
		fmt.Printf("Error: %v\n", err)
		return err
	}
	err := s.matchRepo.Save(entities.Match{
		Id:       matchId,
		BotId1:   s.currentMatch.BotId1,
		BotId2:   s.currentMatch.BotId2,
		BotName1: s.currentMatch.BotName1,
		BotName2: s.currentMatch.BotName2,
		Date:     time.Now(),
		Game:     game,
		Score1:   score1,
		Score2:   score2,
	})
	s.isRunning = false
	if err != nil {
		rgdebug.VPrintf("Error: %v\n", err)
	}
	err = s.StartDebouncedMatch()
	if err != nil {
		rgdebug.VPrintf("Error: %v\n", err)
	}
	return err
}

func (s *MatchmakerService) CancelMatch(matchId uuid.UUID, err error) bool {
	rgdebug.VPrintf("cancel %v because of %v\n", matchId, err)
	// TODO: cancel match
	return false
}

func (s *MatchmakerService) KillMatch() error {
	defer func(s *MatchmakerService) { s.isRunning = false }(s)
	return s.refereeMS.KillMatch()
}

func (s *MatchmakerService) AddMatchToQueue(blueName string, redName string) (bool, error) {
	rgdebug.VPrintf("Add match to queue: %s - %s\n", blueName, redName)
	// TODO: better system to handle queue size and check on elements added
	if s.matchQueue.IsFull() {
		return false, nil
	}
	blueId, err := s.botRepo.GetIdFromName(blueName)
	if err != nil {
		return false, err
	}
	redId, err := s.botRepo.GetIdFromName(redName)
	if err != nil {
		return false, err
	}
	pendingMatch := entities.PendingMatch{
		Id:       uuid.New(),
		BotId1:   blueId,
		BotId2:   redId,
		BotName1: blueName,
		BotName2: redName,
	}
	added := s.matchQueue.Push(pendingMatch)
	if !s.isRunning {
		err = s.StartDebouncedMatch()
	}
	return added, err
}

func (s *MatchmakerService) StartMatch(pendingMatch entities.PendingMatch) error {
	s.currentMatch = pendingMatch
	s.isRunning = true
	rgdebug.VPrintln("Match started, waiting for the result...")
	return s.refereeMS.StartMatch(pendingMatch.Id, pendingMatch.BotId1, pendingMatch.BotId2)
}

func (s *MatchmakerService) StartDebouncedMatch() error {
	rgdebug.VPrintln("Debounced start...")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isRunning {
		rgdebug.VPrintln("Abort, already start")
		return nil
	}
	if s.debounceTimer != nil {
		s.debounceTimer.Stop()
	}
	pendingMatch, err := s.matchQueue.Pop()
	if err != nil {
		rgdebug.VPrintf("Abort: %v\n", err)
		return err
	}
	s.debounceTimer = time.AfterFunc(MATCH_TIMEOUT, func() {
		s.KillMatch()
		fmt.Println("Kill match because it took too long")
		s.StartDebouncedMatch()
	})
	return s.StartMatch(pendingMatch)
}
