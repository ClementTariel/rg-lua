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
	"github.com/ClementTariel/rg-lua/rgcore/rgconst"
	"github.com/ClementTariel/rg-lua/rgcore/rgentities"
	"github.com/ClementTariel/rg-lua/rgcore/rgutils"
)

const MATCH_TIMEOUT = 2 *
	// Convert from nanoseconds to milliseconds
	1000000 *
	// Init time
	((rgconst.BOT_INIT_TIME_BUDGET *
		// number of bots per wave
		2 * rgconst.SPAWN_COUNT *
		// number of waves
		rgconst.MAX_TURN / rgconst.SPAWN_DELAY) +

		// Action time
		(rgconst.BOT_ACTION_TIME_BUDGET *
			// duration of a wave
			rgconst.SPAWN_DELAY *
			// sum of the max number of bots per wave
			((rgconst.MAX_TURN / rgconst.SPAWN_DELAY) *
				((rgconst.MAX_TURN / rgconst.SPAWN_DELAY) + 1) / 2) *
			2 * rgconst.SPAWN_COUNT))

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

func (s *MatchmakerService) printGrid(currentGameState map[int]rgentities.BotState) {
	gameStateAsStr := ""
	for i := 0; i < rgconst.GRID_SIZE; i++ {
		for j := 0; j < rgconst.GRID_SIZE; j++ {
			tile := " "
			if rgutils.GetLocationType(j, i) == rgconst.OBSTACLE {
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
		if botState.Bot.PlayerId == rgconst.RED_ID {
			tile = "X"
			blueCount -= 1
			redCount += 1
		}
		tileIndex := ((2*rgconst.GRID_SIZE+1)*botState.Bot.Y + (2 * botState.Bot.X))
		gameStateAsStr = gameStateAsStr[:tileIndex] + tile + gameStateAsStr[tileIndex+1:]
	}
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "# ", "\033[40m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "X ", "\033[41m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "O ", "\033[46m  \033[47m")
	gameStateAsStr = strings.ReplaceAll(gameStateAsStr, "\n", "\033[0m\n\033[47m")
	fmt.Printf("%d - %d\n", blueCount, redCount)
	fmt.Printf("\033[47m%s\033[0m\n", gameStateAsStr)
}

func (s *MatchmakerService) SaveMatch(matchId uuid.UUID, game []map[int]rgentities.BotState) error {
	for i, state := range game {
		fmt.Printf("turn %d\n", i+1)
		s.printGrid(state)
	}
	score1 := 0
	score2 := 0
	for _, botState := range game[len(game)-1] {
		if botState.Bot.PlayerId == rgconst.BLUE_ID {
			score1 += 1
		} else {
			score2 += 1
		}
	}
	fmt.Printf("%v - %v\n", score1, score2)
	s.mu.Lock()
	defer s.mu.Unlock()
	if matchId != s.currentMatch.Id {
		err := errors.New("Corrupted Match ID")
		fmt.Printf("Error: %v\n", err)
		return err
	}
	s.isRunning = false
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
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	go s.StartDebouncedMatch()
	return err
}

func (s *MatchmakerService) CancelMatch(matchId uuid.UUID, err error) error {
	// TODO: cancel match
	return nil
}

func (s *MatchmakerService) KillMatch() error {
	s.mu.Lock()
	defer func() {
		s.isRunning = false
		s.mu.Unlock()
	}()
	return s.refereeMS.KillMatch()
}

func (s *MatchmakerService) AddMatchToQueue(blueName string, redName string) (bool, error) {
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
	go s.StartDebouncedMatch()
	return added, nil
}

func (s *MatchmakerService) StartMatch(pendingMatch entities.PendingMatch) error {
	s.currentMatch = pendingMatch
	s.isRunning = true
	return s.refereeMS.StartMatch(pendingMatch.Id, pendingMatch.BotId1, pendingMatch.BotId2)
}

func (s *MatchmakerService) StartDebouncedMatch() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isRunning {
		return nil
	}
	if s.debounceTimer != nil {
		s.debounceTimer.Stop()
	}
	pendingMatch, err := s.matchQueue.Pop()
	if err != nil {
		return err
	}
	s.debounceTimer = time.AfterFunc(MATCH_TIMEOUT, func() {
		s.KillMatch()
		fmt.Println("Kill match because it took too long")
		err := s.StartDebouncedMatch()
		for err != nil {
			fmt.Printf("Error: %v", err)
			if s.matchQueue.IsEmpty() {
				fmt.Println("Queue is empty")
				return
			}
			fmt.Println("Try next match in queue")
			err = s.StartDebouncedMatch()
		}
	})
	return s.StartMatch(pendingMatch)
}
