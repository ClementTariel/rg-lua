package entities

import (
	"time"

	"github.com/google/uuid"

	"github.com/ClementTariel/rg-lua/rgcore/rgentities"
)

type Game []map[int]rgentities.BotState

type Match struct {
	Id             uuid.UUID
	BotId1         uuid.UUID
	BotId2         uuid.UUID
	BotName1       string
	BotName2       string
	UserName1      string
	UserName2      string
	Date           time.Time
	CompressedGame []byte
	Score1         int
	Score2         int
}
