package interfaces

type GetMatchRequest struct {
	MatchId string `param:"id"`
}

type AddPendingMatchRequest struct {
	BlueName string `json:"blueName"`
	RedName  string `json:"redName"`
}
