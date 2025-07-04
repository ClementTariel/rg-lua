package services

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/ClementTariel/rg-lua/rgcore"
	"github.com/ClementTariel/rg-lua/rgcore/lua"
	"github.com/ClementTariel/rg-lua/rgcore/rgconst"
	"github.com/ClementTariel/rg-lua/rgcore/rgdebug"
	"github.com/ClementTariel/rg-lua/rgcore/rgentities"
	"github.com/ClementTariel/rg-lua/rgcore/rgerrors"
)

type PlayerService struct {
	L       unsafe.Pointer
	Running bool
}

func NewPlayerService() PlayerService {
	return PlayerService{
		L:       nil,
		Running: false,
	}
}

func (s *PlayerService) CreateState() error {
	var err error
	if !s.Running {
		s.L, err = rgcore.NewRGState()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		s.Running = true
	}
	return nil
}

func (s *PlayerService) KillCurrentMatch() bool {
	rgdebug.VPrintln("kill")
	killed := false
	if s.Running {
		killed = true
		lua.CloseState(s.L)
		s.Running = false
	}
	rgdebug.VPrintf("kill status: %v\n", killed)
	return killed
}

func (s *PlayerService) InitNewMatch(name string, script string) (int, error) {
	s.KillCurrentMatch()
	err := s.CreateState()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return rgconst.WARNING_TOLERANCE + 1, err
	}
	lua.PushFunction(s.L, rgdebug.GetPrintInLuaFunctionPointer(), "print")
	warningCount, err := rgcore.InitRG(s.L, script, name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return warningCount, err
	}
	rgdebug.VPrintln("[Successfully initialized]")
	return warningCount, nil
}

func (s *PlayerService) ResetGame(turn int) error {
	return lua.RunScript(s.L, rgcore.GetResetScript(turn), "[reset game data]")
}

func (s *PlayerService) LoadGameBot(bot rgentities.Bot) error {
	botId := "nil"
	if (bot.Id) > 0 {
		botId = fmt.Sprintf("%d", bot.Id)
	}
	botDescription := fmt.Sprintf("bot %s", botId)
	if botId == "nil" {
		botDescription = "enemy bot"
	}
	return lua.RunScript(s.L,
		rgcore.GetLoadBotScript(bot.X, bot.Y, bot.Hp, bot.PlayerId, botId),
		fmt.Sprintf("[loading game data - %s]", botDescription))
}

func (s *PlayerService) LoadSelf(bot rgentities.Bot) error {
	return lua.RunScript(s.L,
		rgcore.GetLoadSelfScript(bot.X, bot.Y, bot.Hp, bot.PlayerId, bot.Id),
		fmt.Sprintf("[loading self data - bot %d]", bot.Id))
}

func (s *PlayerService) PlayTurn(turn int, allies []rgentities.Bot, enemies []rgentities.Bot, warningCount int) (map[int]rgentities.Action, int) {
	err := s.ResetGame(turn)
	if err != nil {
		fmt.Printf("error when reseting game: %v\n", err)
		warningCount = rgconst.WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
	}
	actions := map[int]rgentities.Action{}
	if !(warningCount > rgconst.WARNING_TOLERANCE) {
		for _, bot := range append(allies, enemies...) {
			err = s.LoadGameBot(bot)
			if err != nil {
				fmt.Printf("error when loading game bot %v: %v\n", bot, err)
				warningCount = rgconst.WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
				break
			}
		}
	}
	for _, bot := range allies {
		actions[bot.Id] = rgentities.Action{
			ActionType: rgconst.SUICIDE,
			X:          -1,
			Y:          -1,
		}
		if warningCount > rgconst.WARNING_TOLERANCE {
			continue
		}
		err = s.LoadSelf(bot)
		if err != nil {
			fmt.Printf("error when loading self (bot %v) %v\n", bot, err)
			warningCount = rgconst.WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
			continue
		}
		action, err := rgcore.GetActionWithTimeout(s.L, bot)
		switch true {
		case errors.Is(err, rgerrors.TIMEOUT_ERROR):
			warningCount++
			if warningCount > rgconst.WARNING_TOLERANCE {
				break
			}
			fallthrough
		case errors.Is(err, rgerrors.INVALID_MOVE_ERROR):
			action.ActionType = rgconst.GUARD
			action.X = -1
			action.Y = -1
		case errors.Unwrap(err) != nil:
			fmt.Printf("disqualified because of %v\n", err)
			warningCount = rgconst.WARNING_TOLERANCE + 1 // error => instantly triggers all warnings
		}
		if warningCount > rgconst.WARNING_TOLERANCE {
			continue
		}
		actions[bot.Id] = action
	}
	return actions, warningCount
}
