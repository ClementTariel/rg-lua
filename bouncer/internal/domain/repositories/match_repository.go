package repositories

import (
	"github.com/ClementTariel/rg-lua/bouncer/internal/domain/entities"
	"github.com/google/uuid"
)

type MatchRepository interface {
	GetById(id uuid.UUID) (entities.Match, error)
	// TODO: WIP
}
