package entities

import "github.com/google/uuid"

type Bot struct {
	Id     uuid.UUID
	Name   string
	Script string
}
