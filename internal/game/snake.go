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

func (s *Snake) Grow() {
	panic("implement me")
}

func (s *Snake) Update() {
	s.moveTimer++
	if s.moveTimer < s.moveInterval {
		return
	}
	s.moveTimer = 0
	s.move()
}

func (s *Snake) move() {
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

	if s.AteApple {
		s.AteApple = false
	} else {
		s.Body = s.Body[0 : len(s.Body)-1]
	}
}
