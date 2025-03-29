package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ClementTariel/rg-lua/referee/internal/domain/external"
	"github.com/ClementTariel/rg-lua/rgcore"
)

type PlayerMS struct {
}

const (
	BLUE_PLAYER_HOST = "blue_player"
	RED_PLAYER_HOST  = "red_player"

	BLUE_PLAYER_PORT = 1111
	RED_PLAYER_PORT  = 2222
)

func NewPlayerMS() external.PlayerMS {
	return PlayerMS{}
}

func (PlayerMS) Kill(isBlue bool) (bool, error) {
	var port int
	var host string
	if isBlue {
		port = BLUE_PLAYER_PORT
		host = BLUE_PLAYER_HOST
	} else {
		port = RED_PLAYER_PORT
		host = RED_PLAYER_HOST
	}
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/kill", host, port), "", nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var playerKillResponse external.PlayerKillResponse
	err = json.Unmarshal(body, &playerKillResponse)
	if err != nil {
		return false, err
	}
	return playerKillResponse.Killed, nil
}

func (PlayerMS) Init(isBlue bool, name string, script string) (int, error) {
	var port int
	var host string
	if isBlue {
		port = BLUE_PLAYER_PORT
		host = BLUE_PLAYER_HOST
	} else {
		port = RED_PLAYER_PORT
		host = RED_PLAYER_HOST
	}
	postBody, _ := json.Marshal(external.PlayerInitRequest{
		Name:   name,
		Script: script,
	})
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/init", host, port), "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return rgcore.WARNING_TOLERANCE + 1, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rgcore.WARNING_TOLERANCE + 1, err
	}
	var initResponse external.PlayerInitResponse
	err = json.Unmarshal(body, &initResponse)
	if err != nil {
		return rgcore.WARNING_TOLERANCE + 1, err
	}
	return initResponse.WarningCount, nil
}

func (PlayerMS) PlayTurn(isBlue bool, turn int, allies []rgcore.Bot, enemies []rgcore.Bot, warningCount int) (map[int]rgcore.Action, int, error) {
	var port int
	var host string
	if isBlue {
		port = BLUE_PLAYER_PORT
		host = BLUE_PLAYER_HOST
	} else {
		port = RED_PLAYER_PORT
		host = RED_PLAYER_HOST
	}
	postBody, _ := json.Marshal(external.PlayerPlayTurnRequest{
		Turn:         turn,
		Allies:       allies,
		Enemies:      enemies,
		WarningCount: warningCount,
	})
	var playResponse external.PlayerPlayTurnResponse
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/play", host, port), "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return map[int]rgcore.Action{}, rgcore.WARNING_TOLERANCE + 1, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[int]rgcore.Action{}, rgcore.WARNING_TOLERANCE + 1, err
	}
	err = json.Unmarshal(body, &playResponse)
	if err != nil {
		return map[int]rgcore.Action{}, rgcore.WARNING_TOLERANCE + 1, err
	}
	return playResponse.Actions, playResponse.WarningCount, nil
}
