package interfaces

import "github.com/google/uuid"

type StopResponse struct {
	MatchId uuid.UUID `json:"matchId"`
}
