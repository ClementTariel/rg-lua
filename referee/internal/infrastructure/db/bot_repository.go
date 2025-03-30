package db

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/ClementTariel/rg-lua/referee/internal/domain/entities"
	"github.com/ClementTariel/rg-lua/referee/internal/domain/repositories"
)

type BotRepository struct {
	db *sql.DB
}

func NewBotRepository(db *sql.DB) repositories.BotRepository {
	return &BotRepository{
		db: db,
	}
}

func (br *BotRepository) GetById(id uuid.UUID) (entities.Bot, error) {
	stmt, err := br.db.Prepare("SELECT name, script FROM bots WHERE id= $1")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return entities.Bot{}, err
	}
	var bot entities.Bot
	err = stmt.QueryRow(id).Scan(&bot.Name, &bot.Script)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return entities.Bot{}, err
	}
	return bot, nil

}
