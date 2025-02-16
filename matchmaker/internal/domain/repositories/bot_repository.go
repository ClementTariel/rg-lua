package repositories

import "github.com/google/uuid"

type BotRepository interface {
	GetIdFromName(name string) (uuid.UUID, error)
}
