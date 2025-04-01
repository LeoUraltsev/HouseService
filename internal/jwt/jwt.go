package jwt

import (
	"fmt"
	"time"

	"github.com/LeoUraltsev/HauseService/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secret   []byte
	duration time.Duration
}

type CustomClaims struct {
	UserType models.UserType `json:"user_type"`
	jwt.RegisteredClaims
}

func New(duration time.Duration, secret []byte) *JWT {
	return &JWT{
		secret:   secret,
		duration: duration,
	}
}

func (j *JWT) NewToken(ut models.UserType) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserType: ut,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(j.duration),
			},
		},
	})

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("failed token generated: %w", err)
	}

	return tokenString, nil
}

func (j *JWT) ParseToken(token string) (*CustomClaims, error) {

	t, err := jwt.ParseWithClaims(
		token,
		&CustomClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return j.secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := t.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("unknown claims type")
	}

	return claims, nil
}
