package main

import (
	gptservice "github.com/search-platform/gpt-service/internal/gpt-service"
	"github.com/search-platform/gpt-service/internal/pkg/config"
	"github.com/search-platform/gpt-service/internal/pkg/db"
	"github.com/search-platform/gpt-service/internal/pkg/grpcserver"
	"github.com/search-platform/gpt-service/internal/pkg/httpserver"
	"go.uber.org/fx"
)

type AppConfig struct {
	fx.Out

	*db.DBConfig
	GptServiceConfig *gptservice.ServiceConfig
	*httpserver.Config
	GRPC *grpcserver.Config
}

func NewAppConfig() (cfg AppConfig, err error) {
	cfg.DBConfig = &db.DBConfig{}
	cfg.GptServiceConfig = &gptservice.ServiceConfig{}
	cfg.GRPC = &grpcserver.Config{}
	ep := config.RealEnvParser{}
	if err := config.Parse(ep, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
