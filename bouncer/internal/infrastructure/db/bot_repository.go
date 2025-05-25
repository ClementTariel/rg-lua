package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/ClementTariel/rg-lua/bouncer/internal/domain/repositories"
)

type BotRepository struct {
	db *sql.DB
}

func NewBotRepository(db *sql.DB) repositories.BotRepository {
	return &BotRepository{
		db: db,
	}
}
