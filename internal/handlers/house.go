package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/LeoUraltsev/HouseService/internal/gen"
	mv "github.com/LeoUraltsev/HouseService/internal/middleware"
	"github.com/LeoUraltsev/HouseService/internal/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type HouseService interface {
	HouseCreate(ctx context.Context, house models.House) (*models.House, error)
	HouseByID(ctx context.Context, id int) (*models.House, error)
	FlatsInHouseByID(ctx context.Context, id int, ut models.UserType) ([]models.Flat, error)
}

type HouseFlatsResponse struct {
	Flats []gen.Flat `json:"flats"`
}

// GetHouseId implements gen.ServerInterface.
func (h *Handler) GetHouseId(w http.ResponseWriter, r *http.Request, id gen.HouseId) {
	const op = "handlers.GetHouseId"

	reqID := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
		slog.Int("id", id),
	)

	log.Info("attempting getting flats in house")

	userType := r.Context().Value(mv.UserTypeContextKey).(models.UserType)

	fl, err := h.HouseService.FlatsInHouseByID(context.Background(), id, userType)
	if errors.Is(err, models.ErrFlatNotFound) {
		log.Warn("flats not found")
		render.JSON(w, r, []gen.Flat{})
	}

	if err != nil {
		log.Error("internal error", slog.String("err", err.Error()))
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("success getting flats")

	flats := make([]gen.Flat, len(fl))
	for i, v := range fl {
		flats[i] = gen.Flat{
			HouseId: v.HouseID,
			Id:      v.ID,
			Price:   gen.Price(v.Price),
			Rooms:   gen.Rooms(v.Rooms),
			Status:  gen.Status(v.Status),
		}
	}

	render.JSON(w, r, HouseFlatsResponse{
		Flats: flats,
	})

}

// PostHouseCreate implements gen.ServerInterface.
func (h *Handler) PostHouseCreate(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostHouseCreate"
	var req gen.PostHouseCreateJSONRequestBody

	reqID := middleware.GetReqID(r.Context())

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)

	log.Info("attempting create house")

	userType := r.Context().Value(mv.UserTypeContextKey).(models.UserType)
	if userType != models.Moderator {
		log.Warn("unautorized", slog.String("user_type", string(userType)))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Warn("invilide json", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := validateHouse(gen.PostHouseCreateJSONBody(req))
	if err != nil {
		log.Warn("validation params", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	house, err := h.HouseService.HouseCreate(context.Background(), models.House{
		Address:   req.Address,
		Year:      uint(req.Year),
		Developer: req.Developer,
	})

	if err != nil {
		log.Error("internal error", slog.String("err", err.Error()))
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("create hause", slog.Int("hause_id", house.UID))

	render.JSON(w, r, gen.House{
		Address:   house.Address,
		CreatedAt: &house.CreatedAt,
		Developer: house.Developer,
		Id:        house.UID,
		UpdateAt:  &house.LastFlatAddAt,
		Year:      gen.Year(house.Year),
	})
}

// PostHouseIdSubscribe implements gen.ServerInterface.
func (h *Handler) PostHouseIdSubscribe(w http.ResponseWriter, r *http.Request, id gen.HouseId) {
	panic("unimplemented")
}
