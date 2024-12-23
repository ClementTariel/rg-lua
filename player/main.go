package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/rgcore"
)

var PORT int

type PlayRequest struct {
	Turn         int          `json:"turn"`
	Allies       []rgcore.Bot `json:"allies"`
	Enemies      []rgcore.Bot `json:"enemies"`
	WarningCount int          `json:"warningCount"`
}

type PlayResponse struct {
	Actions      map[int]rgcore.Action `json:"actions"`
	WarningCount int                   `json:"warningCount"`
}

type InitRequest struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

type InitResponse struct {
	WarningCount int `json:"warningCount"`
}

func SetFlags() {
	flag.BoolVar(&rgcore.VERBOSE, "v", false, "")
	flag.BoolVar(&rgcore.VERBOSE, "verbose", false, "Show more logs")
	flag.IntVar(&PORT, "p", PORT, "")
	flag.Parse()
}

func main() {
	SetFlags()

	e := echo.New()

	e.GET("/kill", func(c echo.Context) error {
		killed := KillCurrentMatch()
		var respString string
		if killed {
			respString = "1"
		} else {
			respString = "0"
		}
		rgcore.VPrintf("kill response: %s\n", respString)
		return c.String(http.StatusOK, respString)
	})

	e.POST("/init", func(c echo.Context) error {
		var initRequest InitRequest
		err = c.Bind(&initRequest)
		if err != nil {
			fmt.Println("bad request")
			return c.String(http.StatusBadRequest, "bad request")
		}
		rgcore.VPrintf("name : %s\n", initRequest.Name)
		rgcore.VPrintf("script:\n%s\n", initRequest.Script)
		var warningCount int
		warningCount, err = InitNewMatch(initRequest.Name, initRequest.Script)
		if err != nil {
			// TODO: WIP make distinction between error caused by internal bug VS bad lua code
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			initResponse := &InitResponse{
				WarningCount: warningCount,
			}
			return c.JSON(http.StatusOK, initResponse)
		}
	})

	e.POST("/play", func(c echo.Context) error {
		var playRequest PlayRequest
		err := c.Bind(&playRequest)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		actions, warningCount := PlayTurn(playRequest.Turn,
			playRequest.Allies,
			playRequest.Enemies,
			playRequest.WarningCount)

		playResponse := &PlayResponse{
			Actions:      actions,
			WarningCount: warningCount,
		}
		return c.JSON(http.StatusOK, playResponse)

	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
