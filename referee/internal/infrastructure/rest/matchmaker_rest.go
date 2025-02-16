package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/ClementTariel/rg-lua/referee/internal/domain/external"
	"github.com/ClementTariel/rg-lua/rgcore"
)

const MATCHMAKER_PORT = 4444

type MatchmakerMS struct {
}

func NewMatchmakerMS() external.MatchmakerMS {
	return MatchmakerMS{}
}

func (MatchmakerMS) SaveMatch(matchId uuid.UUID, match []map[int]rgcore.BotState) error {
	postBody, _ := json.Marshal(external.MatchmakerSaveMatchRequest{
		MatchId: matchId,
		Match:   match,
	})
	_, err := http.Post(fmt.Sprintf("http://localhost:%d/save-match", MATCHMAKER_PORT), "application/json", bytes.NewBuffer(postBody))
	return err
}

func (MatchmakerMS) CancelMatch(matchId uuid.UUID, error error) error {
	postBody, _ := json.Marshal(external.MatchmakerCancelMatchRequest{
		MatchId: matchId,
		Error:   error,
	})
	_, err := http.Post(fmt.Sprintf("http://localhost:%d/cancel-match", MATCHMAKER_PORT), "application/json", bytes.NewBuffer(postBody))
	return err
}
