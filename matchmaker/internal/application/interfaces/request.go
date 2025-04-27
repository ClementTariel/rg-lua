package interfaces

import (
	"github.com/google/uuid"

	"github.com/ClementTariel/rg-lua/rgcore"
)

type SaveMatchRequest struct {
	MatchId uuid.UUID                 `json:"matchId"`
	Game    []map[int]rgcore.BotState `json:"game"`
}

type CancelMatchRequest struct {
	MatchId uuid.UUID `json:"matchId"`
	Error   error     `json:"error"`
}

type AddPendingMatchRequest struct {
	BlueName string `json:"blueName"`
	RedName  string `json:"redName"`
}
