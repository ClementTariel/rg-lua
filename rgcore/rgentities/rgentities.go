package rgentities

/*
#include "rgentities.h"
*/
import "C"

type (
	ActionType int
	Action     struct {
		ActionType ActionType
		X          int
		Y          int
	}

	Location struct {
		X int
		Y int
	}
	LocationType int

	Bot struct {
		X        int
		Y        int
		Hp       int
		Id       int
		PlayerId int
	}
	BotState struct {
		Bot    Bot
		Action Action
	}
)
