package controllers

import (
	"fmt"
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
	e.POST("/request-match", controller.AddMatchToQueue)

	return controller
}

func (mc *MatchmakerController) SaveMatch(c echo.Context) error {
	var saveMatchRequest interfaces.SaveMatchRequest
	err := c.Bind(&saveMatchRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	saved := mc.matchmakerService.SaveMatch(saveMatchRequest.MatchId, saveMatchRequest.Match)
	rgdebug.VPrintf("saved : %v\n", saved)
	return c.NoContent(http.StatusNoContent)
}

func (mc *MatchmakerController) CancelMatch(c echo.Context) error {
	var cancelMatchRequest interfaces.CancelMatchRequest
	err := c.Bind(&cancelMatchRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	canceled := mc.matchmakerService.CancelMatch(cancelMatchRequest.MatchId, cancelMatchRequest.Error)
	rgdebug.VPrintf("canceled : %v\n", canceled)
	return c.NoContent(http.StatusNoContent)
}

func (mc *MatchmakerController) AddMatchToQueue(c echo.Context) error {
	var addPendingMatchRequest interfaces.AddPendingMatchRequest
	err := c.Bind(&addPendingMatchRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	added, err := mc.matchmakerService.AddMatchToQueue(addPendingMatchRequest.BlueName, addPendingMatchRequest.RedName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad Request: Invalid bot names")
	}
	if added {
		return c.NoContent(http.StatusAccepted)
	} else {
		return c.NoContent(http.StatusServiceUnavailable)
	}
}
