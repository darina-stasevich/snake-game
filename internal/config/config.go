package config

import (
	"log/slog"
)

type Config struct {
	ScreenWidth           int
	ScreenHeight          int
	TileSize              int
	InitialSnakeLen       int
	InitialSpeed          int
	SpeedIncreaseInterval int
	SpeedIncreaseAmount   int
	Logger                *slog.Logger
}

func LoadConfig() *Config {
	return &Config{
		ScreenWidth:           480,
		ScreenHeight:          360,
		TileSize:              20,
		InitialSnakeLen:       2,
		InitialSpeed:          45,
		SpeedIncreaseInterval: 5,
		SpeedIncreaseAmount:   5,
	}
}

func (config *Config) SetLogger(logger *slog.Logger) {
	config.Logger = logger
}
