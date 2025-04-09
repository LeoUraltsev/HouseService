package models

import (
	"errors"

	"github.com/google/uuid"
)

type UserType string

const (
	Client    UserType = "client"
	Moderator UserType = "moderator"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	UserType     UserType
}
