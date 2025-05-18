package main

import (
	"flag"
	"strconv"

	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/player/internal/application/controllers"
	"github.com/ClementTariel/rg-lua/player/internal/application/services"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
)

var (
	DEFAULT_PRINT_MEMORY_BUDGET = (1 << 15)
	PORT                        = 1111
)

var PRINT_MEMORY_BUDGET = DEFAULT_PRINT_MEMORY_BUDGET

func SetFlags() {
	flag.BoolVar(&rgdebug.VERBOSE, "v", false, "")
	flag.BoolVar(&rgdebug.VERBOSE, "verbose", false, "Show more logs")
	flag.IntVar(&PRINT_MEMORY_BUDGET, "m", DEFAULT_PRINT_MEMORY_BUDGET, "")
	flag.IntVar(&PRINT_MEMORY_BUDGET, "memory", DEFAULT_PRINT_MEMORY_BUDGET, "Memory budget")
	flag.IntVar(&PORT, "p", PORT, "Port")
	flag.Parse()
	if !rgdebug.VERBOSE {
		PRINT_MEMORY_BUDGET = 0
	}
}

func main() {
	SetFlags()

	e := echo.New()

	playerService := services.NewPlayerService()
	controllers.NewPlayerController(e, playerService)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
