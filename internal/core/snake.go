package core

import (
	"fmt"
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
		Direction:       Right,
		NextDirection:   Right,
		IsAlive:         true,
		moveInterval:    moveInterval,
		minMoveInterval: minMoveInterval,
		moveTimer:       0,
	}

	return &snake, nil
}

func (s *Snake) SetNextDirection(direction Direction) {
	var isOpposite = false

	switch direction {
	case Up:
		if s.Direction == Down {
			isOpposite = true
		}
	case Down:
		if s.Direction == Up {
			isOpposite = true
		}
	case Left:
		if s.Direction == Right {
			isOpposite = true
		}
	case Right:
		if s.Direction == Left {
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
	case Up:
		newHeadPos.Y--
	case Left:
		newHeadPos.X--
	case Down:
		newHeadPos.Y++
	case Right:
		newHeadPos.X++
	}
	newHead := SnakeSegment{newHeadPos}
	s.Body = append([]SnakeSegment{newHead}, s.Body...)
}

func (s *Snake) CutTail() error {
	if len(s.Body) > 1 {
		s.Body = s.Body[0 : len(s.Body)-1]
		return nil
	}
	return fmt.Errorf("invalid command: can't cut tail of snake with size less than 2")
}

func (s *Snake) DecreaseMoveInterval(x int) {
	s.moveInterval = max(s.moveInterval-x, s.minMoveInterval)
}

func (s *Snake) CheckCollisionsWithSelf() {
	for i := 1; i < len(s.Body); i++ {
		if s.Body[i].Position == s.Body[0].Position {
			s.IsAlive = false
			break
		}
	}
}
