package models

import "github.com/google/uuid"

type UserType string

const (
	Client    UserType = "client"
	Moderator UserType = "moderator"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	UserType     UserType
}
