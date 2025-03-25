package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LeoUraltsev/HauseService/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/pgtype"
)

type HouseDTO struct {
	UID           int64
	Address       string
	Year          uint
	Developer     string
	CreatedAt     pgtype.Timestamp
	LastFlatAddAt pgtype.Timestamp
}

func (p *Storage) InsertHouse(ctx context.Context, house models.House) (*models.House, error) {
	const op = "storage.postgres.InsertHouse"

	var h HouseDTO

	log := p.log.With(slog.String("op", op))

	log.Debug("attemping adding insert house in db")

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.
		Insert("house").
		Columns("address", "year", "developer", "created_at").
		Values(house.Address, house.Year, house.Developer, time.Now().UTC()).
		Suffix(`RETURNING *`).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed creating query sql", slog.String("err", err.Error()))
		return nil, err
	}

	log.Debug(query)

	err = p.Pool.QueryRow(ctx, query, args...).Scan(&h.UID, &h.Address, &h.Year, &h.Developer, &h.CreatedAt, &h.LastFlatAddAt)
	if err != nil {
		log.Error("failed adding new house", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed adding house")
	}

	log.Info("adding new house in db", slog.Int64("id", h.UID))
	return &models.House{
		UID:           h.UID,
		Address:       h.Address,
		Year:          h.Year,
		Developer:     h.Developer,
		CreatedAt:     h.CreatedAt.Time,
		LastFlatAddAt: h.LastFlatAddAt.Time,
	}, nil
}

func (p *Storage) SelectHouseByID(ctx context.Context, id int64) (*models.House, error) {
	const op = "storage.postgres.SelectHouseByID"

	log := p.log.With(slog.String("op", op))
	log.Debug("attempting getting house", slog.Int64("id", id))

	var h HouseDTO
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, _, err := psql.Select("*").From("house").Where("id=$1").ToSql()
	if err != nil {
		log.Error("failed create sql query", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed getting house: %v", err)
	}
	if err := p.Pool.QueryRow(ctx, query, id).Scan(&h.UID, &h.Address, &h.Year, &h.Developer, &h.CreatedAt, &h.LastFlatAddAt); err != nil {
		log.Error("failed getting house by id", slog.Int64("id", id), slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed getting house: %v", err)
	}

	log.Info("success getting house", slog.Int64("id", id))

	return &models.House{
		UID:           h.UID,
		Address:       h.Address,
		Year:          h.Year,
		Developer:     h.Developer,
		CreatedAt:     h.CreatedAt.Time,
		LastFlatAddAt: h.LastFlatAddAt.Time,
	}, nil
}
