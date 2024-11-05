package psql

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var defaultPingOptions = pingOptions{
	enabled:      true,
	tickInterval: time.Second,
	deadline:     time.Second * 30,
}

type pingOptions struct {
	enabled      bool
	tickInterval time.Duration
	deadline     time.Duration
}

// OptionPingFunc - тип опции
type OptionPingFunc func(*psqlOptions)

// WithTickInterval - устанавливает время интервала, между попытками пинга БД
func WithTickInterval(d time.Duration) OptionPingFunc {
	return func(o *psqlOptions) {
		o.connectionWaiting.tickInterval = d
	}
}

// WithDeadline - устанавливает максимальное время для пинга БД
func WithDeadline(d time.Duration) OptionPingFunc {
	return func(o *psqlOptions) {
		o.connectionWaiting.deadline = d
	}
}

// pingDB Ждем пока соединение с бд установится, пингуем раз в tickInterval секунд до истечения deadLine секунд.
func (p *pingOptions) pingDB(ctx context.Context, db *pgxpool.Pool) error {
	if db == nil {
		return errors.New("db is nil")
	}

	// Пингуем сразу, что бы не ждать впустую
	if err := db.Ping(ctx); err == nil {
		return nil
	}

	ticker := time.NewTicker(p.tickInterval)
	defer ticker.Stop()
	timer := time.NewTimer(p.deadline)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			if err := db.Ping(ctx); err == nil {
				return nil
			} else {
				slog.Debug("ping", "err", err)
			}

		case <-timer.C:
			return db.Ping(ctx)
		}
	}
}
