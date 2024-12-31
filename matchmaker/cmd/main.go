package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/ClementTariel/rg-lua/matchmaker/internal/application/controllers"
	"github.com/ClementTariel/rg-lua/matchmaker/internal/application/services"
)

type StartRequest struct {
	MatchId string `json:"matchId"`
	Blue    string `json:"blue"`
	Red     string `json:"red"`
}

func SetFlags() {
	flag.Parse()
}

const (
	PORT = 4444
)

// TODO WIP
func main() {
	var err error
	SetFlags()
	tail := flag.Args()
	if len(tail) != 2 {
		fmt.Println("Error: Expected 2 lua files in arguments")
		return
	}

	fileName1 := tail[0]
	fileName2 := tail[1]

	e := echo.New()

	matchmakerService := services.NewMatchmakerService()
	controllers.NewMatchmakerController(e, matchmakerService)

	postBody, _ := json.Marshal(StartRequest{
		Blue:    fileName1,
		Red:     fileName2,
		MatchId: "matchId",
	})
	_, err = http.Post(fmt.Sprintf("http://localhost:%d/start", 3333), "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	// _, err = http.Post(fmt.Sprintf("http://localhost:%d/kill", 3333), "", nil)
	// if err != nil {
	// 	fmt.Printf("%v\n", err)
	// }
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
