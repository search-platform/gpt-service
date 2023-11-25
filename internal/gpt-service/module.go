package gptservice

import (
	api "github.com/search-platform/gpt-service/api/gpt"
	"github.com/search-platform/gpt-service/internal/pkg/db"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"gptservice",
	fx.Provide(
		NewPublicController,
		fx.Annotate(NewService, fx.As(new(api.GptServiceServer))),
	),
	db.ConnectionModule,
)
