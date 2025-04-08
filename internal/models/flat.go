package models

type Status string

const (
	Created      Status = "created"
	Approved     Status = "approved"
	Declined     Status = "declined"
	OnModeration Status = "on moderation"
)

type Flat struct {
	ID      int
	HouseID int
	Price   uint
	Rooms   uint
	Status  Status
}
