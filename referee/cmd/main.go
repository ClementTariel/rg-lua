package main

import (
	"flag"
	"strconv"

	"github.com/ClementTariel/rg-lua/referee/internal/application/controllers"
	"github.com/ClementTariel/rg-lua/referee/internal/application/services"
	"github.com/ClementTariel/rg-lua/referee/internal/infrastructure/db"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
	"github.com/labstack/echo"
)

const (
	DEFAULT_PRINT_MEMORY_BUDGET = (1 << 15)
	MAX_FILE_SIZE               = 1 << 16
	REFEREE_PORT                = 3333
)

var PRINT_MEMORY_BUDGET = DEFAULT_PRINT_MEMORY_BUDGET

func SetFlags() {
	flag.BoolVar(&rgdebug.VERBOSE, "v", false, "")
	flag.BoolVar(&rgdebug.VERBOSE, "verbose", false, "Show more logs")
	flag.IntVar(&PRINT_MEMORY_BUDGET, "m", DEFAULT_PRINT_MEMORY_BUDGET, "")
	flag.IntVar(&PRINT_MEMORY_BUDGET, "memory", DEFAULT_PRINT_MEMORY_BUDGET, "Memory budget")
	flag.Parse()
	if !rgdebug.VERBOSE {
		PRINT_MEMORY_BUDGET = 0
	}
}

func main() {
	SetFlags()
	rgdebug.SetPrintMemoryBudget(PRINT_MEMORY_BUDGET)

	e := echo.New()

	botRepo := db.NewBotRepository()
	refereeService := services.NewRefereeService(botRepo)
	controllers.NewRefereeController(e, refereeService)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(REFEREE_PORT)))
}
