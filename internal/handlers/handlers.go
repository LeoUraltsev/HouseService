package handlers

import (
	"context"
	"net/http"

	"github.com/LeoUraltsev/HauseService/internal/gen"
	"github.com/LeoUraltsev/HauseService/internal/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type HouseService interface {
}

type FlatService interface {
}

type AuthService interface {
	Login(ctx context.Context, user models.User) (string, error)
	Register(ctx context.Context, user models.User) (*uuid.UUID, error)
}

type Handler struct {
	HouseService HouseService
	FlatService  FlatService
	AuthService  AuthService
}

type RegResponse struct {
	UserID gen.UserId `json:"user_id"`
}

func New(
	houseService HouseService,
	flatService FlatService,
	authService AuthService,
) *Handler {
	return &Handler{
		HouseService: houseService,
		FlatService:  flatService,
		AuthService:  authService,
	}
}

// GetDummyLogin implements gen.ServerInterface.
func (h *Handler) GetDummyLogin(w http.ResponseWriter, r *http.Request, params gen.GetDummyLoginParams) {
	panic("unimplemented")
}

// GetHouseId implements gen.ServerInterface.
func (h *Handler) GetHouseId(w http.ResponseWriter, r *http.Request, id gen.HouseId) {
	panic("unimplemented")
}

// PostFlatCreate implements gen.ServerInterface.
func (h *Handler) PostFlatCreate(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostFlatUpdate implements gen.ServerInterface.
func (h *Handler) PostFlatUpdate(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostHouseCreate implements gen.ServerInterface.
func (h *Handler) PostHouseCreate(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// PostHouseIdSubscribe implements gen.ServerInterface.
func (h *Handler) PostHouseIdSubscribe(w http.ResponseWriter, r *http.Request, id gen.HouseId) {
	panic("unimplemented")
}

// PostLogin implements gen.ServerInterface.
func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	var req gen.PostLoginJSONRequestBody

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.Login(context.Background(), models.User{
		ID:           *req.Id,
		PasswordHash: *req.Password,
	})

	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		rID := middleware.GetReqID(r.Context())
		render.JSON(w, r, &gen.N5xx{
			Code:      new(int),
			Message:   "",
			RequestId: &rID,
		})
		return
	}

	render.JSON(w, r, token)
}

// PostRegister implements gen.ServerInterface.
func (h *Handler) PostRegister(w http.ResponseWriter, r *http.Request) {
	var req gen.PostRegisterJSONRequestBody

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	id, err := h.AuthService.Register(context.Background(), models.User{
		Email:        string(*req.Email),
		PasswordHash: *req.Password,
		UserType:     models.UserType(*req.UserType),
	})

	if err != nil {
		requestID := middleware.GetReqID(r.Context())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, gen.N5xx{
			Code:      new(int),
			Message:   "что-то пошло не так",
			RequestId: &requestID,
		})
		return
	}

	render.JSON(w, r, &RegResponse{UserID: *id})

}
