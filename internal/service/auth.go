package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/LeoUraltsev/HouseService/internal/jwt"
	"github.com/LeoUraltsev/HouseService/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	InsertUser(ctx context.Context, user models.User) error
	SelectUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type AuthService struct {
	AuthRepo AuthRepo

	jwt *jwt.JWT
	log *slog.Logger
}

func NewAuthService(authRepo AuthRepo, jwt *jwt.JWT, log *slog.Logger) *AuthService {
	return &AuthService{
		AuthRepo: authRepo,

		jwt: jwt,
		log: log,
	}
}

func (s *AuthService) DummyLogin(ctx context.Context, userType models.UserType) (string, error) {
	uid := generateUUID()
	token, err := s.jwt.NewToken(uid, userType)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, user models.User) (string, error) {
	const op = "service.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", user.ID.String()),
	)

	log.Info("attempting login")

	u, err := s.AuthRepo.SelectUserByID(ctx, user.ID)
	if errors.Is(err, models.ErrUserNotFound) {
		log.Error("user not found")
		return "", err
	}
	if err != nil {
		log.Error("failed getting user by id", slog.String("err", err.Error()))
		return "", err
	}

	log.Info("received user by id")

	if !(strings.EqualFold(u.ID.String(), user.ID.String()) && equalPasswordAndPasswordHash(user.PasswordHash, u.PasswordHash)) {
		log.Warn("invalid credentials")
		return "", models.ErrInvalideCredentials
	}

	log.Info("success login", slog.String("user_type", string(u.UserType)))
	log.Info("attempting create token", slog.String("user_type", string(u.UserType)))
	token, err := s.jwt.NewToken(u.ID, u.UserType)
	if err != nil {
		log.Error("create token", slog.String("err", err.Error()), slog.String("user_type", string(u.UserType)))
		return "", err
	}
	log.Info("success create token", slog.String("user_type", string(u.UserType)))

	return token, nil
}

func (s *AuthService) Register(ctx context.Context, user models.User) (*uuid.UUID, error) {
	const op = "service.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_email", user.Email),
	)

	hp, err := generatePasswordHash(user.PasswordHash)
	if err != nil {
		log.Error(
			"failed generate password hash",
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	uuid := generateUUID()

	u := models.User{
		ID:           uuid,
		Email:        user.Email,
		PasswordHash: hp,
		UserType:     user.UserType,
	}

	err = s.AuthRepo.InsertUser(ctx, u)
	if errors.Is(err, models.ErrUserAlreadyExists) {
		log.Warn("user already exists")
		return nil, err
	}
	if err != nil {
		log.Error("failed registration user", slog.String("err", err.Error()))
		return nil, err
	}

	return &uuid, nil
}

func generateUUID() uuid.UUID {
	return uuid.New()
}

func generatePasswordHash(pass string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(p), nil
}

func equalPasswordAndPasswordHash(pass string, passHash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(pass)); err != nil {
		return false
	}
	return true
}
