package interfaces

type StartResponse struct {
	Started bool `json:"started"`
}

type StopResponse struct {
	MatchId string `json:"matchId"`
}
