package services

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/external"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/repositories"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/infrastructure/rest"
	"github.com/ClementTariel/rg-lua/rgcore"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
)

type MatchmakerService struct {
	botRepo   repositories.BotRepository
	refereeMS external.RefereeMS
}

func NewMatchmakerService(botRepo repositories.BotRepository) MatchmakerService {
	return MatchmakerService{
		botRepo:   botRepo,
		refereeMS: rest.NewRefereeMS(),
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

func (s *MatchmakerService) SaveMatch(matchId uuid.UUID, match []map[int]rgcore.BotState) bool {
	rgdebug.VPrintf("save %v\n", matchId)
	for i, state := range match {
		fmt.Printf("turn %d\n", i+1)
		s.printGrid(state)
	}
	score1 := 0
	score2 := 0
	for _, botState := range match[len(match)-1] {
		if botState.Bot.PlayerId == rgcore.BLUE_ID {
			score1 += 1
		} else {
			score2 += 1
		}
	}
	fmt.Printf("%v - %v\n", score1, score2)
	// TODO: save in db
	return false
}

func (s *MatchmakerService) CancelMatch(matchId uuid.UUID, err error) bool {
	fmt.Printf("cancel %v because of %v\n", matchId, err)
	// TODO: cancel match
	return false
}

// TODO: WIP TEMPORARY
func (s *MatchmakerService) StartMatch(matchId uuid.UUID, blueId uuid.UUID, redId uuid.UUID) error {
	return s.refereeMS.StartMatch(matchId, blueId, redId)
}

// TODO: WIP TEMPORARY
func (s *MatchmakerService) KillMatch() error {
	return s.refereeMS.KillMatch()
}
