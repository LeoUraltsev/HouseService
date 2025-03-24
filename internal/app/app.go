package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/LeoUraltsev/HauseService/internal/config"
	"github.com/LeoUraltsev/HauseService/internal/models"
	"github.com/LeoUraltsev/HauseService/internal/storage/postgres"
)

func Run(log *slog.Logger, cfg *config.Config) error {

	db, err := postgres.New(context.Background(), cfg.PostgresURLConnection, log)
	if err != nil {
		return err
	}

	h, err := db.InsertHouse(context.Background(), models.House{
		Address:       "asd",
		Year:          2024,
		Developer:     "develop",
		CreatedAt:     time.Now(),
		LastFlatAddAt: time.Now(),
	})

	if err != nil {
		return err
	}

	log.Info("new house", slog.Any("house", h))

	return nil
}
