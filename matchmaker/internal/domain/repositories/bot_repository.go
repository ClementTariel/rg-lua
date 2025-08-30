package repositories

import "github.com/google/uuid"

type BotRepository interface {
	GetIdFromName(name string) (uuid.UUID, error)
	GetUserNameFromBotId(id uuid.UUID) (string, error)
}
