package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/referee/internal/application/interfaces"
	"github.com/ClementTariel/rg-lua/referee/internal/application/services"
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
	matchId, err := rc.refereeService.StopMatch()
	if err != nil {
		// TODO: handle that case properly
		fmt.Printf("Error during match stop: %v\n", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	stopResponse := &interfaces.StopResponse{
		MatchId: matchId,
	}
	return c.JSON(http.StatusOK, stopResponse)
}

func (rc *RefereeController) StartMatch(c echo.Context) error {
	var startRequest interfaces.StartRequest
	err := c.Bind(&startRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	err = rc.refereeService.StartMatch(startRequest.MatchId, startRequest.BlueId, startRequest.RedId)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
