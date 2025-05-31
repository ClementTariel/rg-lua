package controllers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/bouncer/internal/application/interfaces"
	"github.com/ClementTariel/rg-lua/bouncer/internal/application/services"
)

type BouncerController struct {
	bouncerService services.BouncerService
}

func NewBouncerController(e *echo.Echo, bouncerService services.BouncerService) *BouncerController {
	controller := &BouncerController{
		bouncerService: bouncerService,
	}

	e.GET("/match/:id", controller.GetMatch)
	e.POST("/request-match", controller.AddMatchToQueue)

	return controller
}

func (bc *BouncerController) GetMatch(c echo.Context) error {
	matchId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	match, err := bc.bouncerService.GetMatch(matchId)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusInternalServerError, "Internal Error")
	}
	stopResponse := &interfaces.GetMatchResponse{
		Match: match,
	}
	return c.JSON(http.StatusOK, stopResponse)
}

func (bc *BouncerController) AddMatchToQueue(c echo.Context) error {
	var addPendingMatchRequest interfaces.AddPendingMatchRequest
	err := c.Bind(&addPendingMatchRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	added, err := bc.bouncerService.AddMatchToQueue(addPendingMatchRequest.BlueName, addPendingMatchRequest.RedName)
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
