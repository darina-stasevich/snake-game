package scenes

import (
	"fmt"
	"snake-game/internal/core"
)

type SnakeSegment struct {
	core.Position
}

func NewSnakeSegment(x, y int) *SnakeSegment {
	return &SnakeSegment{core.Position{X: x, Y: y}}
}

type Snake struct {
	Body          []SnakeSegment
	Direction     core.Direction
	NextDirection core.Direction

	IsAlive bool

	minMoveInterval int
	moveInterval    int
	moveTimer       int
}

func NewSnake(x, y, snakeLength, moveInterval, minMoveInterval int) (*Snake, error) {
	if snakeLength < 2 {
		return nil, fmt.Errorf("invalid snake size: expected greater than 1, received %d", snakeLength)
	}

	if moveInterval <= 0 {
		return nil, fmt.Errorf("invalid extendForward interval: expected positive value, received %d", moveInterval)
	}

	body := make([]SnakeSegment, snakeLength)
	for i := 0; i < snakeLength; i++ {
		body[i] = *NewSnakeSegment(x-i, y)
	}

	snake := Snake{
		Body:            body,
		Direction:       core.Right,
		NextDirection:   core.Right,
		IsAlive:         true,
		moveInterval:    moveInterval,
		minMoveInterval: minMoveInterval,
		moveTimer:       0,
	}

	return &snake, nil
}

func (s *Snake) SetNextDirection(direction core.Direction) {
	var isOpposite = false

	switch direction {
	case core.Up:
		if s.Direction == core.Down {
			isOpposite = true
		}
	case core.Down:
		if s.Direction == core.Up {
			isOpposite = true
		}
	case core.Left:
		if s.Direction == core.Right {
			isOpposite = true
		}
	case core.Right:
		if s.Direction == core.Left {
			isOpposite = true
		}
	}
	if !isOpposite {
		s.NextDirection = direction
	}
}

func (s *Snake) Update() bool {
	s.moveTimer++
	if s.moveTimer < s.moveInterval {
		return false
	}
	s.moveTimer = 0
	s.extendForward()
	return true
}

func (s *Snake) extendForward() {
	s.Direction = s.NextDirection
	oldHead := s.Body[0]
	newHeadPos := oldHead.Position
	switch s.Direction {
	case core.Up:
		newHeadPos.Y--
	case core.Left:
		newHeadPos.X--
	case core.Down:
		newHeadPos.Y++
	case core.Right:
		newHeadPos.X++
	}
	newHead := SnakeSegment{newHeadPos}
	s.Body = append([]SnakeSegment{newHead}, s.Body...)
}

func (s *Snake) cutTail() error {
	if len(s.Body) > 1 {
		s.Body = s.Body[0 : len(s.Body)-1]
		return nil
	}
	return fmt.Errorf("invalid command: can't cut tail of snake with size less than 2")
}

func (s *Snake) DecreaseMoveInterval(x int) {
	s.moveInterval = max(s.moveInterval-x, s.minMoveInterval)
}
