package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
)

type PostgresRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresRepository(connStr string, log *slog.Logger) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &PostgresRepository{
		db:     db,
		logger: log,
	}

	if err := repo.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Info("database created, schema initialized")
	return repo, nil
}

func (r *PostgresRepository) initSchema() error {
	query := `CREATE TABLE IF NOT EXISTS records(
    id SERIAL PRIMARY KEY,
    player_name VARCHAR(50) NOT NULL,
    score INT NOT NULL,
    time_in_seconds INT NOT_NULL,
    level_name VARCHAR(50) NOT NULL,
created_at TIMESTAMP WITH TIMEZONE DEFAULT CURRENT_TIMESTAMP
		)`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) SaveRecord(ctx context.Context, record Record) error {
	query := `INSERT INTO records (player_name, score, time_in_seconds, level_name) VALUES ($1, $2, $3, $4)`

	result, err := r.db.ExecContext(ctx, query, record.playerName, record.score, int(record.time.Seconds()), record.levelName)
	if err != nil {
		r.logger.Error("failed to save record", "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		r.logger.Warn("no rows were affected during saving record")
		return nil
	}

	r.logger.Info("record saved successfully", "player", record.playerName, "score", record.score, "time_in_sec", record.time, "level", record.levelName)
	return nil

}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
