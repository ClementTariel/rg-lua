package external

import "github.com/ClementTariel/rg-lua/rgcore"

type MatchmakerSaveMatchRequest struct {
	MatchId string                    `json:"matchId"`
	Match   []map[int]rgcore.BotState `json:"match"`
}

type MatchmakerCancelMatchRequest struct {
	MatchId string `json:"matchId"`
	Error   error  `json:"error"`
}

type MatchmakerMS interface {
	SaveMatch(matchId string, match []map[int]rgcore.BotState) error
	CancelMatch(matchId string, err error) error
}
