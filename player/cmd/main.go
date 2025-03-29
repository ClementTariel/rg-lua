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
	PORT = 1111
)

func SetFlags() {
	flag.BoolVar(&rgdebug.VERBOSE, "v", false, "")
	flag.BoolVar(&rgdebug.VERBOSE, "verbose", false, "Show more logs")
	flag.IntVar(&PORT, "p", PORT, "Port")
	flag.Parse()
}

func main() {
	SetFlags()

	e := echo.New()

	playerService := services.NewPlayerService()
	controllers.NewPlayerController(e, playerService)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
