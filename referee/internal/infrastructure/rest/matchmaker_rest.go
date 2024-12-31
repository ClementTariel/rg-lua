package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ClementTariel/rg-lua/referee/internal/domain/external"
	"github.com/ClementTariel/rg-lua/rgcore"
)

type MatchmakerMS struct {
}

func NewMatchmakerMS() external.MatchmakerMS {
	return MatchmakerMS{}
}

func (MatchmakerMS) PlayTurn(isBlue bool, turn int, allies []rgcore.Bot, enemies []rgcore.Bot, warningCount int) (map[int]rgcore.Action, int, error) {
	var port int
	if isBlue {
		port = rgcore.PORT_PLAYER_BLUE
	} else {
		port = rgcore.PORT_PLAYER_RED
	}
	postBody, _ := json.Marshal(external.PlayerPlayTurnRequest{
		Turn:         turn,
		Allies:       allies,
		Enemies:      enemies,
		WarningCount: warningCount,
	})
	var playResponse external.PlayerPlayTurnResponse
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/play", port), "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return map[int]rgcore.Action{}, rgcore.WARNING_TOLERANCE + 1, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[int]rgcore.Action{}, rgcore.WARNING_TOLERANCE + 1, err
	}
	err = json.Unmarshal(body, &playResponse)
	if err != nil {
		return map[int]rgcore.Action{}, rgcore.WARNING_TOLERANCE + 1, err
	}
	return playResponse.Actions, playResponse.WarningCount, nil
}

func (MatchmakerMS) SaveMatch(matchId string, match []map[int]rgcore.BotState) error {
	port := 4444
	postBody, _ := json.Marshal(external.MatchmakerSaveMatchRequest{
		MatchId: matchId,
		Match:   match,
	})
	_, err := http.Post(fmt.Sprintf("http://localhost:%d/save-match", port), "application/json", bytes.NewBuffer(postBody))
	return err
}

func (MatchmakerMS) CancelMatch(matchId string, error error) error {
	port := 4444
	postBody, _ := json.Marshal(external.MatchmakerCancelMatchRequest{
		MatchId: matchId,
		Error:   error,
	})
	_, err := http.Post(fmt.Sprintf("http://localhost:%d/cancel-match", port), "application/json", bytes.NewBuffer(postBody))
	return err
}
