package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LeoUraltsev/HouseService/internal/models"
	"github.com/google/uuid"
)

type FlatRepo interface {
	InsertFlat(ctx context.Context, flat models.Flat) (*models.Flat, error)
	UpdateStatusFlat(ctx context.Context, flatID int, newStatus models.Status) (*models.Flat, error)
}

type FlatProviderRepo interface {
	SelectFlatByID(ctx context.Context, flatID int) (*models.Flat, error)
}

type ModerationRepo interface {
	SelectModeratorByFlatID(ctx context.Context, flatID int) (uuid.UUID, error)
	InsertModeratorForFlat(ctx context.Context, moderatorID uuid.UUID, flatID int) error
}

type FlatService struct {
	FlatRepo         FlatRepo
	FlatProviderRepo FlatProviderRepo
	ModerationRepo   ModerationRepo
	log              *slog.Logger
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
func (f *FlatService) FlatUpdate(ctx context.Context, flatId int, moderatorID uuid.UUID, newStatus models.Status) (*models.Flat, error) {
	/*
		1. Получить текуший статус
		2. Если текущий статус created
			2.1 Проверяем что пришел статус on moderation
			2.2 Меняем статус, квартиру добавляем в таблицу модерации
		3. Если текущий статус on moderation
			2.1 Проверяем что пришел запрос именно от модератора который переводил в статус on moderation
			2.2 Меняем статус
	*/

	flat, err := f.FlatProviderRepo.SelectFlatByID(ctx, flatId)
	if err != nil {
		return nil, err
	}

	switch flat.Status {
	case models.Created:
		switch newStatus {
		case models.OnModeration:
			nflat, err := f.FlatRepo.UpdateStatusFlat(ctx, flatId, newStatus)
			if err != nil {
				return nil, err
			}
			err = f.ModerationRepo.InsertModeratorForFlat(ctx, moderatorID, flatId)
			if err != nil {
				//todo: нужно продумать ретрай так как статус обновлен но модератор не добавлен,
				//либо отдельный метод который будет делать 2 действия через транзакцию
				return nil, fmt.Errorf("failed insert moderator for flat")
			}
			return nflat, nil
		default:
			return nil, fmt.Errorf("can't change the status without on moderation")
		}
	default:
		mID, err := f.ModerationRepo.SelectModeratorByFlatID(ctx, flatId)
		if err != nil {
			return nil, err
		}

		if mID != moderatorID {
			return nil, fmt.Errorf("flat checked by another moderator")
		}

		nflat, err := f.FlatRepo.UpdateStatusFlat(ctx, flatId, newStatus)
		if err != nil {
			return nil, err
		}
		return nflat, nil
	}
}
