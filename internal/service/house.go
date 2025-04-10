package service

import (
	"context"
	"log/slog"

	"github.com/LeoUraltsev/HouseService/internal/jwt"
	"github.com/LeoUraltsev/HouseService/internal/models"
)

type HouseRepo interface {
	InsertHouse(ctx context.Context, house models.House) (*models.House, error)
	SelectHouseByID(ctx context.Context, id int) (*models.House, error)
	SelectFlatsInHouseByID(ctx context.Context, id int) ([]models.Flat, error)
}

type HouseService struct {
	HouseRepo HouseRepo

	jwt *jwt.JWT
	log *slog.Logger
}

func NewHouseService(h HouseRepo, j *jwt.JWT, log *slog.Logger) *HouseService {
	return &HouseService{
		HouseRepo: h,
		jwt:       j,
		log:       log,
	}
}

// HouseByID implements handlers.HouseService.
func (h *HouseService) HouseByID(ctx context.Context, id int) (*models.House, error) {
	return h.HouseRepo.SelectHouseByID(ctx, id)
}

// HouseCreate implements handlers.HouseService.
func (h *HouseService) HouseCreate(ctx context.Context, house models.House) (*models.House, error) {
	return h.HouseRepo.InsertHouse(ctx, house)
}

// todo: подумать как выделять память в слайс под результат, возможно стоит инициализировать его на половину длины ответа
// или в зависимости от типа пользоватяля, получать дома из базы только с нужными статусами
func (h *HouseService) FlatsInHouseByID(ctx context.Context, id int, ut models.UserType) ([]models.Flat, error) {

	f, err := h.HouseRepo.SelectFlatsInHouseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ut == models.Moderator {
		return f, nil
	}

	var res []models.Flat

	for _, v := range f {
		if v.Status == models.Approved {
			res = append(res, v)
		}
	}
	return res, nil
}
