package db

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/entities"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/repositories"
)

type Game entities.Game

// Implement driver.Valuer interface to allow JsonB encoding
func (g Game) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// Implement sql.Scanner interface to allow JsonB decoding
func (g *Game) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Could not convert value to bytes")
	}
	return json.Unmarshal(bytes, &g)
}

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
	err = stmt.QueryRow(id).Scan(&match)
	return match, err
}

func (br *MatchRepository) Save(match entities.Match) error {
	stmt, err := br.db.Prepare("INSERT INTO matchs (id, botId1, botId2, botName1, botName2, date, game, score1, score2) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);")
	if err != nil {
		return err
	}
	jsonBGame, err := Game(match.Game).Value()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(match.Id, match.BotId1, match.BotId2, match.BotName1, match.BotName2, match.Date, jsonBGame, match.Score1, match.Score2)
	return err
}
