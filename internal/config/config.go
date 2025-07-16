package config

import (
	"log/slog"
)

type Config struct {
	ScreenWidth           int
	ScreenHeight          int
	TopBarHeight          int
	TileSize              int
	InitialSnakeLen       int
	InitialSpeed          int
	SpeedIncreaseInterval int
	SpeedIncreaseAmount   int
	Logger                *slog.Logger
}

func LoadConfig() *Config {
	return &Config{
		ScreenWidth:           2400,
		ScreenHeight:          1200,
		TopBarHeight:          30,
		TileSize:              120,
		InitialSnakeLen:       2,
		InitialSpeed:          30,
		SpeedIncreaseInterval: 5,
		SpeedIncreaseAmount:   5,
	}
}

func (config *Config) SetLogger(logger *slog.Logger) {
	config.Logger = logger
}
