package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/LeoUraltsev/HouseService/internal/models"
	sq "github.com/Masterminds/squirrel"
)

type Flat struct {
	ID      int
	HouseID int
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
		Values(fpg.HouseID, fpg.Price, fpg.Rooms, fpg.Status).Suffix("RETURNING *").ToSql()
	if err != nil {
		return nil, err
	}
	log.Debug("generate sql", slog.String("query", query))

	if err := s.Pool.QueryRow(ctx, query, args...).Scan(&fpg.ID, &fpg.HouseID, &fpg.Rooms, &flat.Price, &fpg.Status); err != nil {
		return nil, err
	}
	t := time.Now().UTC()
	query, args1, err := psql.Update("house").Where("id = ?", fpg.ID).Set("last_flat_add_at", t).ToSql()
	if err != nil {
		return nil, err
	}
	log.Debug("generate sql", slog.String("query", query))

	_, err = s.Pool.Exec(ctx, query, args1...)
	if err != nil {
		log.Error(
			"failed update last_flat_add_at in table hause",
			slog.String("err", err.Error()),
			slog.Int("flat_id", fpg.ID),
			slog.Int("house_id", fpg.HouseID),
		)
	}

	return ConvertFromPGFlat(fpg), nil
}

func (s *Storage) UpdateStatusFlat(ctx context.Context, flatID int, newStatus models.Status) (*models.Flat, error) {
	const op = "storage.postgres.UpdateStatusFlat"
	var f Flat

	log := s.log.With(
		slog.String("op", op),
		slog.Int("flat_id", flatID),
	)

	log.Info("attempting update status flat")

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.
		Update("flat").
		Where("id = ?", flatID).
		Set("status", newStatus.String()).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}
	log.Debug("generate sql", slog.String("query", query))

	if err := s.Pool.QueryRow(ctx, query, args...).Scan(&f.ID, &f.HouseID, &f.Price, &f.Rooms, &f.Status); err != nil {
		log.Error("failed change status", slog.String("err", err.Error()))

		return nil, err
	}

	log.Info("success change status", slog.String("new_status", f.Status))
	return ConvertFromPGFlat(&f), nil
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
