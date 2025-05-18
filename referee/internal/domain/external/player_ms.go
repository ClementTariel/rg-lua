package external

import "github.com/ClementTariel/rg-lua/rgcore/rgentities"

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
	Turn         int              `json:"turn"`
	Allies       []rgentities.Bot `json:"allies"`
	Enemies      []rgentities.Bot `json:"enemies"`
	WarningCount int              `json:"warningCount"`
}

type PlayerPlayTurnResponse struct {
	Actions      map[int]rgentities.Action `json:"actions"`
	WarningCount int                       `json:"warningCount"`
}

type PlayerMS interface {
	Kill(isBlue bool) (bool, error)
	Init(isBlue bool, name string, script string) (int, error)
	PlayTurn(isBlue bool, turn int, allies []rgentities.Bot, enemies []rgentities.Bot, warningCount int) (map[int]rgentities.Action, int, error)
}
