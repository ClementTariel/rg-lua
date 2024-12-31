package services

import (
	"fmt"
	"strings"

	"github.com/ClementTariel/rg-lua/rgcore"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
)

type MatchmakerService struct {
}

func NewMatchmakerService() MatchmakerService {
	return MatchmakerService{}
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

func (s *MatchmakerService) SaveMatch(matchId string, match []map[int]rgcore.BotState) bool {
	rgdebug.VPrintf("save %v\n", matchId)
	for i, state := range match {
		fmt.Printf("turn %d\n", i+1)
		s.printGrid(state)
	}
	score1 := 0
	score2 := 0
	for _, botState := range match[len(match)-1] {
		if botState.Bot.PlayerId == 1 {
			score1 += 1
		} else {
			score2 += 1
		}
	}
	fmt.Printf("%v - %v\n", score1, score2)
	// TODO: save in db
	return false
}

func (s *MatchmakerService) CancelMatch(matchId string, err error) bool {
	fmt.Printf("cancel %v because of %v\n", matchId, err)
	// TODO: cancel match
	return false
}
