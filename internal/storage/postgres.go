package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"time"
)

type PostgresRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresRepository(connStr string, log *slog.Logger) (Repository, error) {
	log.Debug("im in NewPostgresRepository")
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	log.Debug("2. im in NewPostgresRepository")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Важно вызывать cancel, чтобы освободить ресурсы контекста

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database within 5s: %w", err)
	}
	log.Debug("3. im in NewPostgresRepository")

	repo := &PostgresRepository{
		db:     db,
		logger: log,
	}

	log.Debug("4. im in NewPostgresRepository")

	if err := repo.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Debug("5. im in NewPostgresRepository")

	log.Info("database created, schema initialized")
	return repo, nil
}

func (r *PostgresRepository) initSchema() error {
	query := `CREATE TABLE IF NOT EXISTS records(
    id SERIAL PRIMARY KEY,
    player_name VARCHAR(50) NOT NULL,
    score INT NOT NULL,
    time_in_seconds INT NOT NULL,
    level_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) SaveRecord(ctx context.Context, record Record) error {
	query := `INSERT INTO records (player_name, score, time_in_seconds, level_name) VALUES ($1, $2, $3, $4)`

	result, err := r.db.ExecContext(ctx, query, record.PlayerName, record.Score, int(record.Time.Seconds()), record.LevelName)
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

	r.logger.Info("record saved successfully", "player", record.PlayerName, "score", record.Score, "time_in_sec", record.Time, "level", record.LevelName)
	return nil

}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
