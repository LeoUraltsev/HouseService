package handlers

import (
	"context"
	"log/slog"
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
	DummyLogin(ctx context.Context, userType models.UserType) (string, error)
}

type Handler struct {
	HouseService HouseService
	FlatService  FlatService
	AuthService  AuthService

	log *slog.Logger
}

type RegResponse struct {
	UserID gen.UserId `json:"user_id"`
}

type LoginResponse struct {
	Token gen.Token `json:"token"`
}

func New(
	houseService HouseService,
	flatService FlatService,
	authService AuthService,
	log *slog.Logger,
) *Handler {
	return &Handler{
		HouseService: houseService,
		FlatService:  flatService,
		AuthService:  authService,
		log:          log,
	}
}

// GetDummyLogin implements gen.ServerInterface.
func (h *Handler) GetDummyLogin(w http.ResponseWriter, r *http.Request, params gen.GetDummyLoginParams) {
	token, err := h.AuthService.DummyLogin(context.Background(), models.UserType(params.UserType))
	reqID := middleware.GetReqID(r.Context())
	if err != nil {
		code := http.StatusInternalServerError
		render.Status(r, code)
		render.JSON(w, r, gen.N5xx{
			Code:      &code,
			Message:   "что-то пошло не так",
			RequestId: &reqID,
		})
		return
	}
	render.JSON(w, r, LoginResponse{
		Token: token,
	})
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
	const op = "handlers.PostLogin"
	var req gen.PostLoginJSONRequestBody

	reqID := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)

	log.Info("attempting login")

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("decode json", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Id == nil || req.Password == nil {
		log.Warn("requered filds is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug("success decode json", slog.String("id", req.Id.String()))

	token, err := h.AuthService.Login(context.Background(), models.User{
		ID:           *req.Id,
		PasswordHash: *req.Password,
	})
	if err != nil {
		log.Error("login", slog.String("err", err.Error()))
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}
	log.Debug("success get token", slog.String("token", token))
	log.Info("success login", slog.String("user_id", req.Id.String()))

	render.JSON(w, r, &LoginResponse{
		Token: token,
	})
}

// PostRegister implements gen.ServerInterface.
func (h *Handler) PostRegister(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostRegister"

	log := h.log.With(slog.String("op", op))

	log.Info("attempting registration")

	var req gen.PostRegisterJSONRequestBody

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Warn("incorect json", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.AuthService.Register(context.Background(), models.User{
		Email:        string(*req.Email),
		PasswordHash: *req.Password,
		UserType:     models.UserType(*req.UserType),
	})

	if err != nil {
		log.Error("internal error", slog.String("err", err.Error()))
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("success registation", slog.String("user_id", id.String()))
	render.JSON(w, r, &RegResponse{UserID: *id})
}

func respError(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	requestID := middleware.GetReqID(r.Context())
	render.Status(r, statusCode)
	render.JSON(w, r, gen.N5xx{
		Code:      &statusCode,
		Message:   msg,
		RequestId: &requestID,
	})
}
