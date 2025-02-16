package interfaces

import "github.com/google/uuid"

type StartRequest struct {
	MatchId uuid.UUID `json:"matchId"`
	BlueId  uuid.UUID `json:"blueId"`
	RedId   uuid.UUID `json:"redId"`
}
