package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	httpv1 "go-boilerplate/internal/controller/http_v1"
	tg "go-boilerplate/internal/entity/telegram"
	"go-boilerplate/internal/infrastructure/repository"

	"go.uber.org/zap"

	"go-boilerplate/pkg/psql"

	"github.com/wanomir/d"
	"github.com/wanomir/e"
	"github.com/wanomir/l"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	exitStatusOk     = 0
	exitStatusFailed = 1
)

type App struct {
	config *Config
	logger *zap.Logger

	ctx     context.Context
	errChan chan error

	server *http.Server
	tg     *tg.Telegram
	http   *httpv1.HttpController
}

func NewApp() (*App, error) {
	a := new(App)

	if err := a.init(); err != nil {
		return nil, e.Wrap("failed to init app", err)
	}

	return a, nil
}

func (a *App) Run() (exitCode int) {
	defer a.recoverFromPanic(&exitCode)
	var err error

	ctx, stop := signal.NotifyContext(a.ctx, os.Interrupt, os.Kill)
	defer stop()

	// run telegram
	go a.tg.Run(ctx)

	// listen for incoming api requests
	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errChan <- err
		}
	}()

	go d.Run(a.config.Debug.ServerAddr)

	select {
	case err = <-a.errChan:
		a.logger.Error(e.Wrap("fatal error, service shutdown", err).Error())
		exitCode = exitStatusFailed
	case <-ctx.Done():
		a.logger.Info("service shutdown")
	}

	return exitStatusOk
}

func (a *App) init() (err error) {
	// config
	if err = a.readConfig(); err != nil {
		return e.Wrap("failed to read config", err)
	}

	a.ctx = context.Background()
	a.errChan = make(chan error)

	l.BuildLogger(a.config.Log.Level)
	a.logger = l.Logger()

	// database
	pool, err := psql.Connect(a.ctx,
		psql.WithHost(a.config.PG.Host),
		psql.WithDatabase(a.config.PG.Database),
		psql.WithUser(a.config.PG.User),
		psql.WithPassword(a.config.PG.Password),
		psql.WithUserAdmin(a.config.PG.UserAdmin),
		psql.WithPasswordAdmin(a.config.PG.PasswordAdmin),
		psql.WithMigrations(os.DirFS("db/migrations")),
		psql.WithLogger(a.logger),
	)
	if err != nil {
		return e.Wrap("failed to init db", err)
	}

	_ = repository.NewPostgresDB(pool)

	// telegram service
	if a.tg, err = tg.NewTelegram(a.config.TG.Token, a.logger); err != nil {
		return e.Wrap("failed to init telegram", err)
	}

	// http server
	a.server = &http.Server{
		Addr:         a.config.Target.Addr,
		Handler:      a.routes(),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return nil
}

func (a *App) readConfig() (err error) {
	a.config = new(Config)
	if err = cleanenv.ReadEnv(a.config); err != nil {
		return err
	}

	return nil
}

func (a *App) recoverFromPanic(exitCode *int) {
	if panicErr := recover(); panicErr != nil {
		a.logger.Error(fmt.Sprintf("recover from panic: %v, stacktrace: %s", panicErr, string(debug.Stack())))
		*exitCode = exitStatusFailed
	}
}
