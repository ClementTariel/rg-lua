package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/ClementTariel/rg-lua/referee/internal/domain/external"
	"github.com/ClementTariel/rg-lua/rgcore/rgentities"
)

const (
	MATCHMAKER_HOST = "matchmaker"
	MATCHMAKER_PORT = 4444
)

type MatchmakerMS struct {
}

func NewMatchmakerMS() external.MatchmakerMS {
	return MatchmakerMS{}
}

func (MatchmakerMS) SaveMatch(matchId uuid.UUID, game []map[int]rgentities.BotState) error {
	postBody, _ := json.Marshal(external.MatchmakerSaveMatchRequest{
		MatchId: matchId,
		Game:    game,
	})
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/save-match", MATCHMAKER_HOST, MATCHMAKER_PORT), "application/json", bytes.NewBuffer(postBody))
	resp.Body.Close()
	return err
}

func (MatchmakerMS) CancelMatch(matchId uuid.UUID, error error) error {
	postBody, _ := json.Marshal(external.MatchmakerCancelMatchRequest{
		MatchId: matchId,
		Error:   error,
	})
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/cancel-match", MATCHMAKER_HOST, MATCHMAKER_PORT), "application/json", bytes.NewBuffer(postBody))
	resp.Body.Close()
	return err
}
