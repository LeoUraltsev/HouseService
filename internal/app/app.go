package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/LeoUraltsev/HauseService/internal/config"
	"github.com/LeoUraltsev/HauseService/internal/gen"
	"github.com/LeoUraltsev/HauseService/internal/handlers"

	"github.com/LeoUraltsev/HauseService/internal/storage/postgres"
)

func Run(log *slog.Logger, cfg *config.Config) error {

	db, err := postgres.New(context.Background(), cfg.PostgresURLConnection, log)
	if err != nil {
		return err
	}
	defer db.Close()

	r := gen.Handler(handlers.New())

	http.ListenAndServe("localhost:10000", r)
	return nil
}
