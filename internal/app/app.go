package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pressly/goose/v3"

	"github.com/mathbdw/subscription-service/config"
	"github.com/mathbdw/subscription-service/internal/infrastructure/httpserver"
	"github.com/mathbdw/subscription-service/internal/infrastructure/observability/logger/zerolog"
	"github.com/mathbdw/subscription-service/internal/infrastructure/persistence/postgres"
	"github.com/mathbdw/subscription-service/internal/infrastructure/persistence/postgres/repositories"
	httpimp "github.com/mathbdw/subscription-service/internal/interfaces/http"
	"github.com/mathbdw/subscription-service/internal/interfaces/observability"
	"github.com/mathbdw/subscription-service/internal/usecases/subscription"
)

// initLogger - initializing logger
func initLogger(cfg *config.Config) observability.Logger {
	logger := zerolog.New(cfg)
	logger.Debug("app.initLogger: config running ", map[string]any{"config": cfg})

	return logger
}

// initPostgres - initializing postgres
func initPostgres(cfg *config.Config, logger observability.Logger) *postgres.Postgres {
	pg, err := postgres.New(
		logger,
		postgres.Dsn(cfg.Database),
		postgres.Driver(cfg.Database.Driver),
		postgres.MaxOpenConns(cfg.Database.MaxOpenConns),
		postgres.MaxIdleConns(cfg.Database.MaxIdleConns),
		postgres.ConnMaxIdleTime(cfg.Database.ConnMaxIdleTime),
		postgres.ConnMaxLifetime(cfg.Database.ConnMaxLifetime),
	)
	if err != nil {
		logger.Fatal("app.initPostgres: PG new", map[string]any{"error": err})
	}

	return pg
}

// applyMigration - apply migration
func applyMigration(cfg *config.Config, pg *postgres.Postgres, logger observability.Logger) {
	if err := goose.Up(pg.Sqlx.DB, cfg.Database.Migrations); err != nil {
		logger.Fatal("app.applyMigration: failed migration", map[string]any{"err": err})
	}
}

// RunApp - run application
func RunApp(cfg *config.Config) {
	logger := initLogger(cfg)
	pg := initPostgres(cfg, logger)
	defer pg.Sqlx.Close()

	applyMigration(cfg, pg, logger)

	repoSub := repositories.NewUserRepository(pg.Sqlx, pg.Builder, logger)
	usSub := subscription.NewSubscriptionUsecase(repoSub, logger)

	httpServer := httpserver.New(
		httpserver.Address(cfg.Rest.Host, cfg.Rest.Port),
		httpserver.Prefork(cfg.Rest.Prefork),
		httpserver.ReadTimeout(cfg.Rest.ReadTimeout),
		httpserver.WriteTimeout(cfg.Rest.WriteTimeout),
		httpserver.ShutdownTimeout(cfg.Rest.ShutdownTimeout),
	)
	httpimp.NewRouter(httpServer.App, &cfg.Rest, usSub, logger)

	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Error("app.RunApp:", map[string]any{"signal": s.String()})
	case err := <-httpServer.Notify():
		logger.Error("app.RunApp: httpServer.Notify", map[string]any{"error": err.Error()})
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		logger.Error("app.RunApp: httpServer.Shutdown", map[string]any{"error": err.Error()})
	} else {
		logger.Error("app.RunApp: httpServer shutting down...", nil)
	}
}
