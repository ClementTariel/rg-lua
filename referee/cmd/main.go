package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ClementTariel/rg-lua/referee/internal/application/controllers"
	"github.com/ClementTariel/rg-lua/referee/internal/application/services"
	"github.com/ClementTariel/rg-lua/referee/internal/infrastructure/db"
	"github.com/labstack/echo"
)

const (
	MAX_FILE_SIZE = 1 << 16
	REFEREE_PORT  = 3333
)

func main() {
	e := echo.New()

	// TODO: load user and password from conf or env
	user := "referee_user"
	password := "referee_temporary_password"
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user,
		password,
		"rglua_db",
		5432,
		"rglua")
	postgresDb, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Referee could not connect to DB: %v\n", err)
		return
	}

	botRepo := db.NewBotRepository(postgresDb)
	refereeService := services.NewRefereeService(botRepo)
	controllers.NewRefereeController(e, refereeService)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(REFEREE_PORT)))
}
