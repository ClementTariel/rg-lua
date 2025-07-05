package entities

import (
	"time"

	"github.com/google/uuid"
)

type MatchSummary struct {
	Id        uuid.UUID
	BotId1    uuid.UUID
	BotId2    uuid.UUID
	BotName1  string
	BotName2  string
	UserName1 string
	UserName2 string
	Date      time.Time
	Score1    int
	Score2    int
}

type Match struct {
	MatchSummary
	CompressedGame []byte
}
