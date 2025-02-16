package external

import "github.com/google/uuid"

type RefereeStartMatchRequest struct {
	MatchId uuid.UUID `json:"matchId"`
	BlueId  uuid.UUID `json:"blueId"`
	RedId   uuid.UUID `json:"redId"`
}

type RefereeMS interface {
	StartMatch(matchId uuid.UUID, blueId uuid.UUID, redId uuid.UUID) error
	KillMatch() error
}
