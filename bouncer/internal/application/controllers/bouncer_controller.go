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
	e.GET("/highlighted-match", controller.GetHighlightedMatch)
	e.GET("/matchs", controller.GetSummaries)
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
	getMatchResponse := &interfaces.GetMatchResponse{
		Match: match,
	}
	return c.JSON(http.StatusOK, getMatchResponse)
}

func (bc *BouncerController) GetHighlightedMatch(c echo.Context) error {
	match, err := bc.bouncerService.GetHighlightedMatch()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusInternalServerError, "Internal Error")
	}
	if match == nil {
		return c.String(http.StatusNotFound, "Not found")
	}
	getMatchResponse := &interfaces.GetMatchResponse{
		Match: *match,
	}
	return c.JSON(http.StatusOK, getMatchResponse)
}

func (bc *BouncerController) GetSummaries(c echo.Context) error {
	var getSummariesRequest interfaces.GetSummariesRequest
	err := c.Bind(&getSummariesRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	summaries, err := bc.bouncerService.GetSummaries(getSummariesRequest.Start, getSummariesRequest.Size)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c.String(http.StatusInternalServerError, "Internal Error")
	}
	getSummariesResponse := &interfaces.GetSummariesResponse{
		Summaries: summaries,
	}
	return c.JSON(http.StatusOK, getSummariesResponse)
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
