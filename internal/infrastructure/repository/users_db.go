package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	conn *pgxpool.Pool
}

func NewPostgresDB(pool *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{
		conn: pool,
	}
}

func (db *PostgresDB) InsertRamData(ctx context.Context, ts time.Time, host string, ramData float64) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `INSERT INTO host_data.ram_data (time, host, ram_usage) VALUES($1, $2, $3);`

	_, err = db.conn.Exec(ctx, query, ts, host, ramData)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to execute insert cpu_data query: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
