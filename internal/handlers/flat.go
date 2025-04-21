package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/LeoUraltsev/HouseService/internal/gen"
	mv "github.com/LeoUraltsev/HouseService/internal/middleware"
	"github.com/LeoUraltsev/HouseService/internal/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type FlatService interface {
	FlatCreate(ctx context.Context, flat models.Flat) (*models.Flat, error)
	FlatUpdate(ctx context.Context, flatId int, moderatorID uuid.UUID, newStatus models.Status) (*models.Flat, error)
}

// PostFlatCreate implements gen.ServerInterface.
func (h *Handler) PostFlatCreate(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostFlatCreate"
	reqID := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)

	log.Info("attempting create flat")

	var f gen.PostFlatCreateJSONRequestBody
	if err := render.DecodeJSON(r.Body, &f); err != nil {
		log.Warn("incorrect json body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug("success decode json", slog.Any("req", f))

	flat, err := h.FlatService.FlatCreate(context.Background(), models.Flat{
		HouseID: f.HouseId,
		Price:   uint(f.Price),
		Rooms:   uint(*f.Rooms),
		Status:  models.Created,
	})
	if err != nil {
		log.Error("create flat", slog.String("err", err.Error()))

		render.Status(r, http.StatusInternalServerError)
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("success adding flat")

	render.JSON(w, r, gen.Flat{
		HouseId: gen.HouseId(flat.HouseID),
		Id:      gen.FlatId(flat.ID),
		Price:   gen.Price(flat.Price),
		Rooms:   gen.Rooms(flat.Rooms),
		Status:  gen.Status(flat.Status),
	})

}

// PostFlatUpdate implements gen.ServerInterface.
func (h *Handler) PostFlatUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostFlatUpdate"

	reqID := middleware.GetReqID(r.Context())

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)

	log.Info("attempting update status flat")

	userType, ok := r.Context().Value(mv.UserTypeContextKey).(models.UserType)
	if !ok {
		log.Error("failed get user type")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, ok := r.Context().Value(mv.UserIDContextKey).(uuid.UUID)
	if !ok {
		log.Error("failed get user id")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if userType != models.Moderator {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var f gen.PostFlatUpdateJSONRequestBody
	if err := render.DecodeJSON(r.Body, &f); err != nil {
		log.Warn(
			"incorrect json body",
			slog.String("err", err.Error()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := validateFlatUpdate(gen.PostFlatUpdateJSONBody(f))
	if err != nil {
		log.Warn("incorrect json", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	flat, err := h.FlatService.FlatUpdate(
		context.Background(),
		f.Id,
		userID,
		models.Status(*f.Status),
	)
	if err != nil {
		log.Error("failed update", slog.String("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("success update status flat")

	render.JSON(w, r, gen.Flat{
		HouseId: gen.HouseId(flat.HouseID),
		Id:      gen.FlatId(flat.ID),
		Price:   gen.Price(flat.Price),
		Rooms:   gen.Rooms(flat.Rooms),
		Status:  gen.Status(flat.Status),
	})

}

func validateFlatUpdate(user gen.PostFlatUpdateJSONBody) error {
	rules := map[string]string{
		"Id":     "uuid",
		"Status": "eq=approved|eq=created|eq=declined|eq=on moderation",
	}

	validate.RegisterStructValidationMapRules(rules, gen.PostFlatUpdateJSONBody{})
	err := validate.Struct(user)

	return err
}
