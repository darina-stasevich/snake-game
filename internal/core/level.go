package core

type Level struct {
	Name string `json:"name"`

	GridWidth  int `json:"grid_width"`
	GridHeight int `json:"grid_height"`

	Walls []Wall `json:"walls"`
}

func NewEmptyLevel(width, height int) *Level {
	return &Level{
		Name:       "unknown",
		GridWidth:  width,
		GridHeight: height,
	}
}
