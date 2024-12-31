package main

import (
	"flag"
	"strconv"

	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/player/internal/application/controllers"
	"github.com/ClementTariel/rg-lua/player/internal/application/services"
	"github.com/ClementTariel/rg-lua/rgcore"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
)

var (
	PORT      int
	TRUE_BLUE bool
)

func SetFlags() {
	flag.BoolVar(&rgdebug.VERBOSE, "v", false, "")
	flag.BoolVar(&rgdebug.VERBOSE, "verbose", false, "Show more logs")
	flag.BoolVar(&TRUE_BLUE, "blue", false, "determine if the player is blue or red")
	flag.Parse()
	if TRUE_BLUE {
		PORT = rgcore.PORT_PLAYER_BLUE
	} else {
		PORT = rgcore.PORT_PLAYER_RED
	}
}

func main() {
	SetFlags()

	e := echo.New()

	playerService := services.NewPlayerService()
	controllers.NewPlayerController(e, playerService)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
