package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/hiteshwadhwani/go-rest/internal/config"
	"github.com/hiteshwadhwani/go-rest/internal/healthcheck"
	"github.com/hiteshwadhwani/go-rest/pkg/log"
	_ "github.com/lib/pq"
)

var Version = "1.0"

var flagConfig = flag.String("config", "./config/local.yml", "config file")

func main() {
	flag.Parse()

	logger := log.New().With(nil, "version", Version)

	loadedConfig, err := config.Load(*flagConfig, logger)

	if err != nil {
		logger.Errorf("failed to load config: %v", err)
		os.Exit(-1)
	}

	db, err := dbx.MustOpen("postgres", loadedConfig.DSN)

	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("failed to close database connection: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("failed to connect to database: %v", err)
		os.Exit(-1)
	}

	db.QueryLogFunc = logDBQuery(logger)
	db.ExecLogFunc = logDbExecution(logger)

	address := fmt.Sprintf(":%v", loadedConfig.ServerPort)

	hs := http.Server{
		Addr:    address,
		Handler: buildHandler(loadedConfig),
	}

	go routing.GracefulShutdown(&hs, 10*time.Second, logger.Infof)

	logger.Infof("server %v is running at %v", Version, address)

	if err := hs.ListenAndServe(); err != nil {
		logger.Errorf("failed to start server: %v", err)
		os.Exit(-1)
	}

}

func buildHandler(config *config.Config) http.Handler {
	router := routing.New()

	// register healthcheck handler
	healthcheck.RegisterHealthCheckHandler(router)

	return router
}

func logDBQuery(logger log.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

func logDbExecution(logger log.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}
