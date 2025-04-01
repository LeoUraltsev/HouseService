package service

import (
	"context"

	"github.com/LeoUraltsev/HauseService/internal/jwt"
	"github.com/LeoUraltsev/HauseService/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	InsertUser(ctx context.Context, user models.User) error
}

type AuthService struct {
	AuthRepo AuthRepo

	jwt *jwt.JWT
}

func NewAuthService(authRepo AuthRepo, jwt *jwt.JWT) *AuthService {
	return &AuthService{
		AuthRepo: authRepo,

		jwt: jwt,
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
