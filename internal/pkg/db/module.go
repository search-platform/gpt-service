package db

import (
	"context"

	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

var ConnectionModule = fx.Provide(func(lc fx.Lifecycle, cfg *DBConfig) *bun.DB {
	conn := NewConnection(cfg)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return conn.PingContext(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return conn
})
