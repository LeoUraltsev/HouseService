package models

import (
	"errors"
	"time"
)

/*
уникальный номер дома
адрес
год постройки
застройщик (у 50% домов)
дата создания дома в базе
дата последнего добавления новой квартиры дома
*/

var (
	ErrHouseNotFound = errors.New("house not found")
)

type House struct {
	UID           int
	Address       string
	Year          uint
	Developer     *string
	CreatedAt     time.Time
	LastFlatAddAt time.Time
}
