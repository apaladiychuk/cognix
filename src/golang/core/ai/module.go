package ai

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
)

var ChunkingModule = fx.Options(
	fx.Provide(func() (*ChunkingConfig, error) {
		cfg := ChunkingConfig{}
		if err := utils.ReadConfig(&cfg); err != nil {
			return nil, err
		}
		if err := cfg.Validate(); err != nil {
			return nil, err
		}
		return &cfg, nil
	},
		newChunking,
	),
)

func newChunking(cfg *ChunkingConfig) Chunking {
	if cfg.Strategy == StrategyLLM {
		return NewLLMChunking()
	}
	return NewStaticChunking(cfg)
}
