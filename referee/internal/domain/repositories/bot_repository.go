package repositories

import (
	"github.com/ClementTariel/rg-lua/referee/internal/domain/entities"
)

type BotRepository interface {
	GetByName(name string) (entities.Bot, error)
}
