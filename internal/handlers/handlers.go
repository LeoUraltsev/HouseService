package handlers

import (
	"log/slog"
	"net/http"

	"github.com/LeoUraltsev/HouseService/internal/gen"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Handler struct {
	HouseService HouseService
	FlatService  FlatService
	AuthService  AuthService

	log *slog.Logger
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

func respError(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	requestID := middleware.GetReqID(r.Context())
	render.Status(r, statusCode)
	render.JSON(w, r, gen.N5xx{
		Code:      &statusCode,
		Message:   msg,
		RequestId: &requestID,
	})
}
