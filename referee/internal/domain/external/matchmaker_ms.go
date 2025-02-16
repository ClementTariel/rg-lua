package external

import (
	"github.com/google/uuid"

	"github.com/ClementTariel/rg-lua/rgcore"
)

type MatchmakerSaveMatchRequest struct {
	MatchId uuid.UUID                 `json:"matchId"`
	Match   []map[int]rgcore.BotState `json:"match"`
}

type MatchmakerCancelMatchRequest struct {
	MatchId uuid.UUID `json:"matchId"`
	Error   error     `json:"error"`
}

type MatchmakerMS interface {
	SaveMatch(matchId uuid.UUID, match []map[int]rgcore.BotState) error
	CancelMatch(matchId uuid.UUID, err error) error
}
