package external

import (
	"github.com/ClementTariel/rg-lua/rgcore/rgentities"
	"github.com/google/uuid"
)

type MatchmakerSaveMatchRequest struct {
	MatchId uuid.UUID                     `json:"matchId"`
	Game    []map[int]rgentities.BotState `json:"game"`
}

type MatchmakerCancelMatchRequest struct {
	MatchId uuid.UUID `json:"matchId"`
	Error   error     `json:"error"`
}

type MatchmakerMS interface {
	SaveMatch(matchId uuid.UUID, match []map[int]rgentities.BotState) error
	CancelMatch(matchId uuid.UUID, err error) error
}
