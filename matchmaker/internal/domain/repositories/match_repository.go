package repositories

import (
	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/entities"
)

type MatchRepository interface {
	Save(match entities.Match) error
}
