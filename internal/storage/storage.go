package storage

import (
	"context"
	"time"
)

type Record struct {
	ID         int64
	PlayerName string
	Score      int
	Time       time.Duration
	LevelName  string
	CreatedAt  time.Time
}

type Repository interface {
	SaveRecord(ctx context.Context, record Record) error
	Close() error
}
