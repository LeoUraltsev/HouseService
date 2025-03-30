package service

import (
	"context"

	"github.com/LeoUraltsev/HauseService/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	InsertUser(ctx context.Context, user models.User) error
}

type AuthService struct {
	AuthRepo AuthRepo
}

func NewAuthService(authRepo AuthRepo) *AuthService {
	return &AuthService{
		AuthRepo: authRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, user models.User) (string, error) {

	return "", nil
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
