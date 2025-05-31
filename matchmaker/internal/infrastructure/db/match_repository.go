package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/entities"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/repositories"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) repositories.MatchRepository {
	return &MatchRepository{
		db: db,
	}
}

func (br *MatchRepository) Save(match entities.Match) error {
	stmt, err := br.db.Prepare("INSERT INTO matchs (id, botId1, botId2, botName1, botName2, date, compressedGame, score1, score2) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(match.Id, match.BotId1, match.BotId2, match.BotName1, match.BotName2, match.Date, match.CompressedGame, match.Score1, match.Score2)
	return err
}
