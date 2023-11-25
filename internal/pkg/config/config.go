package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type EnvParser interface {
	Parse(cfg interface{}) error
}

type RealEnvParser struct{}

func (RealEnvParser) Parse(cfg interface{}) error {
	return env.Parse(cfg)
}

func Parse(ep EnvParser, cfg interface{}) error {
	return ep.Parse(cfg)
}

func MustParse(ep EnvParser, cfg interface{}) {
	if err := ep.Parse(cfg); err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}
}
