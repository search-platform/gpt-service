package gptservice

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	api "github.com/search-platform/gpt-service/api/gpt"
	gptrepository "github.com/search-platform/gpt-service/internal/gpt-service/gptRepository"
	"github.com/search-platform/gpt-service/internal/pkg/errdetails"

	"github.com/uptrace/bun"
)

var _ api.GptServiceServer = (*Service)(nil)

type Service struct {
	api.UnimplementedGptServiceServer

	cfg ServiceConfig

	db      *bun.DB
	gptRepo *gptrepository.GptRepo
}

func NewService(cfg *ServiceConfig, db *bun.DB) (*Service, error) {
	return &Service{
		cfg:     *cfg,
		db:      db,
		gptRepo: gptrepository.NewGptRepo(cfg.ApiKey),
	}, nil
}

func (s *Service) FindBankInformation(ctx context.Context, req *api.FindBankInformationRequest) (*api.BankInfo, error) {
	errg := errdetails.NewBadRequestBuilder()
	if req.Country == "" {
		errg.Required("Country", "country is required")
	}
	if req.Name == "" {
		errg.Required("Name", "name is required")
	}
	if errg.NotEmpty() {
		return nil, errg.AsError()
	}

	bank, err := s.gptRepo.GetBankInfo(ctx, req.Name, req.Country)
	if err != nil {
		return nil, err
	}

	log.Info().Msg(fmt.Sprintf("got bank info: %v", bank))

	return bank.ToAPI(), nil
}
