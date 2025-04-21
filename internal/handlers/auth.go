package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/LeoUraltsev/HouseService/internal/gen"
	"github.com/LeoUraltsev/HouseService/internal/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, user models.User) (string, error)
	Register(ctx context.Context, user models.User) (*uuid.UUID, error)
	DummyLogin(ctx context.Context, userType models.UserType) (string, error)
}

type RegResponse struct {
	UserID gen.UserId `json:"user_id"`
}

type LoginResponse struct {
	Token gen.Token `json:"token"`
}

// GetDummyLogin implements gen.ServerInterface.
func (h *Handler) GetDummyLogin(w http.ResponseWriter, r *http.Request, params gen.GetDummyLoginParams) {
	const op = "handlers.GetDummyLogin"

	reqID := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
		slog.String("type", string(params.UserType)),
	)

	log.Info("attempting getting token")

	userType := params.UserType
	if !(userType == gen.Client || userType == gen.Moderator) {
		log.Warn("params invalide")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.DummyLogin(context.Background(), models.UserType(params.UserType))
	if err != nil {
		log.Error("internal error", slog.String("err", err.Error()))
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("success geting token")
	render.JSON(w, r, LoginResponse{
		Token: token,
	})
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

	err := validateLogin(gen.PostLoginJSONBody(req))
	if err != nil {
		log.Warn("failed validate", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log = log.With(slog.String("user_id", req.Id.String()))

	log.Debug("success decode json", slog.String("id", req.Id.String()))

	token, err := h.AuthService.Login(context.Background(), models.User{
		ID:           *req.Id,
		PasswordHash: *req.Password,
	})

	if errors.Is(err, models.ErrUserNotFound) {
		log.Warn("user not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if errors.Is(err, models.ErrInvalideCredentials) {
		log.Warn("invalide credential", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	reqID := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)

	log.Info("attempting registration")

	var req gen.PostRegisterJSONRequestBody

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Warn("incorect json", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := validateRegister(gen.PostRegisterJSONBody(req))
	if err != nil {
		log.Warn("validate failed", slog.String("err", err.Error()))
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

func validateRegister(user gen.PostRegisterJSONBody) error {

	rules := map[string]string{
		"Email":    "required,email",
		"UserType": "eq=client|eq=moderator",
		"Password": "min=6",
	}

	validate.RegisterStructValidationMapRules(rules, gen.PostRegisterJSONBody{})
	err := validate.Struct(user)

	return err
}

func validateLogin(user gen.PostLoginJSONBody) error {
	rules := map[string]string{
		"Id":       "uuid",
		"Password": "min=6",
	}

	validate.RegisterStructValidationMapRules(rules, gen.PostRegisterJSONBody{})
	err := validate.Struct(user)

	return err
}
