package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
		MatchId: matchId,
	})
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/start", REFEREE_HOST, REFEREE_PORT), "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			var errMessage string
			err = json.Unmarshal(body, &errMessage)
			if err == nil {
				err = errors.New(errMessage)
			}
		}
	}
	return err
}

func (RefereeMS) KillMatch() error {
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/stop", REFEREE_HOST, REFEREE_PORT), "", nil)
	resp.Body.Close()
	return err
}
