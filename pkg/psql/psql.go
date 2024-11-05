package psql

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/pressly/goose/v3"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3/database"
	"github.com/pressly/goose/v3/lock"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect возвращает настроенный pgxpool к БД
func Connect(ctx context.Context, opts ...OptionFunc) (*pgxpool.Pool, error) {
	options := defaultOptions

	for _, opt := range opts {
		opt(&options)
	}

	connString := (&url.URL{
		Scheme: "postgres",
		User:   getUserInfo(options.user, options.password),
		Host:   fmt.Sprintf("%s:%d", options.host, options.port),
		Path:   options.database,
	}).String()

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	poolConfig.MinConns = options.maxIdleConns
	poolConfig.MaxConns = options.maxOpenConns
	poolConfig.MaxConnIdleTime = options.connMaxIdleTime
	poolConfig.MaxConnLifetime = options.connMaxLifetime

	var db *pgxpool.Pool
	db, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if options.connectionWaiting.enabled {
		if err = options.connectionWaiting.pingDB(ctx, db); err != nil {
			return nil, fmt.Errorf("failed to ping database: %w", err)
		}
	}

	if options.migrations != nil {
		migrationConnString := (&url.URL{
			Scheme: "postgres",
			User:   getUserInfo(options.userAdmin, options.passwordAdmin),
			Host:   fmt.Sprintf("%s:%d", options.host, options.port),
			Path:   options.database,
		}).String()

		config, err := pgx.ParseConfig(migrationConnString)
		if err != nil {
			return nil, err
		}

		stdDb := stdlib.OpenDB(*config)
		stdDb.SetMaxIdleConns(0)
		defer func() { _ = stdDb.Close() }()

		err = stdDb.Ping()
		if err != nil {
			return nil, err
		}

		locker, err := lock.NewPostgresSessionLocker()
		if err != nil {
			return nil, err
		}

		g, err := goose.NewProvider(
			database.DialectPostgres,
			stdDb,
			options.migrations,
			goose.WithSessionLocker(locker),
			goose.WithAllowOutofOrder(true),
		)
		if err != nil {
			return nil, err
		}

		if options.logger != nil {
			version, err := g.GetDBVersion(ctx)
			if err != nil {
				return nil, err
			}

			options.logger.Info(fmt.Sprintf("db migrations version %d", version))
		}

		res, err := g.Up(ctx)
		if err != nil {
			return nil, err
		}

		if options.logger != nil {
			for _, i := range res {
				options.logger.Info(fmt.Sprintf("migration file %s, duration %s", i.Source.Path, i.Duration.String()))
			}
		}
	}

	return db, nil
}

func getUserInfo(user, password string) *url.Userinfo {
	if user != "" {
		if password != "" {
			return url.UserPassword(user, password)
		}
		return url.User(user)
	}
	return nil
}
