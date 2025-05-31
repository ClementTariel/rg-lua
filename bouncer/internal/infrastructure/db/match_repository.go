package db

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/ClementTariel/rg-lua/bouncer/internal/domain/entities"
	"github.com/ClementTariel/rg-lua/bouncer/internal/domain/repositories"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) repositories.MatchRepository {
	return &MatchRepository{
		db: db,
	}
}

func (br *MatchRepository) GetById(id uuid.UUID) (entities.Match, error) {
	var match entities.Match
	stmt, err := br.db.Prepare("SELECT * FROM matchs WHERE id=$1")
	if err != nil {
		return match, err
	}
	err = stmt.QueryRow(id).Scan(&match.Id, &match.BotId1, &match.BotId2, &match.BotName1, &match.BotName2, &match.Date, &match.CompressedGame, &match.Score1, &match.Score2)
	if err != nil {
		return match, err
	}
	return match, err
}
