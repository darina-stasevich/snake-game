package storage

import (
	"context"
	"time"
)

type Record struct {
	ID         int64
	playerName string
	score      int
	time       time.Duration
	levelName  string
	createdAt  time.Time
}

type Repository interface {
	SaveRecord(ctx context.Context, record Record) error
	Close() error
}
