package gptservice

import (
	api "github.com/search-platform/gpt-service/api/gpt"
	"github.com/uptrace/bun"
)

var _ api.GptServiceServer = (*Service)(nil)

type Service struct {
	api.UnimplementedGptServiceServer

	cfg ServiceConfig

	db *bun.DB
}

func NewService(cfg *ServiceConfig, db *bun.DB) (*Service, error) {
	return &Service{
		cfg: *cfg,
		db:  db,
	}, nil
}
