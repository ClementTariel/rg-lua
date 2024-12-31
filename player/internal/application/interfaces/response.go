package interfaces

import (
	"github.com/ClementTariel/rg-lua/rgcore"
)

type PlayResponse struct {
	Actions      map[int]rgcore.Action `json:"actions"`
	WarningCount int                   `json:"warningCount"`
}

type InitResponse struct {
	WarningCount int `json:"warningCount"`
}

type KillResponse struct {
	Killed bool `json:"killed"`
}