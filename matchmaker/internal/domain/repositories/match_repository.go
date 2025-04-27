package repositories

import (
	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/entities"
	"github.com/google/uuid"
)

type MatchRepository interface {
	GetById(id uuid.UUID) (entities.Match, error)
	Save(match entities.Match) error
}
