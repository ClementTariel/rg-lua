package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/player/internal/application/interfaces"
	"github.com/ClementTariel/rg-lua/player/internal/application/services"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
)

type PlayerController struct {
	playerService services.PlayerService
}

func NewPlayerController(e *echo.Echo, playerService services.PlayerService) *PlayerController {
	controller := &PlayerController{
		playerService: playerService,
	}

	e.POST("/kill", controller.KillPlayer)
	e.POST("/init", controller.InitPlayer)
	e.POST("/play", controller.Play)

	return controller
}

func (pc *PlayerController) KillPlayer(c echo.Context) error {
	killed := pc.playerService.KillCurrentMatch()
	rgdebug.VPrintf("killed : %v\n", killed)
	killResponse := &interfaces.KillResponse{
		Killed: killed,
	}
	return c.JSON(http.StatusOK, killResponse)
}

func (pc *PlayerController) InitPlayer(c echo.Context) error {
	var initRequest interfaces.InitRequest
	err := c.Bind(&initRequest)
	if err != nil {
		fmt.Println("Could not bind Init request data")
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	rgdebug.VPrintf("name : %s\n", initRequest.Name)
	rgdebug.VPrintf("script :\n%s\n", initRequest.Script)
	warningCount, err := pc.playerService.InitNewMatch(initRequest.Name, initRequest.Script)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	initResponse := &interfaces.InitResponse{
		WarningCount: warningCount,
	}
	return c.JSON(http.StatusOK, initResponse)
}

func (pc *PlayerController) Play(c echo.Context) error {
	var playRequest interfaces.PlayRequest
	err := c.Bind(&playRequest)
	if err != nil {
		fmt.Println("Could not bind Play request data")
		return c.String(http.StatusBadRequest, "bad request")
	}
	actions, warningCount := pc.playerService.PlayTurn(playRequest.Turn,
		playRequest.Allies,
		playRequest.Enemies,
		playRequest.WarningCount)
	playResponse := &interfaces.PlayResponse{
		Actions:      actions,
		WarningCount: warningCount,
	}
	return c.JSON(http.StatusOK, playResponse)
}
