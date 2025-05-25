package controllers

import (
	"fmt"
	"net/http"

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

	e.POST("/request-match", controller.AddMatchToQueue)

	return controller
}

func (mc *BouncerController) AddMatchToQueue(c echo.Context) error {
	var addPendingMatchRequest interfaces.AddPendingMatchRequest
	err := c.Bind(&addPendingMatchRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	added, err := mc.bouncerService.AddMatchToQueue(addPendingMatchRequest.BlueName, addPendingMatchRequest.RedName)
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
