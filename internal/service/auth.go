package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/LeoUraltsev/HauseService/internal/jwt"
	"github.com/LeoUraltsev/HauseService/internal/models"
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
	token, err := s.jwt.NewToken(userType)
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
	if err != nil {
		log.Error("getting user by id")
		return "", err
	}

	log.Info("received user by id")

	if !(strings.EqualFold(u.ID.String(), user.ID.String()) && equalPasswordAndPasswordHash(user.PasswordHash, u.PasswordHash)) {
		log.Warn("invalid credentials")
		return "", fmt.Errorf("invalid credentials")
	}

	log.Info("success login", slog.String("user_type", string(u.UserType)))
	log.Info("attempting create token", slog.String("user_type", string(u.UserType)))
	token, err := s.jwt.NewToken(u.UserType)
	if err != nil {
		log.Error("create token", slog.String("err", err.Error()), slog.String("user_type", string(u.UserType)))
		return "", err
	}
	log.Info("success create token", slog.String("user_type", string(u.UserType)))

	return token, nil
}
func (s *AuthService) Register(ctx context.Context, user models.User) (*uuid.UUID, error) {
	hp, err := generatePasswordHash(user.PasswordHash)
	uuid := generateUUID()
	if err != nil {
		return nil, err
	}
	u := models.User{
		ID:           uuid,
		Email:        user.Email,
		PasswordHash: hp,
		UserType:     user.UserType,
	}
	if err := s.AuthRepo.InsertUser(ctx, u); err != nil {
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
