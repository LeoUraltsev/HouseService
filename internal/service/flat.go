package service

import (
	"context"
	"log/slog"

	"github.com/LeoUraltsev/HouseService/internal/models"
)

type FlatRepo interface {
	InsertFlat(ctx context.Context, flat models.Flat) (*models.Flat, error)
	UpdateStatusFlat(ctx context.Context, flatID int, newStatus models.Status) (*models.Flat, error)
}

type FlatService struct {
	FlatRepo FlatRepo

	log *slog.Logger
}

func NewFlatService(f FlatRepo, l *slog.Logger) *FlatService {
	return &FlatService{
		FlatRepo: f,
		log:      l,
	}
}

// FlatCreate implements handlers.FlatService.
func (f *FlatService) FlatCreate(ctx context.Context, flat models.Flat) (*models.Flat, error) {
	return f.FlatRepo.InsertFlat(ctx, flat)
}

// FlatUpdate implements handlers.FlatService.
func (f *FlatService) FlatUpdate(ctx context.Context, flatId int, newStatus models.Status) (*models.Flat, error) {
	return f.FlatRepo.UpdateStatusFlat(ctx, flatId, newStatus)
}
