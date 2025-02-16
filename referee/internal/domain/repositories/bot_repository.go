package repositories

import (
	"github.com/google/uuid"

	"github.com/ClementTariel/rg-lua/referee/internal/domain/entities"
)

type BotRepository interface {
	GetById(id uuid.UUID) (entities.Bot, error)
}
