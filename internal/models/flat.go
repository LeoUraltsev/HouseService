package models

type Status int

const (
	Created Status = iota
	Approved
	Declined
	OnModeration
)

type Flat struct {
	ID      int
	HouseID int
	Price   uint
	Rooms   uint
	Status  Status
}

var status = map[Status]string{
	Created:      "created",
	Approved:     "approved",
	Declined:     "declined",
	OnModeration: "on moderation",
}

func (s Status) String() string {
	return status[s]
}

func (s Status) Status(value string) Status {
	reverseMap := make(map[string]Status)
	for key, val := range status {
		reverseMap[val] = key
	}
	return reverseMap[value]
}
