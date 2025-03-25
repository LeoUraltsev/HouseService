package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/LeoUraltsev/HauseService/internal/models"
	sq "github.com/Masterminds/squirrel"
)

type Flat struct {
	ID      int64
	HouseID int64
	Price   uint
	Rooms   uint
	Status  string
}

func (s *Storage) InsertFlat(ctx context.Context, flat models.Flat) (*models.Flat, error) {
	const op = "storage.postgres.InsertFlat"

	log := s.log.With(slog.String("op", op))

	log.Debug("attempting adding flat")

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	fpg := ConvertToPGFlat(&flat)

	query, args, err := psql.Insert("flat").Columns("house_id", "price", "rooms", "status").
		Values(fpg.HouseID, fpg.Price, fpg.Rooms, fpg.Status).Prefix("RETURNING *").ToSql()
	if err != nil {
		return nil, err
	}
	log.Debug("generate sql", slog.String("query", query))

	if err := s.Pool.QueryRow(ctx, query, args...).Scan(&fpg); err != nil {
		return nil, err
	}
	t := time.Now().UTC()
	query, _, err = psql.Update("hause").Where("id = $1").Set("last_flat_add_at", t).ToSql()
	if err != nil {
		return nil, err
	}
	log.Debug("generate sql", slog.String("query", query))

	_, err = s.Pool.Exec(ctx, query, fpg.HouseID)
	if err != nil {
		log.Error(
			"failed update last_flat_add_at in table flat",
			slog.String("err", err.Error()),
			slog.Int64("flat_id", fpg.ID),
			slog.Int64("house_id", fpg.HouseID),
		)
	}

	return ConvertFromPGFlat(fpg), nil
}

func ConvertToPGFlat(flat *models.Flat) *Flat {
	return &Flat{
		ID:      flat.ID,
		HouseID: flat.HouseID,
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  flat.Status.String(),
	}
}

func ConvertFromPGFlat(flat *Flat) *models.Flat {
	var status models.Status

	return &models.Flat{
		ID:      flat.ID,
		HouseID: flat.HouseID,
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  status.Status(flat.Status),
	}
}
