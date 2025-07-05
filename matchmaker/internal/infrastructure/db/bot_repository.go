package db

import (
	"database/sql"

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
	id := uuid.Nil
	stmt, err := br.db.Prepare("SELECT id FROM bots WHERE name= $1")
	if err != nil {
		return id, err
	}
	err = stmt.QueryRow(name).Scan(&id)
	return id, err
}

func (br *BotRepository) GetUserNameFromBotId(id uuid.UUID) (string, error) {
	var name string
	stmt, err := br.db.Prepare("SELECT name FROM users AS u JOIN (SELECT bots.userId FROM bots WHERE id=$1) AS b ON u.id=b.userId")
	if err != nil {
		return name, err
	}
	err = stmt.QueryRow(id).Scan(&name)
	return name, err
}
