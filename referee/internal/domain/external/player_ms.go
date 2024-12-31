package external

import "github.com/ClementTariel/rg-lua/rgcore"

type PlayerKillResponse struct {
	Killed bool `json:"killed"`
}

type PlayerInitRequest struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}

type PlayerInitResponse struct {
	WarningCount int `json:"warningCount"`
}

type PlayerPlayTurnRequest struct {
	Turn         int          `json:"turn"`
	Allies       []rgcore.Bot `json:"allies"`
	Enemies      []rgcore.Bot `json:"enemies"`
	WarningCount int          `json:"warningCount"`
}

type PlayerPlayTurnResponse struct {
	Actions      map[int]rgcore.Action `json:"actions"`
	WarningCount int                   `json:"warningCount"`
}

type PlayerMS interface {
	Kill(isBlue bool) (bool, error)
	Init(isBlue bool, name string, script string) (int, error)
	PlayTurn(isBlue bool, turn int, allies []rgcore.Bot, enemies []rgcore.Bot, warningCount int) (map[int]rgcore.Action, int, error)
}
