package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/ClementTariel/rg-lua/matchmaker/internal/domain/external"
)

const (
	REFEREE_HOST = "referee"
	REFEREE_PORT = 3333
)

type RefereeMS struct {
}

func NewRefereeMS() external.RefereeMS {
	return RefereeMS{}
}

func (RefereeMS) StartMatch(matchId uuid.UUID, blueId uuid.UUID, redId uuid.UUID) error {
	postBody, _ := json.Marshal(external.RefereeStartMatchRequest{
		BlueId:  blueId,
		RedId:   redId,
		MatchId: uuid.New(),
	})
	_, err := http.Post(fmt.Sprintf("http://%s:%d/start", REFEREE_HOST, REFEREE_PORT), "application/json", bytes.NewBuffer(postBody))
	return err
}

func (RefereeMS) KillMatch() error {
	_, err := http.Post(fmt.Sprintf("http://%s:%d/stop", REFEREE_HOST, REFEREE_PORT), "", nil)
	return err
}
