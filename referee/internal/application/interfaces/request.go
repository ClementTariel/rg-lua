package interfaces

type StartRequest struct {
	MatchId string `json:"matchId"`
	Blue    string `json:"blue"`
	Red     string `json:"red"`
}
