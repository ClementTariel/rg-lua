package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/matchmaker/internal/application/interfaces"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/application/services"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
)

type MatchmakerController struct {
	matchmakerService services.MatchmakerService
}

func NewMatchmakerController(e *echo.Echo, matchmakerService services.MatchmakerService) *MatchmakerController {
	controller := &MatchmakerController{
		matchmakerService: matchmakerService,
	}

	e.POST("/save-match", controller.SaveMatch)
	e.POST("/cancel-match", controller.CancelMatch)

	return controller
}

func (mc *MatchmakerController) SaveMatch(c echo.Context) error {
	var saveMatchRequest interfaces.SaveMatchRequest
	err := c.Bind(&saveMatchRequest)
	if err != nil {
		fmt.Println("Could not bind Save Match request data")
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	saved := mc.matchmakerService.SaveMatch(saveMatchRequest.MatchId, saveMatchRequest.Match)
	rgdebug.VPrintf("saved : %v\n", saved)
	log.Fatal("WIP Force exit")
	return c.JSON(http.StatusOK, nil)
}

func (mc *MatchmakerController) CancelMatch(c echo.Context) error {
	var cancelMatchRequest interfaces.CancelMatchRequest
	err := c.Bind(&cancelMatchRequest)
	if err != nil {
		fmt.Println("Could not bind Cancel Match request data")
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	canceled := mc.matchmakerService.CancelMatch(cancelMatchRequest.MatchId, cancelMatchRequest.Error)
	rgdebug.VPrintf("canceled : %v\n", canceled)
	return c.JSON(http.StatusOK, nil)
}
