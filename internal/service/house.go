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
