package core

import "math"

type GameState int

const (
	MainMenuState GameState = iota
	GamePlayingState
	GamePauseState
	GameOverState
	LevelCreateState
	BestScoresState
)

type Position struct {
	X, Y int
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func GetDirection(first Position, second Position) Direction {
	if first.X == second.X {
		if first.Y < second.Y {
			return Up
		} else {
			return Down
		}
	} else {
		if first.X < second.X {
			return Left
		} else {
			return Right
		}
	}
}

func DirectionToRotationAngle(direction Direction) (rotation float64) {
	switch direction {
	case Right:
		rotation = 0
	case Left:
		rotation = math.Pi
	case Up:
		rotation = 3 * math.Pi / 2
	case Down:
		rotation = math.Pi / 2
	}
	return
}

func CornerToRotationAngle(oldDirection Direction, newDirection Direction) (rotation float64) {
	switch oldDirection {
	case Right:
		{
			if newDirection == Up {
				rotation = 0
			} else {
				rotation = 3 * math.Pi / 2
			}
		}
	case Left:
		if newDirection == Up {
			rotation = math.Pi / 2
		} else {
			rotation = math.Pi
		}
	case Up:
		if newDirection == Left {
			rotation = 3 * math.Pi / 2
		} else {
			rotation = math.Pi
		}
	case Down:
		if newDirection == Left {
			rotation = 0
		} else {
			rotation = math.Pi / 2
		}
	}
	return
}
