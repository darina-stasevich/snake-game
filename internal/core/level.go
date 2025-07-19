package core

type Level struct {
	Name string `json:"name"`

	GridWidth  int `json:"grid_width"`
	GridHeight int `json:"grid_height"`

	Walls []Wall `json:"walls"`
}

func NewLevel(name string, width, height int, walls []Wall) *Level {
	return &Level{
		Name:       name,
		GridWidth:  width,
		GridHeight: height,
		Walls:      walls,
	}
}
