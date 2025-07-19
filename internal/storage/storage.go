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

func NewRecord(playerName string, score int, time time.Duration, levelName string, created_at time.Time) *Record {
	if playerName == " " {
		playerName = "undefined"
	}
	return &Record{
		PlayerName: playerName,
		Score:      score,
		Time:       time,
		LevelName:  levelName,
		CreatedAt:  created_at,
	}
}

type Filter struct {
	playerNamePrefix string
	levelName        string
	isScoreAsc       bool
	isTimeAsc        bool
	playersMaxNumber int
}

func NewFilter(playerNamePrefix, levelName string, isScoreAsc, isTimeAsc bool, playersMaxNumber int) *Filter {
	return &Filter{
		playerNamePrefix: playerNamePrefix,
		levelName:        levelName,
		isScoreAsc:       isScoreAsc,
		isTimeAsc:        isTimeAsc,
		playersMaxNumber: playersMaxNumber,
	}
}

type Repository interface {
	SaveRecord(ctx context.Context, record *Record) error
	GetTopRecords(ctx context.Context, filter Filter) ([]Record, error)
	Close() error
}
