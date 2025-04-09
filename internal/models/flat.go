package models

import "errors"

type Status string

const (
	Created      Status = "created"
	Approved     Status = "approved"
	Declined     Status = "declined"
	OnModeration Status = "on moderation"
)

var (
	ErrFlatNotFound = errors.New("flat not found")
)

type Flat struct {
	ID      int
	HouseID int
	Price   uint
	Rooms   uint
	Status  Status
}
