package main

import (
	"database/sql"
	"flag"
	"fmt"
	"strconv"

	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/matchmaker/internal/application/controllers"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/application/services"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/infrastructure/db"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
)

const (
	PORT = 4444
)

func SetFlags() {
	flag.BoolVar(&rgdebug.VERBOSE, "v", false, "")
	flag.BoolVar(&rgdebug.VERBOSE, "verbose", false, "Show more logs")
	flag.Parse()
}

func main() {
	SetFlags()

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
