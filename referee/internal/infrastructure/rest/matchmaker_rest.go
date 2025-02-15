package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ClementTariel/rg-lua/referee/internal/domain/external"
	"github.com/ClementTariel/rg-lua/rgcore"
)

type MatchmakerMS struct {
}

func NewMatchmakerMS() external.MatchmakerMS {
	return MatchmakerMS{}
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
