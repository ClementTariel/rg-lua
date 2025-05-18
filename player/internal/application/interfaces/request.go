package interfaces

import (
	"github.com/ClementTariel/rg-lua/rgcore/rgentities"
)

type PlayRequest struct {
	Turn         int              `json:"turn"`
	Allies       []rgentities.Bot `json:"allies"`
	Enemies      []rgentities.Bot `json:"enemies"`
	WarningCount int              `json:"warningCount"`
}

type InitRequest struct {
	Name   string `json:"name"`
	Script string `json:"script"`
}
