package game

import (
	"fmt"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type SnakeSegment struct {
	Position
}

func NewSnakeSegment(x, y int) *SnakeSegment {
	return &SnakeSegment{Position{X: x, Y: y}}
}

type Snake struct {
	Body          []SnakeSegment
	Direction     Direction
	NextDirection Direction

	IsAlive  bool
	AteApple bool

	moveInterval int
	moveTimer    int
}

func NewSnake(x, y, snakeLength, moveInterval int) (*Snake, error) {
	if snakeLength < 2 {
		return nil, fmt.Errorf("invalid snake size: expected greater than 1, received %d", snakeLength)
	}

	if moveInterval <= 0 {
		return nil, fmt.Errorf("invalid move interval: expected positive value, received %d", moveInterval)
	}

	body := make([]SnakeSegment, snakeLength)
	for i := 0; i < snakeLength; i++ {
		body[i] = *NewSnakeSegment(x-i, y)
	}

	snake := Snake{
		Body:          body,
		Direction:     Right,
		NextDirection: Right,
		IsAlive:       true,
		AteApple:      false,
		moveInterval:  moveInterval,
		moveTimer:     0,
	}

	return &snake, nil
}
