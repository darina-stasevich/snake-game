package config

type Config struct {
	ScreenWidth     int
	ScreenHeight    int
	TileSize        int
	InitialSnakeLen int
	InitialSpeed    int
}

func LoadConfig() *Config {
	return &Config{
		ScreenWidth:     1080,
		ScreenHeight:    720,
		TileSize:        20,
		InitialSnakeLen: 4,
		InitialSpeed:    30,
	}
}
