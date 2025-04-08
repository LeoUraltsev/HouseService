package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LeoUraltsev/HouseService/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/pgtype"
)

type House struct {
	UID           int
	Address       string
	Year          uint
	Developer     string
	CreatedAt     pgtype.Timestamp
	LastFlatAddAt pgtype.Timestamp
}

func (p *Storage) InsertHouse(ctx context.Context, house models.House) (*models.House, error) {
	const op = "storage.postgres.InsertHouse"

	var h House

	log := p.log.With(slog.String("op", op))

	log.Debug("attemping adding insert house in db")

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.
		Insert("house").
		Columns("address", "year", "developer", "created_at", "last_flat_add_at").
		Values(house.Address, house.Year, house.Developer, time.Now().UTC(), time.Now().UTC()).
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

	log.Info("adding new house in db", slog.Int("id", h.UID))
	return &models.House{
		UID:           h.UID,
		Address:       h.Address,
		Year:          h.Year,
		Developer:     &h.Developer,
		CreatedAt:     h.CreatedAt.Time,
		LastFlatAddAt: h.LastFlatAddAt.Time,
	}, nil
}

func (p *Storage) SelectHouseByID(ctx context.Context, id int) (*models.House, error) {
	const op = "storage.postgres.SelectHouseByID"

	log := p.log.With(slog.String("op", op))
	log.Debug("attempting getting house", slog.Int("id", id))

	var h House
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, _, err := psql.Select("*").From("house").Where("id=$1").ToSql()
	if err != nil {
		log.Error("failed create sql query", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed getting house: %v", err)
	}
	if err := p.Pool.QueryRow(ctx, query, id).Scan(&h.UID, &h.Address, &h.Year, &h.Developer, &h.CreatedAt, &h.LastFlatAddAt); err != nil {
		log.Error("failed getting house by id", slog.Int("id", id), slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed getting house: %v", err)
	}

	log.Info("success getting house", slog.Int("id", id))

	return &models.House{
		UID:           h.UID,
		Address:       h.Address,
		Year:          h.Year,
		Developer:     &h.Developer,
		CreatedAt:     h.CreatedAt.Time,
		LastFlatAddAt: h.LastFlatAddAt.Time,
	}, nil
}
