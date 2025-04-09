package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/LeoUraltsev/HouseService/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	Type         string
}

const usersTableName = "users"

func (s *Storage) InsertUser(ctx context.Context, user models.User) error {
	const op = "storage.postgers.InsertUser"
	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", user.ID.String()),
		slog.String("user_email", user.Email),
	)

	log.Info("attempting insert user")

	u := ConvertToPGUser(&user)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Insert(usersTableName).
		Columns("id", "email", "password_hash", "type").
		Values(u.ID, u.Email, u.PasswordHash, u.Type).
		ToSql()

	if err != nil {
		return err
	}

	var pgErr *pgconn.PgError

	_, err = s.Pool.Exec(ctx, query, args...)
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		log.Error("user already exists")
		return models.ErrUserAlreadyExists
	}
	if err != nil {
		log.Error("failed insert user", slog.String("err", err.Error()))
		return err
	}

	log.Info("sucess insert user")
	return nil
}

func (s *Storage) SelectUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	const op = "storage.postgers.SelectUserByID"
	var u User
	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", id.String()),
	)

	log.Info("attempting getting user by id from db")

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("*").From(usersTableName).Where("id = ?", id).
		ToSql()

	log.Debug("creating sql query", slog.String("q", query))
	if err != nil {
		log.Error("sql string", slog.String("err", err.Error()))
		return nil, err
	}

	err = s.Pool.QueryRow(ctx, query, args...).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Type)
	if errors.Is(err, sql.ErrNoRows) {
		log.Warn("user not found")
		return nil, models.ErrUserNotFound
	}
	if err != nil {
		log.Error("failed getting user", slog.String("err", err.Error()))
		return nil, err
	}

	log.Info("success getting user from db")

	return &models.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		UserType:     models.UserType(u.Type),
	}, nil
}

func ConvertToPGUser(user *models.User) *User {
	return &User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Type:         string(user.UserType),
	}
}
