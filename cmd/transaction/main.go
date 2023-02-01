package main

import (
	"bwg_test/internal/transaction"
	"bwg_test/internal/transaction/config"
	"bwg_test/internal/transaction/http"
	"bwg_test/internal/transaction/http/handlers"
	"bwg_test/internal/transaction/storage"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
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

	connectionStr := ParseDBConnStr(
		cfg.DataBase.Login,
		cfg.DataBase.Password,
		cfg.DataBase.DBName,
		cfg.DataBase.SslMode,
	)

	conn, err := sqlx.Connect("postgres", connectionStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect db")
	}

	db := storage.New(conn)

	svc := transaction.New(context.Background(), logger, db)
	handler := handlers.New(logger, svc)

	server := http.New(cfg.Server.Host, handler)

	go func() {
		if err = server.Run(); err != nil {
			logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info().Msg("http server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("server shutdown error")
	}
}

func ParseDBConnStr(login, password, db, sslmode string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		login,
		password,
		db,
		sslmode,
	)
}
