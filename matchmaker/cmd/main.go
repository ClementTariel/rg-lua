package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/matchmaker/internal/application/controllers"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/application/services"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/infrastructure/db"
)

type StartRequest struct {
	MatchId uuid.UUID `json:"matchId"`
	BlueId  uuid.UUID `json:"blueId"`
	RedId   uuid.UUID `json:"redId"`
}

const (
	PORT = 4444
)

// TODO WIP
func main() {
	var err error

	// TODO: load user and password from conf or env
	user := "matchmaker_user"
	password := "matchmaker_temporary_password"
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user,
		password,
		"rglua_db",
		5432,
		"rglua")
	postgresDb, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("matchmaker could not connect to DB: %v\n", err)
		return
	}

	botRepo := db.NewBotRepository(postgresDb)

	e := echo.New()

	matchmakerService := services.NewMatchmakerService(botRepo)
	controllers.NewMatchmakerController(e, matchmakerService)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
