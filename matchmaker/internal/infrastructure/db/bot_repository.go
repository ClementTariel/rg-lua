package db

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/repositories"
)

type BotRepository struct {
	db *sql.DB
}

func NewBotRepository(db *sql.DB) repositories.BotRepository {
	return &BotRepository{
		db: db,
	}
}

func (br *BotRepository) GetIdFromName(name string) (uuid.UUID, error) {
	stmt, err := br.db.Prepare("SELECT id FROM bots WHERE name= $1")
	if err != nil {
		fmt.Printf("%v\n", err)
		return uuid.Nil, err
	}
	var id uuid.UUID
	err = stmt.QueryRow(name).Scan(&id)
	if err != nil {
		fmt.Printf("%v\n", err)
		return uuid.Nil, err
	}
	return id, nil

}
