package interfaces

import (
	"github.com/ClementTariel/rg-lua/rgcore"
)

type PlayRequest struct {
	Turn         int          `json:"turn"`
	Allies       []rgcore.Bot `json:"allies"`
	Enemies      []rgcore.Bot `json:"enemies"`
	WarningCount int          `json:"warningCount"`
}

type InitRequest struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}
