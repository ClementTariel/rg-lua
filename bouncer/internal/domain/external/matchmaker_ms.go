package external

type MatchmakerAddMatchToQueueRequest struct {
	BlueName string `json:"blueName"`
	RedName  string `json:"redName"`
}

type MatchmakerMS interface {
	AddMatchToQueue(string, string) (bool, error)
}
