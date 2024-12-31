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

func NewPlayerMS() external.PlayerMS {
	return PlayerMS{}
}

func (PlayerMS) Kill(isBlue bool) (bool, error) {
	var port int
	if isBlue {
		port = rgcore.PORT_PLAYER_BLUE
	} else {
		port = rgcore.PORT_PLAYER_RED
	}
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/kill", port), "", nil)
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
	if isBlue {
		port = rgcore.PORT_PLAYER_BLUE
	} else {
		port = rgcore.PORT_PLAYER_RED
	}
	postBody, _ := json.Marshal(external.PlayerInitRequest{
		Name:   name,
		Script: script,
	})
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/init", port), "application/json", bytes.NewBuffer(postBody))
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
	if isBlue {
		port = rgcore.PORT_PLAYER_BLUE
	} else {
		port = rgcore.PORT_PLAYER_RED
	}
	postBody, _ := json.Marshal(external.PlayerPlayTurnRequest{
		Turn:         turn,
		Allies:       allies,
		Enemies:      enemies,
		WarningCount: warningCount,
	})
	var playResponse external.PlayerPlayTurnResponse
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d/play", port), "application/json", bytes.NewBuffer(postBody))
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
