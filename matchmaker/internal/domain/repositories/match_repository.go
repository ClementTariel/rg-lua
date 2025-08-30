package repositories

import (
	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/entities"
	"github.com/google/uuid"
)

type MatchRepository interface {
	Save(match entities.Match) error
	// TODO: WIP
	GetById(id uuid.UUID) (entities.Match, error)
}
