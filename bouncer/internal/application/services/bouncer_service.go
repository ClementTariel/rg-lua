package services

import (
	"github.com/ClementTariel/rg-lua/bouncer/internal/domain/entities"
	"github.com/ClementTariel/rg-lua/bouncer/internal/domain/external"
	"github.com/ClementTariel/rg-lua/bouncer/internal/domain/repositories"
	"github.com/ClementTariel/rg-lua/bouncer/internal/infrastructure/rest"
	"github.com/google/uuid"
)

type BouncerService struct {
	botRepo      repositories.BotRepository
	matchRepo    repositories.MatchRepository
	matchmakerMS external.MatchmakerMS
}

func NewBouncerService(botRepo repositories.BotRepository, matchRepo repositories.MatchRepository) BouncerService {
	return BouncerService{
		botRepo:      botRepo,
		matchRepo:    matchRepo,
		matchmakerMS: rest.NewMatchmakerMS(),
	}
}

func (s *BouncerService) AddMatchToQueue(blueName string, redName string) (bool, error) {
	return s.matchmakerMS.AddMatchToQueue(blueName, redName)
}

func (s *BouncerService) GetMatch(matchId uuid.UUID) (entities.Match, error) {
	return s.matchRepo.GetById(matchId)
}

func (s *BouncerService) GetSummaries(start int, size int) ([]entities.MatchSummary, error) {
	return s.matchRepo.GetSummaries(start, size)
}
