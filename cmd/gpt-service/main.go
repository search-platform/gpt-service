package main

import (
	"context"

	"github.com/rs/zerolog/log"
	api "github.com/search-platform/gpt-service/api/gpt"
	gptservice "github.com/search-platform/gpt-service/internal/gpt-service"
	"github.com/search-platform/gpt-service/internal/pkg/grpcserver"
	"github.com/search-platform/gpt-service/internal/pkg/httpserver"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		httpserver.Module,
		gptservice.Module,
		grpcserver.Module,
		fx.Provide(NewAppConfig),
		fx.Invoke(registerHTTP, registerGRPC),
	)
	app.Run()
}

func registerHTTP(
	lfc fx.Lifecycle,
	srv *httpserver.Server,
	lpc *gptservice.PublicController,
) {
	// srv.RegisterMiddleware()
	lpc.RegisterController(srv.App())

	lfc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.Run(); err != nil {
					log.Error().Err(err).Msg("HTTP server error")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown()
		},
	})
}

func registerGRPC(lfc fx.Lifecycle, srv *grpcserver.Server, licensesService api.GptServiceServer) {
	api.RegisterGptServiceServer(srv.Server(), licensesService)

	lfc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msg("Starting GRPC Server")
			go func() {
				if err := srv.Start(); err != nil {
					log.Error().Err(err).Msg("GRPC server error")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := srv.GracefulStop(); err != nil {
				log.Error().Err(err).Msg("GRPC server graceful stop error")
				return err
			}
			return nil
		},
	})
}
