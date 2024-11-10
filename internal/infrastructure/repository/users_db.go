package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"go-boilerplate/pkg/e"

	"github.com/jackc/pgx/v5/pgxpool"
)

var errNoRows = errors.New("no rows in result set")

type PostgresDB struct {
	conn    *pgxpool.Pool
	timeout time.Duration
}

func NewPostgresDB(pool *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{
		conn:    pool,
		timeout: 1 * time.Second,
	}
}

func (db *PostgresDB) InsertUser(ctx context.Context, tgId, tgChatId int) (userId int, err error) {
	query := `INSERT INTO users (telegram_id, telegram_chat_id) VALUES($1, $2) RETURNING id;`

	if err = db.conn.QueryRow(ctx, query, tgId, tgChatId).Scan(&userId); err != nil {
		return 0, e.Wrap("failed to execute insert user query", err)
	}

	return userId, nil
}

func (db *PostgresDB) GetUserByTelegramId(ctx context.Context, tgId int) (userId int, err error) {
	query := `SELECT id FROM users WHERE telegram_id = $1;`

	if err = db.conn.QueryRow(ctx, query, tgId).Scan(&userId); errors.Is(err, pgx.ErrNoRows) {
		return 0, errNoRows
	}

	if err != nil {
		return 0, e.Wrap("failed to query user by telegram id", err)
	}

	return userId, nil
}

func (db *PostgresDB) GetTelegramChatIdByUserId(ctx context.Context, userId int) (tgChatId int, err error) {
	query := `SELECT telegram_chat_id FROM users WHERE id = $1;`

	if err = db.conn.QueryRow(ctx, query, userId).Scan(&tgChatId); errors.Is(err, pgx.ErrNoRows) {
		return 0, errNoRows
	}

	if err != nil {
		return 0, e.Wrap("failed to query telegram chat by user id", err)
	}

	return tgChatId, nil
}
