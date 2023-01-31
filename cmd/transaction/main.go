package main

import (
	"bwg_test/internal/transaction"
	"bwg_test/internal/transaction/config"
	"bwg_test/internal/transaction/http"
	"bwg_test/internal/transaction/http/handlers"
	"bwg_test/internal/transaction/storage"
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func main() {
	out := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.StampMilli,
	}

	logger := zerolog.New(out).With().Caller().Logger().With().Timestamp().Logger()

	cfg, err := config.Parse()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")
	}

	conntectionStr := "user=admin password=admin dbname=postgres sslmode=disable"

	conn, err := sqlx.Connect("postgres", conntectionStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect db")
	}

	db := storage.New(conn)

	svc := transaction.New(context.Background(), logger, db)
	handler := handlers.New(logger, svc)

	server := http.New(cfg.Server.Host, handler)

	if err = server.Run(); err != nil {
		return
	}
}
