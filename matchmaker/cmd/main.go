package main

import (
	"database/sql"
	"flag"
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

func SetFlags() {
	flag.Parse()
}

const (
	PORT = 4444
)

// TODO WIP
func main() {
	var err error
	SetFlags()
	tail := flag.Args()
	if len(tail) != 2 {
		fmt.Println("Error: Expected 2 lua files in arguments")
		return
	}
	fileName1 := tail[0]
	fileName2 := tail[1]

	// TODO: WIP connection to db
	// TODO: load user and password from conf or env
	user := "matchmaker_user"
	password := "matchmaker_temporary_password"
	connStr := fmt.Sprintf("user=%s password=%s dbname=rglua sslmode=disable", user, password)
	postgresDb, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	botRepo := db.NewBotRepository(postgresDb)
	blueId, err := botRepo.GetIdFromName(fileName1)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	redId, err := botRepo.GetIdFromName(fileName2)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	matchId := uuid.New()

	e := echo.New()

	matchmakerService := services.NewMatchmakerService(botRepo)
	controllers.NewMatchmakerController(e, matchmakerService)

	matchmakerService.StartMatch(matchId, blueId, redId)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
