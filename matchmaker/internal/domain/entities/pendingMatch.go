package entities

import (
	"github.com/google/uuid"
)

type PendingMatch struct {
	Id        uuid.UUID
	BotId1    uuid.UUID
	BotId2    uuid.UUID
	BotName1  string
	BotName2  string
	UserName1 string
	UserName2 string
}
