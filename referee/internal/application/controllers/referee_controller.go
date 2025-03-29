package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/referee/internal/application/interfaces"
	"github.com/ClementTariel/rg-lua/referee/internal/application/services"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
)

type RefereeController struct {
	refereeService services.RefereeService
}

func NewRefereeController(e *echo.Echo, refereeService services.RefereeService) *RefereeController {
	controller := &RefereeController{
		refereeService: refereeService,
	}

	e.POST("/stop", controller.StopMatch)
	e.POST("/start", controller.StartMatch)

	return controller
}

func (rc *RefereeController) StopMatch(c echo.Context) error {
	fmt.Println("Forced to stop current match")
	matchId, err := rc.refereeService.StopMatch()
	if err != nil {
		// TODO: handle that case properly
		fmt.Printf("%v\n", err)
		fmt.Println("WIP Throw an Internal Error")
		return c.String(http.StatusInternalServerError, "Internal Error")
	}
	rgdebug.VPrintf("match stopped : %v\n", matchId)
	stopResponse := &interfaces.StopResponse{
		MatchId: matchId,
	}
	return c.JSON(http.StatusOK, stopResponse)
}

func (rc *RefereeController) StartMatch(c echo.Context) error {
	fmt.Println("Start match")
	var startRequest interfaces.StartRequest
	err := c.Bind(&startRequest)
	if err != nil {
		fmt.Printf("%v\n", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	startStatus := rc.refereeService.StartMatch(startRequest.MatchId, startRequest.BlueId, startRequest.RedId)
	startResponse := &interfaces.StartResponse{
		Started: startStatus,
	}
	return c.JSON(http.StatusOK, startResponse)
}
