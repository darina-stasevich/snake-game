package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
	"strings"
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

func (r *PostgresRepository) SaveRecord(ctx context.Context, record *Record) error {
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

func (r *PostgresRepository) GetTopRecords(ctx context.Context, filter Filter) ([]Record, error) {
	baseQuery := `SELECT player_name, score, time_in_seconds, level_name, created_at FROM records`

	var whereClauses []string
	var args []interface{}
	argID := 1

	if filter.playerNamePrefix != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("player_name LIKE $%d", argID))
		args = append(args, filter.playerNamePrefix+"%")
		argID++
	}

	if filter.levelName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("level_name = $%d", argID))
		args = append(args, filter.levelName)
		argID++
	}

	query := baseQuery
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY score " + isAsc(filter.isScoreAsc) + ", time_in_seconds " + isAsc(filter.isTimeAsc)

	if filter.playersMaxNumber > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argID)
		args = append(args, filter.playersMaxNumber)
		argID++
	}

	r.logger.Info("made query", "query", query)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to execute query", "error", err)
		return nil, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var (
			playerName  string
			score       int
			time_in_sec int
			levelName   string
			created_at  time.Time
		)

		err = rows.Scan(&playerName, &score, &time_in_sec, &levelName, &created_at)
		if err != nil {
			r.logger.Error("failed to scan row", "err", err)
			return nil, err
		}

		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}

		record := NewRecord(playerName, score, time.Duration(time_in_sec)*time.Second, levelName, created_at)
		records = append(records, *record)
	}

	return records, nil
}

func isAsc(isAsc bool) string {
	if isAsc == true {
		return "ASC"
	} else {
		return "DESC"
	}
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
