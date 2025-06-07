package interfaces

import (
	"github.com/ClementTariel/rg-lua/bouncer/internal/domain/entities"
)

type GetMatchResponse struct {
	Match entities.Match `json:"match"`
}

type GetSummariesResponse struct {
	Summaries []entities.MatchSummary `json:"summaries"`
}
