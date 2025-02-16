package interfaces

import "github.com/google/uuid"

type StartResponse struct {
	Started bool `json:"started"`
}

type StopResponse struct {
	MatchId uuid.UUID `json:"matchId"`
}
