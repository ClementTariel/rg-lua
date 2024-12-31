package db

import (
	"fmt"
	"os"

	"github.com/ClementTariel/rg-lua/referee/internal/domain/entities"
	"github.com/ClementTariel/rg-lua/referee/internal/domain/repositories"
)

type BotRepository struct {
	// TODO: WIP add info necessary to use db
}

func NewBotRepository() repositories.BotRepository {
	// TODO: WIP add info necessary to use db
	return &BotRepository{}
}

func (br *BotRepository) GetByName(name string) (entities.Bot, error) {
	// TODO: WIP retreive from db
	b, err := os.ReadFile(name)
	if err != nil {
		fmt.Printf("%v\n", err)
		return entities.Bot{}, err
	}
	script := string(b)
	return entities.Bot{
		Name:     name,
		Script:   script,
		UserName: "examples",
	}, nil
}
