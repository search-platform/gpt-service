package gptservice

import (
	"context"

	api "github.com/search-platform/gpt-service/api/gpt"
	gptrepository "github.com/search-platform/gpt-service/internal/gpt-service/gptRepository"
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

func (s *Service) FindBankInformation(ctx context.Context, req *api.FindBankInformationRequest) (*api.FindBankInformationResponse, error) {
	return nil, nil
}
