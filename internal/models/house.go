package models

import "time"

/*
уникальный номер дома
адрес
год постройки
застройщик (у 50% домов)
дата создания дома в базе
дата последнего добавления новой квартиры дома
*/
type House struct {
	UID           int64
	Address       string
	Year          uint
	Developer     string
	CreatedAt     time.Time
	LastFlatAddAt time.Time
}
