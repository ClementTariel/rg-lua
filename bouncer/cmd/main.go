package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/bouncer/internal/application/controllers"
	"github.com/ClementTariel/rg-lua/bouncer/internal/application/services"
	"github.com/ClementTariel/rg-lua/bouncer/internal/infrastructure/db"
)

const (
	PORT = 5555
)

func main() {
	e := echo.New()

	// TODO: load user and password from conf or env
	user := "rglua_user"
	password := "rglua_temporary_password"
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user,
		password,
		"rglua_db",
		5432,
		"rglua")
	postgresDb, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("bouncer could not connect to DB: %v\n", err)
		return
	}

	botRepo := db.NewBotRepository(postgresDb)
	matchRepo := db.NewMatchRepository(postgresDb)

	bouncerService := services.NewBouncerService(botRepo, matchRepo)
	controllers.NewBouncerController(e, bouncerService)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
