package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ClementTariel/rg-lua/bouncer/internal/domain/external"
)

const (
	MATCHMAKER_HOST = "matchmaker"
	MATCHMAKER_PORT = 4444
)

type MatchmakerMS struct {
}

func NewMatchmakerMS() external.MatchmakerMS {
	return MatchmakerMS{}
}

func (MatchmakerMS) AddMatchToQueue(blueName string, redName string) (bool, error) {
	postBody, _ := json.Marshal(external.MatchmakerAddMatchToQueueRequest{
		BlueName: blueName,
		RedName:  redName,
	})
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/request-match", MATCHMAKER_HOST, MATCHMAKER_PORT), "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	added := (resp.StatusCode != http.StatusNoContent)
	return added, err
}
