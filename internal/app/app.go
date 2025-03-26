package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/LeoUraltsev/HauseService/internal/config"
	"github.com/LeoUraltsev/HauseService/internal/models"
	"github.com/LeoUraltsev/HauseService/internal/storage/postgres"
	"github.com/google/uuid"
)

func Run(log *slog.Logger, cfg *config.Config) error {

	db, err := postgres.New(context.Background(), cfg.PostgresURLConnection, log)
	if err != nil {
		return err
	}
	defer db.Close()

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

	log.Info("new house 1", slog.Any("house", h))

	_, err = db.InsertFlat(context.Background(), models.Flat{
		HouseID: 2,
		Price:   12,
		Rooms:   2,
		Status:  models.Created,
	})

	if err != nil {
		log.Error("error", slog.Any("err", err))
		os.Exit(1)
	}

	db.UpdateStatusFlat(context.Background(), 2, models.OnModeration)

	db.InsertUser(context.Background(), models.User{
		ID:           uuid.New(),
		Email:        "emailtest@email.ru",
		PasswordHash: "asdasdasdq123",
		UserType:     models.Moderator,
	})

	return nil
}
